package kitutil

import (
	"sync/atomic"

	kitlog "github.com/go-kit/kit/log"
	kitloglevel "github.com/go-kit/kit/log/level"
)

type Logger struct {
	Logger     kitlog.Logger
	ErrorCount int64
}

func (l *Logger) Log(keyvals ...interface{}) error {
	return l.Logger.Log(keyvals...)
}
func (l *Logger) Debug(keyvals ...interface{}) {
	l.countError(kitloglevel.Debug(l.Logger).Log(keyvals...))
}
func (l *Logger) Info(keyvals ...interface{}) {
	l.countError(kitloglevel.Info(l.Logger).Log(keyvals...))
}
func (l *Logger) Warn(keyvals ...interface{}) {
	l.countError(kitloglevel.Warn(l.Logger).Log(keyvals...))
}
func (l *Logger) Error(keyvals ...interface{}) {
	l.countError(kitloglevel.Error(l.Logger).Log(keyvals...))
}
func (l *Logger) With(keyvals ...interface{}) *Logger {
	return &Logger{Logger: kitlog.With(l.Logger, keyvals...)}
}
func (l *Logger) WithPrefix(keyvals ...interface{}) *Logger {
	return &Logger{Logger: kitlog.WithPrefix(l.Logger, keyvals...)}
}
func (l *Logger) countError(err error) {
	if err != nil {
		atomic.AddInt64(&l.ErrorCount, 1)
	}
}
