package main

import (
	"context"
	"embed"
	"io/fs"
	"pan/features/app"
	"pan/features/extfs"
	"pan/lib/log"
	"pan/lib/web"

	"pan/lib/runtime"
)

//go:generate go run script/build_web.go
//go:embed gui/web/dist
var embedFS embed.FS // Declare embedded file system for web assets.

func main() {

	// add web assets
	assetsFS, err := fs.Sub(embedFS, "gui/web/dist")
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
		log.Default().Log(context.Background(), log.LevelError, err.Error())
	}
}
