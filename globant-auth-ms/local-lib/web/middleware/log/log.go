package log

import (
	"context"

	"globant-auth-ms/local-lib/log"
)

// Info logs a message at level Info on the logger retrieved by the context.
func Info(ctx context.Context, step string, options ...Option) {
	fields := mapOptions(options...)

	FromContext(ctx).Info(ctx, step, fields...)
}

// Error logs a message at level Error on the logger retrieved by the context.
func Error(ctx context.Context, step string, options ...Option) {
	fields := mapOptions(options...)
	FromContext(ctx).Error(ctx, step, fields...)
}

// Warning logs a message at level Warning on the logger retrieved by the context.
func Warning(ctx context.Context, step string, options ...Option) {
	fields := mapOptions(options...)

	FromContext(ctx).Warning(ctx, step, fields...)
}

// Debug logs a message at level Debug on the logger retrieved by the context.
func Debug(ctx context.Context, step string, options ...Option) {
	fields := mapOptions(options...)

	FromContext(ctx).Debug(ctx, step, fields...)
}

func mapOptions(options ...Option) []log.Options {
	var fields []log.Options

	option := option{
		fields: map[string]interface{}{},
	}
	for _, o := range options {
		o(&option)
	}

	for name, value := range option.fields {
		fields = append(fields, log.Field(name, value))
	}

	return fields
}
