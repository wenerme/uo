package kitutil

import (
	"fmt"
	"os"
	"sync/atomic"

	kitlog "github.com/go-kit/kit/log"
	kitloglevel "github.com/go-kit/kit/log/level"
)

type loggerWrapper struct {
	Logger     kitlog.Logger
	ErrorCount int64
}

func NewLogger(logger kitlog.Logger) Logger {
	if l, ok := logger.(*loggerWrapper); ok {
		return l
	}
	if logger == nil {
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}
	return &loggerWrapper{Logger: logger}
}

func (l *loggerWrapper) Log(keyvals ...interface{}) error {
	return l.Logger.Log(keyvals...)
}
func (l *loggerWrapper) Debug(keyvals ...interface{}) {
	l.countError(kitloglevel.Debug(l.Logger).Log(keyvals...))
}
func (l *loggerWrapper) Info(keyvals ...interface{}) {
	l.countError(kitloglevel.Info(l.Logger).Log(keyvals...))
}
func (l *loggerWrapper) Warn(keyvals ...interface{}) {
	l.countError(kitloglevel.Warn(l.Logger).Log(keyvals...))
}
func (l *loggerWrapper) Error(keyvals ...interface{}) {
	l.countError(kitloglevel.Error(l.Logger).Log(keyvals...))
}
func (l *loggerWrapper) With(keyvals ...interface{}) Logger {
	return &loggerWrapper{Logger: kitlog.With(l.Logger, keyvals...)}
}
func (l *loggerWrapper) WithPrefix(keyvals ...interface{}) Logger {
	return &loggerWrapper{Logger: kitlog.WithPrefix(l.Logger, keyvals...)}
}
func (l *loggerWrapper) countError(err error) {
	if err != nil {
		atomic.AddInt64(&l.ErrorCount, 1)
		fmt.Println("Logger Error", err)
	}
}

type Logger interface {
	Log(keyvals ...interface{}) error
	Debug(keyvals ...interface{})
	Info(keyvals ...interface{})
	Warn(keyvals ...interface{})
	Error(keyvals ...interface{})
	With(keyvals ...interface{}) Logger
	WithPrefix(keyvals ...interface{}) Logger
}
