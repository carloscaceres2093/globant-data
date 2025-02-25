package handlers

import (
	"net/http"

	"globant-api/internal/health"
	"globant-api/local-lib/web"
)

// HealthService represents the methods available for Health.
type HealthService interface {
	Check() health.Response
}

// Health contains what is necessary to check health.
type Health struct {
	HealthService
}

// New creates a Health.
func New(service HealthService) Health {
	return Health{
		service,
	}
}

// HealthCheck checks Health.
func (rh *Health) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	web.RespondJSON(w, rh.HealthService.Check(), http.StatusOK)
}
