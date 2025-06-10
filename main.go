package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"os"
	"pan/features/app"
	"pan/features/extfs"
	"pan/lib/config"
	"pan/lib/env"
	"pan/lib/log"
	"pan/lib/repository"
	"pan/lib/web"

	"pan/lib/runtime"
)

//go:generate go run script/build_web.go
//go:embed web/dist
var embedFS embed.FS // Declare embedded file system for web assets.

func main() {

	logger := log.Default()
	logger.Debug("main", "begin")
	defer logger.Debug("main", "end")

	// add web assets
	assetsFS, err := fs.Sub(embedFS, "web/dist")
	if err != nil {
		logger.Error("main", "Web Assets Error:"+err.Error())
		return
	}

	if prepare(logger) != nil {
		return
	}

	// web module
	webModule := web.New()
	webAssetsModule := web.NewWebAssets("/", assetsFS)

	// init and bootstrap
	engine, err := runtime.New(
		app.New(webModule),
		extfs.New(),
		webAssetsModule,
		app.Bootstrap(),
	)
	if err == nil {
		ctx := runtime.NewContext()
		err = engine.Bootstrap(ctx)
	}

	if err != nil && !runtime.IsAbort(err) {
		logger.Error("main", "Run Error:"+err.Error())
	}
}

func init() {
	// init logger
	logger := newLogger()
	log.InitDefault(logger)

	// enable web mode
	web.SetWebModeEnabled()
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

func prepare(logger log.Logger) error {
	// init as defaults  for config
	err := config.InitAsDefaults()
	if err != nil {
		logger.Error("main", "config Error:"+err.Error())
		return err
	}

	// init as defaults for repository
	err = repository.InitAsDefaults()
	if err != nil {
		logger.Error("main", "repository Error:"+err.Error())
		return err
	}
	return err
}
