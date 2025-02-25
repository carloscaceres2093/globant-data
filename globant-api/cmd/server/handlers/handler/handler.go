package handlers

import (
	"fmt"
	"globant-api/local-lib/web"
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
)

type Handler struct {
	Service
}

type Service interface {
	UploadFile(userCode string, handler *multipart.FileHeader, request io.Reader) error
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
	file, handler, err := req.FormFile("file")
	if err != nil {
		web.RespondJSON(w, web.Error{Message: ErrRetrieveFile}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Send the file to the internal microservice
	err = rh.Service.UploadFile("test", handler, file)
	if err != nil {
		web.RespondJSON(w, web.Error{Message: ErrProcessing}, http.StatusBadRequest)
		return
	}

	web.RespondJSON(w, []byte("File accepted"), http.StatusOK)
}
