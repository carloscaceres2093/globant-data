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
	traceTag          = "trace_id"
	contextTraceTag   = "context_trace_id"
	stepTag           = "step"
	fileTag           = "file_name"
	functionTag       = "function_name"
	attributesTag     = "attributes"
	datadogTag        = "dd"
	datadogTraceIDTag = "trace_id"
	datadogSpanIDTag  = "span_id"
)

const (
	ErrorLevel Level = iota
	WarningLevel
	InfoLevel
	DebugLevel
)

type Level uint32

var levels = map[Level]logrus.Level{
	ErrorLevel:   logrus.ErrorLevel,
	WarningLevel: logrus.WarnLevel,
	InfoLevel:    logrus.InfoLevel,
	DebugLevel:   logrus.DebugLevel,
}
type Fields = logrus.Fields

type Entry = logrus.Entry

type Logger interface {
	Info(ctx context.Context,  step string, options ...Options)
	Error(ctx context.Context,  step string, options ...Options)
	Warning(ctx context.Context,  step string, options ...Options)
	Debug(ctx context.Context,  step string, options ...Options)
}

type logger interface {
	WithFields(fields Fields) *Entry
}

type Log struct {
	log      logger
	level    Level
}

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

func (l Log) Info(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(ctx,  step, file, function, o)

	l.log.WithFields(fields).Info()

}

func (l Log) Error(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(ctx,  step, file, function, o)
	l.log.WithFields(fields).Error()

}


func (l Log) Warning(ctx context.Context,  step string, options ...Options) {
	o := applyOps(options...)
	file, function := getCaller()

	fields := mapFields(ctx,  step, file, function, o)
	l.log.WithFields(fields).Warning()

}

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
