package logger

//go:generate mockgen -destination=mock/logger.go -package=mock . ILogger

import (
	"github.com/dityuiri/go-baseline/adapter/logger/log"
	"go.uber.org/zap"
)

type ILogger interface {
	Close() error
	Flush() error

	Debug(message string, options ...log.Option)
	Info(message string, options ...log.Option)
	Warn(message string, options ...log.Option)
	Error(message string, options ...log.Option)
	Panic(message string, options ...log.Option)

	SetLevel(level LevelType)
	// SetOutput(output OutputType)

	GetSkip() *int
}

type logger struct {
	name   *string
	logger *zap.Logger
	config *zap.Config

	skip *int
}
