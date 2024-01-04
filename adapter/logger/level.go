package logger

import (
	"reflect"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LevelType
type LevelType struct {
	string
}

type level struct {
	DEBUG LevelType
	INFO  LevelType
	WARN  LevelType
	ERROR LevelType
	PANIC LevelType
}

// Level
var Level = &level{
	DEBUG: LevelType{"DEBUG"},
	INFO:  LevelType{"INFO"},
	WARN:  LevelType{"WARN"},
	ERROR: LevelType{"ERROR"},
	PANIC: LevelType{"PANIC"},
}

func (level LevelType) ToAtomicLevel() zap.AtomicLevel {
	switch level {
	case Level.DEBUG:
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case Level.INFO:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case Level.WARN:
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case Level.ERROR:
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case Level.PANIC:
		return zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}
	return zap.NewAtomicLevelAt(zapcore.DebugLevel)
}

func (level LevelType) ToLevel() zapcore.Level {
	switch level {
	case Level.DEBUG:
		return zap.DebugLevel
	case Level.INFO:
		return zap.InfoLevel
	case Level.WARN:
		return zap.WarnLevel
	case Level.ERROR:
		return zap.ErrorLevel
	case Level.PANIC:
		return zap.PanicLevel
	}
	return zap.DebugLevel
}

func StringToLevel(level string) LevelType {
	value := reflect.ValueOf(*Level)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Interface().(LevelType).string == strings.ToUpper(level) {
			return field.Interface().(LevelType)
		}
	}
	return Level.DEBUG
}
