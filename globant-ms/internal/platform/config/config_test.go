package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm/logger"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		env            string
		expectedConfig Config
		expectedError  error
	}{
		{
			name: "Ok - GetConfigFromScope for local scope ",
			env:  "local",
			expectedConfig: Config{
				Database: Database{
					User:     "postgres",
					Password: "postgres",
					Host:     "127.0.0.1",
					Port:     "5432",
					Name:     "postgres",
					LogLevel: logger.Info,
				},
				Env:        "ENVIRONMENT",
				AppName:    "/globant-ms",
				UploadFile: "test",
			},
		},
		{
			name:          "Error - GetConfigFromScope for unrecognized scope ",
			env:           "unrecognized",
			expectedError: errors.New("config not found for indicated environment"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := GetConfigFromEnvironment(tt.env)
			require.Equal(t, tt.expectedError, err)
			require.Equal(t, tt.expectedConfig, config)
		})
	}
}
