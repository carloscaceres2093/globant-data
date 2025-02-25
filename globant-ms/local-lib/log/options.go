package log

type opts struct {
	fields         map[string]interface{}
	contextTraceID string
}

type Options = func(*opts)

func Field(key string, value interface{}) Options {
	return func(o *opts) {
		o.fields[key] = value
	}
}

func WithContextTraceID(id string) Options {
	return func(o *opts) {
		o.contextTraceID = id
	}
}

func defaultOps() opts {
	return opts{
		fields:         map[string]interface{}{},
		contextTraceID: "",
	}
}

func applyOps(options ...Options) opts {
	o := defaultOps()
	for _, opt := range options {
		opt(&o)
	}

	return o
}
