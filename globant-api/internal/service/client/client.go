package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"globant-api/internal/platform/terror"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"globant-api/internal/platform/config"
)

var (
	ErrUnexpectedMSResponse = terror.Error{
		Message: "unexpected microservice response",
	}
	ErrInvalidResponse = terror.Error{
		Message: "invalid response. %s",
	}
	ErrProcessingFile = terror.Error{
		Message: "error processing file",
	}
	ErrAuth = terror.Error{
		Message: "auth error",
	}
)

const (
	defaultTimeout = 10 * time.Second

	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
	errWrappedFormat  = "%w: %s"
)

type Client struct {
	cfgClient config.Config
	client    *http.Client
}

func NewClient(cfg config.Config) Client {

	client := buildClientHttp(defaultTimeout)

	return Client{
		cfgClient: cfg,
		client:    client,
	}
}

func (r Client) UploadFile(file FileModel) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := writer.WriteField("fileType", file.FileType)
	if err != nil {
		return err
	}
	part, err := writer.CreateFormFile("file", file.FileName)
	if err != nil {
		return err
	}
	part.Write(file.FileBytes)
	contentTypeFile := writer.FormDataContentType()
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, r.cfgClient.UploadFile, body)
	if err != nil {
		return ErrInvalidResponse
	}
	req.Header.Set(contentTypeHeader, contentTypeFile)
	req.Header.Set("X-user-code", file.UserCode)

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ErrProcessingFile.Code = strconv.Itoa(resp.StatusCode)
		return ErrProcessingFile
	}
	return nil
}

func AuthValidation(userCode string, token string, authUrl string) (AuthResponse, error) {
	url := fmt.Sprintf(authUrl, userCode, token)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return AuthResponse{}, err
	}
	client := &http.Client{}
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	resp, err := client.Do(req)
	if err != nil {
		return AuthResponse{}, ErrAuth
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return AuthResponse{}, ErrAuth
	}
	defer resp.Body.Close()

	var data AuthResponse
	if resp.StatusCode != http.StatusOK {
		ErrInvalidResponse.Message = fmt.Sprint(ErrInvalidResponse.Message, responseBody)
		ErrInvalidResponse.Code = strconv.Itoa(resp.StatusCode)
		return AuthResponse{}, ErrInvalidResponse
	}

	if err = json.Unmarshal(responseBody, &data); err != nil {
		return AuthResponse{}, ErrUnexpectedMSResponse
	}

	return data, nil

}

func buildClientHttp(timeoutV time.Duration) *http.Client {
	timeout := defaultTimeout
	if timeoutV != 0 {
		timeout = timeoutV
	}

	return &http.Client{
		Timeout: timeout,
	}
}
