package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
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
				Env:        "ENVIRONMENT",
				AppName:    "/globant-api",
				UploadFile: "http://localhost:8082/globant-ms/v1/upload",
				Auth:       "http://localhost:8083/globant-auth-ms/v1/user/%s?token=%s",
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
