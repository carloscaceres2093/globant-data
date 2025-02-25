package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"globant-auth-ms/internal/auth"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type AuthServiceMock struct {
	mock.Mock
}

func (s *AuthServiceMock) ValidateToken(_ string, _ string) (auth.AuthResponse, error) {
	args := s.Called()
	return args.Get(0).(auth.AuthResponse), args.Error(1)
}
func (s *AuthServiceMock) CreateUser(_ auth.AuthRequest) (auth.AuthResponse, error) {
	args := s.Called()
	return args.Get(0).(auth.AuthResponse), args.Error(1)
}
func TestCreateUserHandler(t *testing.T) {
	requestOk, err := json.Marshal(auth.AuthRequest{
		UserName: "test",
	})
	require.NoError(t, err)
	requestEmpty, err := json.Marshal(auth.AuthRequest{
		UserName: "",
	})
	require.NoError(t, err)

	var tests = []struct {
		name               string
		service            *AuthServiceMock
		request            *bytes.Reader
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			name: "Ok - Create",
			service: func() *AuthServiceMock {
				m := AuthServiceMock{}
				m.On("CreateUser", mock.Anything).Return(nil)
				return &m
			}(),
			request:            bytes.NewReader(requestOk),
			expectedResponse:   "ok",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Fail - Parse Form",
			service: func() *AuthServiceMock {
				m := AuthServiceMock{}
				m.On("CreateUser", mock.Anything).Return(nil)
				return &m
			}(),
			request:            bytes.NewReader(requestEmpty),
			expectedResponse:   "{\"message\":\"Unable to parse form\"}",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := AuthHandler{AuthService: tt.service}
			req := httptest.NewRequest(http.MethodPost, "/v1/user", tt.request)
			rctx := chi.NewRouteContext()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			recorder := httptest.NewRecorder()

			handler := http.HandlerFunc(service.CreateUser)
			handler.ServeHTTP(recorder, req)

			response := recorder.Result()
			resBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedStatusCode, response.StatusCode)
			require.Equal(t, tt.expectedResponse, string(resBody))
		})
	}
}
