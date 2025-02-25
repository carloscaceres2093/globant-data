package service

import (
	"bytes"
	"io"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type StoreMock struct {
	mock.Mock
}

func (s *StoreMock) JobsStore(_ FileModel, _ [][]string, _ interface{}, _ []string) error {
	args := s.Called()
	return args.Error(0)
}
func (s *StoreMock) GetQuarters(_ QueryParams) ([]QuarterMetrics, error) {
	args := s.Called()
	return args.Get(0).([]QuarterMetrics), args.Error(1)
}
func (s *StoreMock) GetHired(_ QueryParams) ([]HiredMetrics, error) {
	args := s.Called()
	return args.Get(0).([]HiredMetrics), args.Error(1)
}

func TestService_UploadFile(t *testing.T) {
	userCode := "test"
	fileType := "employees"
	fileContent := "dummy file content"
	fileReader := io.NopCloser(bytes.NewReader([]byte(fileContent)))
	handler := &multipart.FileHeader{Filename: "employeestest.csv"}

	var tests = []struct {
		name          string
		store         *StoreMock
		userCode      string
		fileHeader    *multipart.FileHeader
		file          io.Reader
		fileType      string
		expectedError error
	}{
		{
			name: "Ok - UploadFile",
			store: func() *StoreMock {
				m := StoreMock{}
				m.On("JobsStore", mock.Anything).Return(nil)
				return &m
			}(),
			userCode:   userCode,
			fileHeader: handler,
			file:       fileReader,
			fileType:   fileType,
		},
		{
			name: "Fail - Saving",
			store: func() *StoreMock {
				m := StoreMock{}
				m.On("JobsStore", mock.Anything).Return(nil)
				return &m
			}(),
			userCode:   userCode,
			fileHeader: handler,
			file:       fileReader,
			fileType:   fileType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.store)
			err := service.UploadFile(tt.userCode, tt.fileHeader, tt.file, tt.fileType)
			require.Equal(t, tt.expectedError, err)
		})
	}
}
