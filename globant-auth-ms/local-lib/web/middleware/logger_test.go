package middleware

import (
	"bytes"
	"globant-auth-ms/local-lib/log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogMiddleware(t *testing.T) {
	log := logrus.StandardLogger()
	mid := Logger(log)
	handler := testHandler{}

	req, err := http.NewRequest(http.MethodGet, "https://www.google.com/search?param=1", bytes.NewReader([]byte(`{"field_1": "hola","field_2": "chau"}`)))
	require.Nil(t, err)

	rec := httptest.NewRecorder()

	mid(&handler).ServeHTTP(rec, req)

	assert.Equal(t, 1, handler.timesCalled)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotNil(t, handler.logger)
	assert.Same(t, log, handler.logger)
}

type testHandler struct {
	timesCalled int
	logger      *logrus.Logger
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.timesCalled++

	h.logger = log.FromContext(r.Context())
	w.WriteHeader(http.StatusOK)
}
