package log

import (
	"fmt"
	"log/slog"
)

type Logger interface {
	Info(tag string, msg string)
	Warn(tag string, msg string)
	Error(tag string, msg string)
	Debug(tag string, msg string)
}

type StdLogger struct {
	logger *slog.Logger
	logFmt string
}

var _ = (Logger)((*StdLogger)(nil))

func NewStdLogger(logger *slog.Logger, logFmt string) *StdLogger {
	stdLogger := &StdLogger{}
	stdLogger.logger = logger
	stdLogger.logFmt = logFmt
	return stdLogger
}

func (l *StdLogger) Info(tag string, msg string) {
	l.logger.Info(fmt.Sprintf(l.logFmt, tag, msg))
}

func (l *StdLogger) Warn(tag string, msg string) {
	l.logger.Warn(fmt.Sprintf(l.logFmt, tag, msg))
}

func (l *StdLogger) Error(tag string, msg string) {
	l.logger.Error(fmt.Sprintf(l.logFmt, tag, msg))
}

func (l *StdLogger) Debug(tag string, msg string) {
	l.logger.Debug(fmt.Sprintf(l.logFmt, tag, msg))
}
