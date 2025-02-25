package handlers

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) UploadFile(_ string, _ *multipart.FileHeader, _ io.Reader) error {
	args := s.Called()
	return args.Error(0)
}

func TestUploadHandler(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte("test file content"))
	writer.Close()

	var tests = []struct {
		name               string
		service            *ServiceMock
		request            *bytes.Buffer
		contentType        string
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			name: "Ok - Upload",
			service: func() *ServiceMock {
				m := ServiceMock{}
				m.On("UploadFile", mock.Anything).Return(nil)
				return &m
			}(),

			request:            body,
			contentType:        writer.FormDataContentType(),
			expectedResponse:   "File accepted",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Fail - Parse Form",
			service: func() *ServiceMock {
				m := ServiceMock{}
				m.On("UploadFile", mock.Anything).Return(nil)
				return &m
			}(),
			request:            body,
			contentType:        "application/csv",
			expectedResponse:   "{\"message\":\"Unable to parse form\"}",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := Handler{Service: tt.service}
			req := httptest.NewRequest(http.MethodPost, "/v1/upload", tt.request)
			req.Header.Set("Content-Type", tt.contentType)
			rctx := chi.NewRouteContext()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			recorder := httptest.NewRecorder()

			handler := http.HandlerFunc(service.UploadHandler)
			handler.ServeHTTP(recorder, req)

			response := recorder.Result()
			resBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedStatusCode, response.StatusCode)
			require.Equal(t, tt.expectedResponse, string(resBody))
		})
	}
}
