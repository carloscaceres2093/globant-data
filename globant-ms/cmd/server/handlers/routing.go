package server

import (
	apiHandlers "globant-ms/cmd/server/handlers/handler"
	health "globant-ms/cmd/server/handlers/health"
	"globant-ms/internal/platform/config"
	"globant-ms/local-lib/web"
)

const (
	healthCheckPattern    = "/health-check"
	microServicePattern   = "/globant-ms"
	v1                    = "/v1"
	apiUploadFilePattern  = "/upload"
	fileTypePattern       = "/{file_type}"
	quarterMetricsPattern = "quarter_metrics"
	hiredMetricsPattern   = "hired_metrics"
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
	v1Group.Get(quarterMetricsPattern, apiHandler.GetQuarterData)
	v1Group.Get(hiredMetricsPattern, apiHandler.GetHiredData)

}
