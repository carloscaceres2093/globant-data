package service

import (
	"bytes"
	"errors"
	apiClient "globant-api/internal/service/client"
	"io"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ClientMock struct {
	mock.Mock
}

func (s *ClientMock) UploadFile(_ apiClient.FileModel) error {
	args := s.Called()
	return args.Error(0)
}

func TestService_UploadFile(t *testing.T) {
	userCode := "test"
	fileContent := "dummy file content"
	fileReader := io.NopCloser(bytes.NewReader([]byte(fileContent)))
	handler := &multipart.FileHeader{Filename: "employeestest.csv"}
	var tests = []struct {
		name          string
		client        *ClientMock
		userCode      string
		fileHeader    *multipart.FileHeader
		file          io.Reader
		expectedError error
	}{
		{
			name: "Ok - UploadFile",
			client: func() *ClientMock {
				m := ClientMock{}
				m.On("UploadFile", mock.Anything).Return(nil)
				return &m
			}(),
			userCode:   userCode,
			fileHeader: handler,
			file:       fileReader,
		},
		{
			name: "Fail - Processing",
			client: func() *ClientMock {
				m := ClientMock{}
				m.On("UploadFile", mock.Anything).Return(apiClient.ErrProcessingFile)
				return &m
			}(),
			userCode:      userCode,
			fileHeader:    handler,
			file:          fileReader,
			expectedError: errors.New("error_reading"),
		},
		{
			name: "Fail - UnexpectedError",
			client: func() *ClientMock {
				m := ClientMock{}
				m.On("UploadFile", mock.Anything).Return(apiClient.ErrUnexpectedMSResponse)
				return &m
			}(),
			userCode:      userCode,
			fileHeader:    handler,
			file:          fileReader,
			expectedError: errors.New("unexpected_response"),
		},
		{
			name: "Fail - InvalidResponse",
			client: func() *ClientMock {
				m := ClientMock{}
				m.On("UploadFile", mock.Anything).Return(apiClient.ErrInvalidResponse)
				return &m
			}(),
			userCode:      userCode,
			fileHeader:    handler,
			file:          fileReader,
			expectedError: errors.New("unexpected_response"),
		},
		{
			name: "Fail - Internal error",
			client: func() *ClientMock {
				m := ClientMock{}
				m.On("UploadFile", mock.Anything).Return(errors.New("random error"))
				return &m
			}(),
			userCode:      userCode,
			fileHeader:    handler,
			file:          fileReader,
			expectedError: errors.New("random error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.client)
			err := service.UploadFile(tt.userCode, tt.fileHeader, tt.file)
			require.Equal(t, tt.expectedError, err)
		})
	}
}
