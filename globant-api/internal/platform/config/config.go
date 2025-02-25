package config

import (
	"errors"
)

const (
	ExitCodeFailTopicBrokerInfoInvalid = iota
	ExitCodeFailTopicNameInfoInvalid
	environment         = "ENVIRONMENT"
	microServicePattern = "/globant-api"
	environmentLocal    = "local"
	environmentProd     = "prod"
)

var ()

var _configs = map[string]Config{
	environmentLocal: {
		AppName:    microServicePattern,
		Env:        environment,
		UploadFile: "http://localhost:8082/globant-ms/v1/upload",
		Auth:       "http://localhost:8083/globant-auth-ms/v1/user/%s?token=%s",
	},
	environmentProd: {
		AppName:    microServicePattern,
		Env:        environment,
		UploadFile: "http://globant-ms:8080/globant-ms/v1/upload",
		Auth:       "http://globant-auth-ms:8080/globant-auth-ms/v1/user/%s?token=%s",
	},
}

func GetConfigFromEnvironment(env string) (Config, error) {
	config, found := _configs[env]
	if !found {
		return Config{}, errors.New("config not found for indicated environment")
	}
	return config, nil
}
