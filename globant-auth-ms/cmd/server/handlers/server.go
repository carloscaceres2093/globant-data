package server

import (
	"fmt"
	"os"

	authHandler "globant-auth-ms/cmd/server/handlers/auth"
	healthhandler "globant-auth-ms/cmd/server/handlers/health"
	authService "globant-auth-ms/internal/auth"
	"globant-auth-ms/internal/health"
	"globant-auth-ms/internal/platform/config"
	"globant-auth-ms/local-lib/web"

	"github.com/labstack/gommon/log"
)

const (
	exitCodeFailToCreateWebApplication = iota
	ExitCodeFailToCreateDBConnection
	exitCodeFailReadConfigs
	environment = "ENVIRONMENT"
	zTAWS       = "us-west-2"
)

type authServices struct {
	healthService healthhandler.HealthService
	authService   authHandler.AuthService
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

func initServices(cfg config.Config) (*authServices, error) {
	healthService := health.NewService()
	postgres, err := authService.NewPostgres(cfg.Database)
	if err != nil {
		log.Error("error db connection ", err)
		os.Exit(ExitCodeFailToCreateDBConnection)
	}
	authService := authService.NewAuthService(postgres, cfg.Salt)

	return &authServices{
		healthService: healthService,
		authService:   authService,
	}, nil
}
