package log

const (
	errorTag = "error"
)

type option struct {
	fields map[string]interface{}
}

type Option = func(*option)

func Field(key string, value interface{}) Option {
	return func(o *option) {
		o.fields[key] = value
	}
}

func Err(err error) Option {
	return Field(errorTag, err)
}
