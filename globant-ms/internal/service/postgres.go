package service

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"globant-ms/internal/platform/config"
	"globant-ms/local-lib/env"
	"globant-ms/local-lib/path"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (job *Job) TableName() string {
	return "jobs"
}

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(cfg config.Database) (Postgres, error) {
	gdb, err := connect(cfg)
	if err != nil {
		return Postgres{}, err
	}
	if err = applyMigrations(cfg); err != nil && err != migrate.ErrNoChange {
		return Postgres{}, fmt.Errorf("an error ocurred applying migrations: %w", err)
	}
	return Postgres{
		DB: gdb,
	}, nil
}

func (r Postgres) JobsStore(employees FileModel, records [][]string, model interface{}, columns []string) error {
	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type().Elem()

	for _, record := range records {

		newModel := reflect.New(modelType).Elem()
		isValid := true
		for i, column := range columns {
			if !isValid {
				continue
			}
			field := newModel.FieldByName(column)
			if !field.IsValid() {
				return fmt.Errorf("field %s not found in model", column)
			}

			recordValue := record[i]
			switch field.Kind() {
			case reflect.Int, reflect.Int64:
				intValue, err := strconv.ParseInt(recordValue, 10, 64)
				if err != nil {
					isValid = false
				}
				field.SetInt(intValue)
			case reflect.String:
				field.SetString(recordValue)
			case reflect.Struct:

				if field.Type() == reflect.TypeOf(time.Time{}) {
					timeValue, err := time.Parse(time.RFC3339, recordValue)
					if err != nil {
						isValid = false
					}
					field.Set(reflect.ValueOf(timeValue))
				}
			default:
				return fmt.Errorf("unsupported field type: %v", field.Kind())
			}

		}
		if isValid {
			createdAtField := newModel.FieldByName("CreatedAt")
			createdAtField.Set(reflect.ValueOf(time.Now()))
			updatedAtField := newModel.FieldByName("UpdatedAt")
			updatedAtField.Set(reflect.ValueOf(time.Now()))
			modelValue.Set(reflect.Append(modelValue, newModel))
		}
	}

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(model).Error; err != nil {
			return fmt.Errorf("failed to insert batch data: %v", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil

}

func (r Postgres) GetQuarters(params QueryParams) ([]QuarterMetrics, error) {
	var values []QuarterMetrics
	var args []interface{}
	query := `SELECT  d.department_name, 
		j.job_name, 
		EXTRACT(YEAR FROM e.hire_time) AS year,
		SUM(CASE WHEN EXTRACT(QUARTER FROM e.hire_time) = 1 THEN 1 ELSE 0 END) AS Q1,
		SUM(CASE WHEN EXTRACT(QUARTER FROM e.hire_time) = 2 THEN 1 ELSE 0 END) AS Q2,
		SUM(CASE WHEN EXTRACT(QUARTER FROM e.hire_time) = 3 THEN 1 ELSE 0 END) AS Q3,
		SUM(CASE WHEN EXTRACT(QUARTER FROM e.hire_time) = 4 THEN 1 ELSE 0 END) AS Q4
	FROM employees as e
	inner join jobs as j on j.job_id=e.job_id
	inner join departments as d on d.department_id=e.department_id
	where 1=1 %s
	GROUP BY 1, 2, 3
	ORDER BY 1,2;`

	finalFilter := ""

	if params.Year != nil {
		finalFilter += " AND EXTRACT(YEAR FROM hire_time) = ?"
		args = append(args, *params.Year)
	}
	if params.DepartmentName != nil {
		finalFilter += " AND department_name = ?"
		args = append(args, *params.DepartmentName)
	}
	if params.JobName != nil {
		finalFilter += " AND job_name = ?"
		args = append(args, *params.JobName)
	}

	finalQuery := fmt.Sprintf(query, finalFilter)
	tx := r.DB.Raw(finalQuery, args...).Scan(&values)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return []QuarterMetrics{}, ErrGettingData
		}
		return []QuarterMetrics{}, tx.Error
	}
	return values, nil

}
func (r Postgres) GetHired(params QueryParams) ([]HiredMetrics, error) {

	var values []HiredMetrics
	var args []interface{}
	query := `SELECT 
			d.department_id,
			d.department_name,
			EXTRACT(YEAR FROM e.hire_time) AS year,
			COUNT(e.id) as hired
		FROM employees as e
		inner join departments as d on d.department_id=e.department_id
		WHERE 1=1 %s 
		GROUP BY 1,2,3
		HAVING COUNT(e.id) > (
			SELECT AVG(employee_count) 
			FROM (
				SELECT COUNT(id) AS employee_count 
				FROM employees 
				WHERE 1=1
				%s
				GROUP BY department_id
			) subquery
		)
		ORDER BY 4 DESC;`
	finalFilter := ""

	if params.Year != nil {
		finalFilter += " AND EXTRACT(YEAR FROM hire_time) = ?"
		args = append(args, *params.Year)
	}
	if params.DepartmentName != nil {
		finalFilter += " AND department_name = ?"
		args = append(args, *params.DepartmentName)
	}

	subQueryFilter := finalFilter
	if params.Year != nil {
		args = append(args, *params.Year)
	}
	if params.DepartmentName != nil {
		args = append(args, *params.DepartmentName)
	}
	finalQuery := fmt.Sprintf(query, finalFilter, subQueryFilter)
	tx := r.DB.Raw(finalQuery, args...).Scan(&values)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return []HiredMetrics{}, ErrGettingData
		}
		return []HiredMetrics{}, tx.Error
	}
	return values, nil

}

func connect(cfg config.Database) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: cfg.LogLevel,
		},
	)

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	instance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	} // esto impide levantar en kingdom
	fmt.Print("error")

	db, err := instance.DB()
	if err != nil {
		return nil, err
	}

	maxIdleSize, err := strconv.Atoi(env.GetEnv("MAX_IDLE_SIZE"))
	if err != nil {
		return nil, err
	}

	maxOpenSize, err := strconv.Atoi(env.GetEnv("MAX_OPEN_SIZE"))
	if err != nil {
		return nil, err
	}

	maxLifeTime, err := time.ParseDuration(env.GetEnv("MAX_LIFE_TIME") + "ms")
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleSize)
	db.SetMaxOpenConns(maxOpenSize)
	db.SetConnMaxLifetime(maxLifeTime)

	return instance, nil
}

func applyMigrations(cfg config.Database) error {
	currentPath := fmt.Sprintf("file:///%s/migrations", path.GetMainPath())
	url := generatePgUrl(cfg)
	m, err := migrate.New(currentPath, url)
	if err != nil {
		return err
	}

	return m.Up()
}

func generatePgUrl(cfg config.Database) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}
