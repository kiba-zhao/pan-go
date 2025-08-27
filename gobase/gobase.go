package gobase

import (
	"log/slog"
	"os"
	"pan/lib/bootstrap"
	"pan/lib/broadcast"
	"pan/lib/env"
	"pan/lib/log"
	"pan/lib/peer"
	"pan/lib/quic"
	"pan/lib/repository"
	"pan/lib/runtime"
)

func Init() {
	// init logger
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

func NewEngine(specifyModules ...interface{}) (*runtime.Engine, error) {

	// base modules
	modules := []interface{}{
		bootstrap.New(),
		repository.New(),
		peer.New(),
		broadcast.New(),
		quic.New(),
	}
	//

	// specific modules for platform
	if len(specifyModules) > 0 {
		modules = append(modules, specifyModules...)
	}
	//

	// base features

	//

	// bootstrap module
	modules = append(modules, bootstrap.Bootstrap())
	//

	return runtime.New(modules)
}
