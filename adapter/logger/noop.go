package logger

import "github.com/dityuiri/go-baseline/adapter/logger/log"

type noopLogger struct{}

func (noopLogger) Close() error { return nil }

func (noopLogger) Flush() error { return nil }

func (noopLogger) Debug(string, ...log.Option) {}

func (noopLogger) Info(string, ...log.Option) {}

func (noopLogger) Warn(string, ...log.Option) {}

func (noopLogger) Error(string, ...log.Option) {}

func (noopLogger) Panic(string, ...log.Option) {}

func (noopLogger) SetLevel(LevelType) {}

func (noopLogger) GetSkip() *int {
	skip := 0
	return &skip
}
