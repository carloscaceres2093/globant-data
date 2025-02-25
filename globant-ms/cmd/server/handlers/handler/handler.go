package handlers

import (
	"fmt"
	"globant-ms/internal/service"
	"globant-ms/local-lib/web"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	ErrInvalidBody         = "invalid body"
	ErrInvalidAccountID    = "invalid account_id"
	ErrInternalServerError = "internal server error"
	ErrStreamFile          = "stream file error"
	fieldRequestID         = "request_id"
	fieldEvent             = "event"
	fieldError             = "error"
	ErrParseForm           = "Unable to parse form"
	ErrRetrieveFile        = "Unable to retrieve file"
	ErrReadFile            = "Unable to read file"
	ErrProcessing          = "Failed processing file "
	ErrYearFilter          = "Year filter is missing"
)

type Handler struct {
	Service
}

type Service interface {
	UploadFile(userCode string, handler *multipart.FileHeader, request io.Reader, fileType string) error
	GetQuarterData(service.QueryParams) ([]service.QuarterMetrics, error)
	GetHiredData(service.QueryParams) ([]service.HiredMetrics, error)
}

func NewHandler(service Service) Handler {
	return Handler{
		service,
	}
}

func (rh *Handler) UploadHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		fmt.Println(err)
		web.RespondJSON(w, web.Error{Message: ErrParseForm}, http.StatusBadRequest)
		return
	}
	fileType := req.FormValue("fileType")
	if fileType == "" {
		http.Error(w, "type field is required", http.StatusBadRequest)
		return
	}
	file, handler, err := req.FormFile("file")
	if err != nil {
		web.RespondJSON(w, web.Error{Message: ErrRetrieveFile}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Send the file to the internal microservice
	err = rh.Service.UploadFile("test", handler, file, fileType)
	if err != nil {
		web.RespondJSON(w, web.Error{Message: ErrProcessing}, http.StatusInternalServerError)
		return
	}

	web.RespondJSON(w, []byte("File uploaded"), http.StatusOK)
}

func (rh *Handler) GetQuarterData(w http.ResponseWriter, req *http.Request) {
	var departmentParam *string
	var jobParam *string
	year := req.URL.Query().Get("year")
	if year == "" {
		web.RespondJSON(w, web.Error{Message: ErrYearFilter}, http.StatusInternalServerError)
		return
	}
	department := req.URL.Query().Get("department_name")
	departmentParam = &department
	if department == "" {
		departmentParam = nil
	}
	job := req.URL.Query().Get("job_name")
	jobParam = &job
	if job == "" {
		jobParam = nil
	}
	parmas := service.QueryParams{
		Year:           &year,
		DepartmentName: departmentParam,
		JobName:        jobParam,
	}
	quarterData, err := rh.Service.GetQuarterData(parmas)
	if err != nil {
		web.RespondJSON(w, web.Error{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	web.RespondJSON(w, quarterData, http.StatusOK)
}

func (rh *Handler) GetHiredData(w http.ResponseWriter, req *http.Request) {
	var departmentParam *string
	year := req.URL.Query().Get("year")
	if year == "" {
		web.RespondJSON(w, web.Error{Message: ErrYearFilter}, http.StatusInternalServerError)
		return
	}
	department := req.URL.Query().Get("department_name")
	departmentParam = &department
	if department == "" {
		departmentParam = nil
	}
	parmas := service.QueryParams{
		Year:           &year,
		DepartmentName: departmentParam,
	}
	hiredData, err := rh.Service.GetHiredData(parmas)
	if err != nil {
		web.RespondJSON(w, web.Error{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	web.RespondJSON(w, hiredData, http.StatusOK)
}
