package middleware

import (
	"globant-ms/local-lib/log"
	"net/http"

	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := log.Context(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
