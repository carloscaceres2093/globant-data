package log

import (
	"context"
	"net/http"

	"globant-api/local-lib/log"
)

type loggerCtxKey struct{}

// Logging provides a middleware that introduces a logger in the context, to allows
func Logging(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			ctx := Context(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}


func Context(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}


func FromContext(ctx context.Context) log.Logger {
	logger, ok := ctx.Value(loggerCtxKey{}).(log.Logger)

	if !ok {
		return log.NoOp{}
	}
	return logger
}
