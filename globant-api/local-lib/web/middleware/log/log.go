package log

import (
	"context"

	"globant-api/local-lib/log"
)

func Info(ctx context.Context, step string, options ...Option) {
	fields := mapOptions(options...)



	FromContext(ctx).Info(ctx, step, fields...)
}

func Error(ctx context.Context, step string, options ...Option) {
	fields := mapOptions(options...)
	FromContext(ctx).Error(ctx,  step, fields...)
}

func Warning(ctx context.Context, step string, options ...Option) {
	fields := mapOptions(options...)

	FromContext(ctx).Warning(ctx, step, fields...)
}

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
