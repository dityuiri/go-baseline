package logger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/dityuiri/go-baseline/adapter/logger/log"
)

func NewLogger(opts ...Option) (ILogger, error) {
	// o := MergeOptions(opts...)

	options := &Options{}

	for _, opt := range opts {
		opt(options)
	}

	if options.isNoop {
		return &noopLogger{}, nil
	}

	config := &zap.Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},

		DisableCaller: true, // use custom caller info

		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			StacktraceKey: "stacktrace",
			LevelKey:      "level",
			TimeKey:       "time",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
		},

		Level: zap.NewAtomicLevelAt(zapcore.DebugLevel),
	}

	if options.level != nil {
		config.Level = options.level.ToAtomicLevel()
	}

	if options.output != nil {
		config.Encoding = options.output.ToEncoding()
	}

	if zapLogger, err := config.Build(); err != nil {
		return nil, err
	} else {
		return &logger{
			name:   options.name,
			logger: zapLogger,
			config: config,

			skip: options.skip,
		}, nil
	}
}

func (l *logger) Flush() error {
	if r := recover(); r != nil {
		l.Error("panic", log.WithError(fmt.Errorf("%v", r)))
	}

	return l.logger.Sync()
}
func (l *logger) Close() error {
	return l.Flush()
}

func (l *logger) Debug(message string, options ...log.Option) {
	l.logger.Debug(message, l.compose(options...)...)
}

func (l *logger) Info(message string, options ...log.Option) {
	l.logger.Info(message, l.compose(options...)...)
}

func (l *logger) Warn(message string, options ...log.Option) {
	l.logger.Warn(message, l.compose(options...)...)
}

func (l *logger) Error(message string, options ...log.Option) {
	l.logger.Error(message, l.compose(options...)...)
}

func (l *logger) Panic(message string, options ...log.Option) {
	l.logger.Panic(message, l.compose(options...)...)
}

func (l *logger) SetLevel(level LevelType) {
	l.config.Level.SetLevel(level.ToLevel())
}

func (l *logger) GetSkip() *int {
	return l.skip
}

func (l *logger) compose(options ...log.Option) []zapcore.Field {
	var fields []zapcore.Field

	if l.name != nil {
		fields = append(fields, zap.String("app", *l.name))
	}

	opts := &log.Options{}

	for _, opt := range options {
		opt(opts)
	}

	if data := opts.GetData(); data != nil {
		fields = append(fields, zap.String("data", *data))
	}
	if err := opts.GetError(); err != nil {
		fields = append(fields, zap.Error(*err))
	}

	var skip int
	if s := opts.GetSkip(); s != nil {
		skip = *s
	} else if l.skip != nil {
		skip = *l.skip
	}

	if ci, err := retrieveCallInfo(skip); err == nil {
		fields = append(fields, []zap.Field{
			zap.String("package", ci.Package),
			zap.String("function", ci.Function),
			zap.String("file", ci.File),
			zap.Int("line", ci.Line),
		}...)
	}

	return fields
}

type callInfo struct {
	Package  string
	Function string
	File     string
	Line     int
}

func retrieveCallInfo(skip int) (*callInfo, error) {
	skip = 3 + skip // omit stacks of logger library call
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil, errors.New("failed to get call info")
	}
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &callInfo{
		Package:  packageName,
		Function: funcName,
		File:     fileName,
		Line:     line,
	}, nil
}
