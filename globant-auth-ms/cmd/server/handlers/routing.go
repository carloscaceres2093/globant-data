package server

import (
	authHandlers "globant-auth-ms/cmd/server/handlers/auth"
	health "globant-auth-ms/cmd/server/handlers/health"
	"globant-auth-ms/internal/platform/config"
	"globant-auth-ms/local-lib/web"
)

const (
	healthCheckPattern  = "/health-check"
	microServicePattern = "/globant-auth-ms"
	v1                  = "/v1"
	userPattern         = "/user"
	userCodePattern     = "/{user_code}"
)

// Associate routes with handlers here
func routes(webApp *web.App, services *authServices, cfg config.Config) {
	healthHandler := health.New(services.healthService)

	// --- Health route ---
	hGroup := webApp.Group("/" + cfg.AppName)
	hGroup.Get(healthCheckPattern, healthHandler.HealthCheck)

	// Endpoints V1
	v1Group := webApp.Group(microServicePattern + v1)
	authHandler := authHandlers.NewAuthHandler(services.authService)
	v1Group.Post(userPattern, authHandler.CreateUser)
	v1Group.Get(userPattern+userCodePattern, authHandler.ValidateToken)

}
