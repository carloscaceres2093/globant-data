package handlers

import (
	"fmt"
	"globant-auth-ms/internal/health"
	"globant-auth-ms/local-lib/web"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	microServiceName   = "/data-api"
	healthCheckPattern = "/health-check"
)

type HealthServiceMock struct {
	mock.Mock
}

func (hm *HealthServiceMock) Check() health.Response {
	args := hm.Called()
	return args.Get(0).(health.Response)
}

func TestHealthCheck(t *testing.T) {
	healthOk := health.Response{
		Status: "ok",
	}
	var tests = []struct {
		name               string
		service            *HealthServiceMock
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			name: "Ok",
			service: func() *HealthServiceMock {
				m := HealthServiceMock{}
				m.On("Check", mock.Anything).Return(healthOk)
				return &m
			}(),
			expectedResponse:   "{\"status\":\"ok\"}",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := web.NewWebApp("local")
			handler := New(tt.service)
			group := app.Group(microServiceName)
			group.Get(healthCheckPattern, handler.HealthCheck)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprint(microServiceName, healthCheckPattern), nil)

			rr := httptest.NewRecorder()
			app.Router.ServeHTTP(rr, r)

			res := rr.Result()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedStatusCode, res.StatusCode)
			require.Equal(t, tt.expectedResponse, string(resBody))
		})
	}
}
