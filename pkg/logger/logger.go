package logger

import (
	"fmt"

	"github.com/gookit/slog"
)

var defaultLogger *slog.Record

func SetDefaultLogger(l *slog.Record) {
	defaultLogger = l
}

func Debug(i ...interface{}) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Debug(i...)
}

func Debugf(msg string, i ...interface{}) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Debugf(msg, i...)
}

func Info(i ...interface{}) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Info(i...)
}

func Infof(msg string, i ...interface{}) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Infof(msg, i...)
}

func Warn(i ...interface{}) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Warn(i...)
}

func Warnf(msg string, i ...interface{}) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Warnf(msg, i...)
}

func Error(i ...interface{}) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Error(i...)
}

func Errorf(msg string, i ...interface{}) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Errorf(msg, i...)
}

func Panic(i ...interface{}) {
	if defaultLogger == nil {
		panic(i)
	}
	defaultLogger.Panic(i...)
}

func Panicf(msg string, i ...interface{}) {
	if defaultLogger == nil {
		panic(fmt.Sprintf(msg, i...))
	}
	defaultLogger.Panicf(msg, i...)
}