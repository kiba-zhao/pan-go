package main

import (
	"context"
	"embed"
	"io/fs"
	"pan/app"
	"pan/app/web"
	"pan/extfs"
	"pan/logger"

	"pan/runtime"
)

//go:generate npm --prefix ./web install
//go:generate npm --prefix ./web run build -- -m production
//go:embed web/dist
var embedFS embed.FS // Declare embedded file system for web assets.

func main() {

	// add web assets
	assetsFS, err := fs.Sub(embedFS, "web/dist")
	if err != nil {
		panic(err)
	}

	// init runtime
	engine := runtime.New()

	// mount modules
	err = engine.Mount(app.New(), extfs.New(), web.NewWebAssets("/", assetsFS), app.Bootstrap())

	if err == nil {
		err = engine.Bootstrap()
	}

	if err != nil {
		logger.Default().Log(context.Background(), logger.LevelError, err.Error())
	}
}
