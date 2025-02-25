package handlers

import (
	"net/http"

	"globant-auth-ms/internal/health"
	"globant-auth-ms/local-lib/web"
)

type HealthService interface {
	Check() health.Response
}
type Health struct {
	HealthService
}

func New(service HealthService) Health {
	return Health{
		service,
	}
}

func (rh *Health) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	web.RespondJSON(w, rh.HealthService.Check(), http.StatusOK)
}
