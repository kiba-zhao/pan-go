//go:build !(android || ios)

package gobase

import (
	"log/slog"
	"os"
	"pan/pkg/env"
	"pan/pkg/log"
)

func init() {
	logger := newLogger()
	log.InitDefault(logger)
}

func newLogger() log.Logger {
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
	return log.NewStdLogger(slogger, "%s %s")
}
