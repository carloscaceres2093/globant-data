package main

import (
	"context"
	"fmt"
	"os"

	"globant-ms/internal/platform/config"
	"globant-ms/local-lib/web/middleware/log"

	server "globant-ms/cmd/server/handlers"
)

const (
	exitCodeFailToCreateWebApplication = iota
	ExitCodeFailInitReportLib
	exitCodeFailReadConfigs
	environment = "ENVIRONMENT"
	loadTag     = "load"
)

func main() {
	env := os.Getenv(environment)
	cfg, err := config.GetConfigFromEnvironment(env)
	if err != nil {

		log.Error(context.Background(), "read_config", log.Field(loadTag,
			fmt.Sprintf("main: can't read config: %v", err.Error())))
		os.Exit(exitCodeFailReadConfigs)
	}

	err = server.StartServer(cfg)
	if err != nil {
		log.Error(context.Background(), "start_server", log.Field(loadTag,
			fmt.Sprintf("main: error creating services: %v", err.Error())))
		os.Exit(exitCodeFailToCreateWebApplication)
	}
}
