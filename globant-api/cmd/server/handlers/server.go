package server

import (
	"fmt"

	apiHandler "globant-api/cmd/server/handlers/handler"
	healthhandler "globant-api/cmd/server/handlers/health"
	"globant-api/internal/health"
	"globant-api/internal/platform/config"
	apiService "globant-api/internal/service"
	apiClient "globant-api/internal/service/client"
	"globant-api/local-lib/web"
)

const (
	exitCodeFailToCreateWebApplication = iota
	ExitCodeFailToCreateDBConnection
	exitCodeFailReadConfigs
	environment = "ENVIRONMENT"
	zTAWS       = "us-west-2"
)

type services struct {
	healthService healthhandler.HealthService
	service       apiHandler.Service
}

func StartServer(cfg config.Config) error {
	r := web.NewWebApp(cfg.Env)
	s, err := initServices(cfg)
	if err != nil {
		fmt.Print("error init services", err)
		return err
	}

	routes(r, s, cfg)

	err = r.Run()
	if err != nil {
		fmt.Print("error booting application", err)
	}
	return err
}

func initServices(cfg config.Config) (*services, error) {
	healthService := health.NewService()
	client := apiClient.NewClient(cfg)
	service := apiService.NewService(client)

	return &services{
		healthService: healthService,
		service:       service,
	}, nil
}
