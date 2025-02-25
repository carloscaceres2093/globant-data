package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"strconv"
)

type Error struct {
	Message        string
	AdditionalInfo string
}

var (
	ErrUnexpectedMSResponse = errors.New("unexpected error")
	ErrFileType             = errors.New("file type error")
	ErrFormatModel          = errors.New("error getting columnsm model")
	ErrParseFile            = errors.New("error parsing file")
	ErrReadFile             = errors.New("error reading file")
	ErrProcessingFile       = errors.New("erorr processing file")
	ErrGettingData          = errors.New("erorr getting data")
	columnMapping           = map[string]string{
		"ID":        "ID",
		"CreatedAt": "CreatedAt",
		"UpdatedAt": "UpdatedAt",
	}
)

type Store interface {
	JobsStore(file FileModel, records [][]string, model interface{}, columns []string) error
	GetQuarters(params QueryParams) ([]QuarterMetrics, error)
	GetHired(params QueryParams) ([]HiredMetrics, error)
}

type Service struct {
	store Store
}

func NewService(store Store) Service {
	return Service{
		store: store,
	}
}

func (receiver Service) UploadFile(userCode string, handler *multipart.FileHeader, file io.Reader, fileType string) error {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return ErrReadFile
	}

	clientFile := FileModel{
		UserCode:  userCode,
		FileBytes: fileBytes,
		FileName:  handler.Filename,
	}

	// Parse the file (e.g., CSV)
	reader := csv.NewReader(bytes.NewReader(fileBytes))
	records, err := reader.ReadAll()
	if err != nil {
		return ErrParseFile
	}

	switch {
	case fileType == "jobs":
		columns, err := getColumnsFromModel(&Job{})
		if err != nil {
			return ErrFormatModel
		}
		var jobs []Job
		err = receiver.store.JobsStore(clientFile, records, &jobs, columns)
	case fileType == "departments":
		columns, err := getColumnsFromModel(&Department{})
		if err != nil {
			return ErrFormatModel
		}
		var departments []Department
		err = receiver.store.JobsStore(clientFile, records, &departments, columns)
	case fileType == "employees":
		columns, err := getColumnsFromModel(&Employee{})
		if err != nil {
			return ErrFormatModel
		}
		var employees []Employee
		err = receiver.store.JobsStore(clientFile, records, &employees, columns)
	default:
		return ErrFileType
	}
	if err != nil {
		if errors.Is(err, ErrProcessingFile) {
			return ErrReadFile
		} else {
			return err
		}
	}

	return nil
}

func (receiver Service) GetQuarterData(queryParams QueryParams) ([]QuarterMetrics, error) {

	values, err := receiver.store.GetQuarters(queryParams)
	if err != nil {
		return []QuarterMetrics{}, err
	}

	return values, nil
}

func (receiver Service) GetHiredData(queryParams QueryParams) ([]HiredMetrics, error) {

	values, err := receiver.store.GetHired(queryParams)
	if err != nil {
		return []HiredMetrics{}, err
	}

	return values, nil
}

func getColumnsFromModel(model interface{}) ([]string, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct or a pointer to a struct")
	}
	var columns []string
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		_, ok := columnMapping[field.Name]
		if !ok {
			columns = append(columns, field.Name)
		}
	}

	return columns, nil
}
func processData(model interface{}, columns []string, records [][]string) (interface{}, error) {
	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type().Elem()

	for _, record := range records {
		newModel := reflect.New(modelType).Elem()
		for i, column := range columns {
			field := newModel.FieldByName(column)
			if !field.IsValid() {
				return model, fmt.Errorf("field %s not found in model", column)
			}

			recordValue := record[i]
			switch field.Kind() {
			case reflect.Int:
				intValue, err := strconv.Atoi(recordValue)
				if err != nil {
					return model, fmt.Errorf("failed to convert %s to int: %v", recordValue, err)
				}
				field.SetInt(int64(intValue))
			case reflect.String:
				field.SetString(recordValue)
			default:
				return model, fmt.Errorf("unsupported field type: %v", field.Kind())
			}
		}

		// Append the new model to the slice
		modelValue.Set(reflect.Append(modelValue, newModel))
	}
	fmt.Println(modelValue)
	// Insert the batch data in a single transaction
	fmt.Println(model)
	return model, nil
}

// func (receiver ReportService) Get(ctx context.Context, organizationCode string, accountCode string, params Params) (clientReport.ListsResponse, error) {
// 	queryParam := queryParamCompose(params)
// 	response, err := receiver.client.GetTransactions(ctx, organizationCode, accountCode, queryParam)
// 	if err != nil {
// 		if errors.Is(err, clientReport.ErrNotRequest) {
// 			errRet := ErrNotFound
// 			errRet.AdditionalInfo = err.Error()
// 			return clientReport.ListsResponse{}, errRet
// 		} else if errors.Is(err, clientReport.ErrUnexpectedMSResponse) {
// 			errRet := ErrUnexpectedResponse
// 			errRet.AdditionalInfo = err.Error()
// 			return clientReport.ListsResponse{}, errRet
// 		}

// 		return clientReport.ListsResponse{}, err
// 	}
// 	return response, nil
// }

// func (receiver ReportService) GetDownload(ctx context.Context, organizationCode string, reportID string) (io.ReadCloser, error) {

// 	response, err := receiver.client.GetDownloadTransactions(ctx, organizationCode, reportID)
// 	if err != nil {
// 		if errors.Is(err, clientReport.ErrNotRequest) {
// 			errRet := ErrNotFound
// 			errRet.AdditionalInfo = err.Error()
// 			return nil, errRet
// 		} else if errors.Is(err, clientReport.ErrUnexpectedMSResponse) {
// 			errRet := ErrUnexpectedResponse
// 			errRet.AdditionalInfo = err.Error()
// 			return nil, errRet
// 		}

// 		return nil, err
// 	}
// 	responseAws, err := receiver.aws.GetStreamFlie(response.Link)
// 	if err != nil {
// 		if errors.Is(err, clientReport.ErrNotRequest) {
// 			errRet := ErrNotFound
// 			errRet.AdditionalInfo = err.Error()
// 			return nil, errRet
// 		} else if errors.Is(err, clientReport.ErrUnexpectedMSResponse) {
// 			errRet := ErrUnexpectedResponse
// 			errRet.AdditionalInfo = err.Error()
// 			return nil, errRet
// 		}

// 		return nil, err
// 	}
// 	return responseAws, nil
// }
// func queryParamCompose(params Params) string {

// 	paramsMap := structs.Map(params)
// 	count := 0
// 	for key, value := range paramsMap {
// 		queryParam += (paramKeys[key].(string) + equal + value.(string) + and)
// 		count++
// 		if count == len(paramsMap)-1 {
// 			and = ""
// 		}
// 	}
// 	return queryParam
// }
