package log

import (
	"log/slog"
	"os"
	"sync"

	"pan/lib/env"
)

var logger Logger
var loggerOnce sync.Once

func Default() Logger {

	loggerOnce.Do(func() {
		logger = newLogger()
	})

	return logger
}

func newLogger() Logger {
	var logLevel slog.Level

	switch env.Mode() {
	case env.DebugMode:
		logLevel = slog.LevelDebug
	case env.TestMode:
		logLevel = slog.LevelDebug
	default:
		logLevel = slog.LevelInfo
	}

	slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	return NewStdLogger(slogger, "%s %s")
}
