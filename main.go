//go:build !(android || ios)

package main

import (
	"pan/gobase"
	"pan/pkg/log"

	"pan/pkg/runtime"
)

func main() {

	logger := log.Default()
	logger.Debug("main", "begin")
	defer logger.Debug("main", "end")

	cfg := gobase.NewSettingsConfig(log.Default())
	module := gobase.New(
		gobase.NewHostModule(),
		gobase.NewSettingsModule(cfg),
	)

	engine, err := runtime.New(module)
	if err == nil {
		ctx := runtime.NewContext()
		err = engine.Bootstrap(ctx)
	}

	if err != nil && !runtime.IsAbort(err) {
		logger.Error("main", "Run Error:"+err.Error())
	}
}
