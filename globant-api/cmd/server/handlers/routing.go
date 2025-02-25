package server

import (
	apiHandlers "globant-api/cmd/server/handlers/handler"
	health "globant-api/cmd/server/handlers/health"
	"globant-api/internal/platform/config"
	"globant-api/local-lib/web"
)

const (
	healthCheckPattern   = "/health-check"
	microServicePattern  = "/globant-api"
	v1                   = "/v1"
	apiUploadFilePattern = "/upload"
	tableNamePattern     = "/{table}"
)

// Associate routes with handlers here
func routes(webApp *web.App, services *services, cfg config.Config) {
	healthHandler := health.New(services.healthService)

	// --- Health route ---
	hGroup := webApp.Group("/" + cfg.AppName)
	hGroup.Get(healthCheckPattern, healthHandler.HealthCheck)

	// Endpoints V1
	v1Group := webApp.Group(microServicePattern + v1)
	apiHandler := apiHandlers.NewHandler(services.service)
	v1Group.Post(apiUploadFilePattern, apiHandler.UploadHandler)
	

}
