package service

import "time"

type Job struct {
	ID        int64     `gorm:"index;column:id;primaryKey"`
	JobID     int64     `gorm:"index;column:job_id;" json:"id"`
	JobName   string    `gorm:"column:job_name;" json:"job_name"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Department struct {
	ID             int64     `gorm:"index;column:id;primaryKey"`
	DepartmentID   int64     `gorm:"index;column:department_id;" json:"id"`
	DepartmentName string    `gorm:"column:department_name;" json:"department_name"`
	CreatedAt      time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Employee struct {
	ID           int64     `gorm:"index;column:id;primaryKey"`
	EmployeeID   int64     `gorm:"index;column:employee_id;" json:"id"`
	EmployeeName string    `gorm:";column:employee_name;" json:"employee_name"`
	HireTime     time.Time `gorm:"column:hire_time;" json:"datetime"`
	DepartmentID int64     `gorm:"column:department_id;" json:"department_id"`
	JobID        int64     `gorm:"column:job_id;" json:"job_id"`
	CreatedAt    time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type FileModel struct {
	UserCode  string
	FileBytes []byte
	FileName  string
}

type QuarterMetrics struct {
	JobName        string `gorm:"column:job_name;" json:"job"`
	DepartmentName string `gorm:"column:department_name;" json:"department"`
	Q1             int64  `gorm:"column:q1;" json:"q_1"`
	Q2             int64  `gorm:"column:q2;" json:"q_2"`
	Q3             int64  `gorm:"column:q3;" json:"q_3"`
	Q4             int64  `gorm:"column:q4;" json:"q_4"`
	Year           string `gorm:"column:year;" json:"year"`
}
type HiredMetrics struct {
	DepartmentID   string `gorm:"column:department_id;" json:"id"`
	DepartmentName string `gorm:"column:department_name;" json:"department"`
	Hired          int64  `gorm:"column:hired;" json:"hired"`
	Year           string `gorm:"column:year;" json:"year"`
}
type QueryParams struct {
	Year           *string `json:"year"`
	DepartmentName *string `json:"department_name"`
	JobName        *string `json:"job_name"`
}
