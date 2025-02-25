package config

import (
	"errors"
	"globant-ms/local-lib/env"

	"gorm.io/gorm/logger"
)

const (
	ExitCodeFailTopicBrokerInfoInvalid = iota
	ExitCodeFailTopicNameInfoInvalid
	environment         = "ENVIRONMENT"
	microServicePattern = "/globant-ms"
	environmentLocal    = "local"
	environmentDev      = "dev"
	environmentSB       = "sb"
	environmentProd     = "prod"
	environmentStg      = "stg"
)

var (
	basicDb = Database{
		User:     env.GetEnv("DB_USER"),
		Password: env.GetEnv("DB_PASSWORD"),
		Host:     env.GetEnv("DB_HOST"),
		Port:     env.GetEnv("DB_PORT"),
		Name:     env.GetEnv("DB_NAME"),
		LogLevel: logger.Silent,
	}
)

var _configs = map[string]Config{
	environmentLocal: {
		Database: Database{
			User:     "admin",
			Password: "admin",
			Host:     "127.0.0.1",
			Port:     "5432",
			Name:     "globant-ms-test",
			LogLevel: logger.Info,
		},
		AppName:    microServicePattern,
		Env:        environment,
		UploadFile: "http://localhost:8084/data-report-ms/v1/reports",
	},
	environmentProd: {
		Database:   basicDb,
		AppName:    microServicePattern,
		Env:        environment,
		UploadFile: "https://internal-prod.y.uno/data-report-ms/v1/reports",
	},
}

func GetConfigFromEnvironment(env string) (Config, error) {
	config, found := _configs[env]
	if !found {
		return Config{}, errors.New("config not found for indicated environment")
	}
	return config, nil
}
