package main

import (
	"embed"
	"io/fs"
	"pan/gobase"
	"pan/lib/log"
	"pan/lib/web"

	"pan/lib/runtime"
)

//go:generate go run script/build_web.go
//go:embed web/dist
var embedFS embed.FS // Declare embedded file system for web assets.

func main() {

	gobase.Init()

	logger := log.Default()
	logger.Debug("main", "begin")
	defer logger.Debug("main", "end")

	// add web assets
	assetsFS, err := fs.Sub(embedFS, "web/dist")
	if err != nil {
		logger.Error("main", "Web Assets Error:"+err.Error())
		return
	}

	engine, err := gobase.NewEngine(
		web.New(),
		web.NewWebAssets("/", assetsFS),
		gobase.New(),
	)
	if err == nil {
		ctx := runtime.NewContext()
		err = engine.Bootstrap(ctx)
	}

	if err != nil && !runtime.IsAbort(err) {
		logger.Error("main", "Run Error:"+err.Error())
	}
}
