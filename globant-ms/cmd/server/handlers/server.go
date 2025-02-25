package server

import (
	"fmt"
	"os"

	apiHandler "globant-ms/cmd/server/handlers/handler"
	healthhandler "globant-ms/cmd/server/handlers/health"
	"globant-ms/internal/health"
	"globant-ms/internal/platform/config"
	apiService "globant-ms/internal/service"
	"globant-ms/local-lib/web"

	"github.com/labstack/gommon/log"
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
	postgres, err := apiService.NewPostgres(cfg.Database)
	if err != nil {
		log.Error("error db connection ", err)
		os.Exit(ExitCodeFailToCreateDBConnection)
	}
	service := apiService.NewService(postgres)

	return &services{
		healthService: healthService,
		service:       service,
	}, nil
}
