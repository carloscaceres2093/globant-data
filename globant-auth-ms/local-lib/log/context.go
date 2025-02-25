package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerCtxKey struct{}

func Context(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

func FromContext(ctx context.Context) *logrus.Logger {
	logger, _ := ctx.Value(loggerCtxKey{}).(*logrus.Logger)
	if logger == nil {
		return logrus.StandardLogger()
	}

	return logger
}
