package service

import (
	"errors"
	"fmt"
	apiClient "globant-api/internal/service/client"
	"io"
	"mime/multipart"
	"strings"
)

type Error struct {
	Message        string
	AdditionalInfo string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.AdditionalInfo)
}

var (
	ErrNotFound           = errors.New("bad_request_ms")
	ErrReadFile           = errors.New("error_reading")
	ErrUnexpectedResponse = errors.New("unexpected_response")
	ErrFileType           = errors.New("unexpected_file")
)

type Client interface {
	UploadFile(file apiClient.FileModel) error
}

type Service struct {
	client Client
}

func NewService(client Client) Service {
	return Service{
		client: client,
	}
}

func (receiver Service) UploadFile(userCode string, handler *multipart.FileHeader, file io.Reader) error {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return ErrReadFile
	}
	clientFile := apiClient.FileModel{
		UserCode:  userCode,
		FileBytes: fileBytes,
		FileName:  handler.Filename,
	}
	fileName := strings.ToLower(handler.Filename)
	switch {
	case strings.Contains(fileName, "employees"):
		clientFile.FileType = "employees"
	case strings.Contains(fileName, "departments"):
		clientFile.FileType = "departments"
	case strings.Contains(fileName, "jobs"):
		clientFile.FileType = "jobs"
	default:
		return ErrFileType
	}
	err = receiver.client.UploadFile(clientFile)
	if err != nil {
		if errors.Is(err, apiClient.ErrProcessingFile) {
			return ErrReadFile
		} else if errors.Is(err, apiClient.ErrUnexpectedMSResponse) {
			return ErrUnexpectedResponse
		} else if errors.Is(err, apiClient.ErrInvalidResponse) {
			return ErrUnexpectedResponse
		} else {
			return err
		}
	}

	return nil
}
