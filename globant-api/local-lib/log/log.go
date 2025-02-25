package log

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	contextTraceTag   = "context_trace_id"
	stepTag           = "step"
	fileTag           = "file_name"
	functionTag       = "function_name"
	attributesTag     = "attributes"
)

const (
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel Level = iota

	// WarningLevel logs are more important than Info, but don't need individual
	// human review.
	WarningLevel

	// InfoLevel is the default logging priority.
	InfoLevel

	// DebugLevel logs are typically voluminous, and are usually disabled in production.
	DebugLevel
)

// A Level is a logging priority. Higher levels are more important.
type Level uint32

var levels = map[Level]logrus.Level{
	ErrorLevel:   logrus.ErrorLevel,
	WarningLevel: logrus.WarnLevel,
	InfoLevel:    logrus.InfoLevel,
	DebugLevel:   logrus.DebugLevel,
}

// Fields is an alias for logrus.Fields. Aliasing this type dramatically
// improves the navigability of this package's API documentation.
type Fields = logrus.Fields

// Entry is an alias for logrus.Entry.
type Entry = logrus.Entry

// Logger is an interface that defines the methods that a logger must implement.
type Logger interface {
	Info(ctx context.Context,  step string, options ...Options)
	Error(ctx context.Context,  step string, options ...Options)
	Warning(ctx context.Context,  step string, options ...Options)
	Debug(ctx context.Context,  step string, options ...Options)
}

type logger interface {
	WithFields(fields Fields) *Entry
}

// Log is a Logger. Provides all functionalities related to log applications events.
type Log struct {
	log      logger
	level    Level
}

// NewLogger initializes and retrieves a Log.
// This uses a logrus Logger by default, and if we want to use a custom logger, we can pass it as a parameter.
func NewLogger(cf ...Config) Log {
	config := applyConfig(cf...)
	logger := config.logger
	if logger == nil {
		logger = newLogrusLogger(levels[config.level])
	}
	return Log{
		log:      logger,
		level:    config.level,
	}
}

// Info sends a logs with InfoLevel level.
func (l Log) Info(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(ctx,  step, file, function, o)

	l.log.WithFields(fields).Info()

}

// Error sends a logs with ErrorLevel level.
func (l Log) Error(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(ctx,  step, file, function, o)
	l.log.WithFields(fields).Error()

}

// Warning sends a logs with WarningLevel level.
func (l Log) Warning(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(ctx,  step, file, function, o)
	l.log.WithFields(fields).Warning()

}

// Debug sends a logs with DebugLevel level.
func (l Log) Debug(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(ctx,  step, file, function, o)
	l.log.WithFields(fields).Debug()

}

func mapFields(ctx context.Context,  step, file, function string, opts opts) Fields {
	f := Fields{}

	f[stepTag] = step
	f[fileTag] = file
	f[functionTag] = function

	attributesFields := logrus.Fields{}
	for key, value := range opts.fields {
		attributesFields[key] = value
		// Otherwise errors are ignored by encoding/json
		// https://github.com/sirupsen/logrus/issues/137
		if err, ok := value.(error); ok {
			attributesFields[key] = err.Error()
		}
	}

	f[attributesTag] = attributesFields

	if opts.contextTraceID != "" {
		f[contextTraceTag] = opts.contextTraceID
	}

	return f
}

func newLogrusLogger(level logrus.Level) *logrus.Logger {
	logrusLog := logrus.New()
	logrusLog.SetLevel(level)
	logrusLog.SetReportCaller(false)

	logrusLog.SetFormatter(&JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	return logrusLog
}

func getCaller() (string, string) {
	const (
		maxCallersStackGet = 5
		// skip the call to runtime.CallersFrames(pcs), the call to this function
		// and the before call (Info, Error, Debug, Warning...)
		skippedCallers = 3

		unknownCaller           = "unknown"
		middlewareLogCallerFile = "log.go"
	)
	var (
		file     = unknownCaller
		function = unknownCaller
	)

	pcs := make([]uintptr, maxCallersStackGet)
	n := runtime.Callers(skippedCallers, pcs)
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		frameFile := path.Base(frame.File)
		line := frame.Line
		file = fmt.Sprintf("%s:%d", frameFile, line)

		lastSlash := strings.LastIndexByte(frame.Function, '/')
		function = frame.Function[lastSlash+1:]

		if !strings.Contains(file, middlewareLogCallerFile) {
			break
		}
	}

	return file, function
}
