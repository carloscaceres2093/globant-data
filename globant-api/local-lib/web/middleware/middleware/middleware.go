package auth

import (
	"context"
	"globant-api/internal/platform/config"
	apiClient "globant-api/internal/service/client"
	"net/http"
	"os"
)

const (
	exitCodeFailToCreateWebApplication = iota
	ExitCodeFailInitReportLib
	exitCodeFailReadConfigs
	environment = "ENVIRONMENT"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := os.Getenv(environment)
		cfg, err := config.GetConfigFromEnvironment(env)
		if err != nil {
			os.Exit(exitCodeFailReadConfigs)
		}

		token := r.Header.Get("Authorization")
		userCode := r.Header.Get("X-user")
		if token == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}
		if userCode == "" {
			http.Error(w, "Unauthorized: Missing user", http.StatusUnauthorized)
			return
		}
		userAuth, err := apiClient.AuthValidation(userCode, token, cfg.Auth)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "X-user-code", userAuth.UserCode)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
