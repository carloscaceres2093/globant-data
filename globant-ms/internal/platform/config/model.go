package config

import "gorm.io/gorm/logger"

type Config struct {
	Database   Database
	AppName    string
	Env        string
	UploadFile string
}

type Configs struct {
	Scope map[string]Config
}

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	LogLevel logger.LogLevel
}
