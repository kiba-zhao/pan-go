//go:build !(android || ios)

package gobase

import (
	"log/slog"
	"os"
	"pan/pkg/log"
)

func init() {
	logger := newLogger()
	log.InitDefault(logger)
}

func newLogger() log.Logger {
	var logLevel slog.Level

	switch Mode() {
	case DebugMode:
		logLevel = slog.LevelDebug
	case TestMode:
		logLevel = slog.LevelDebug
	default:
		logLevel = slog.LevelInfo
	}

	slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	return log.NewStdLogger(slogger, "%s %s")
}
