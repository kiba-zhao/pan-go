package main

import (
	"embed"
	"errors"
	"io/fs"
	"pan/features/app"
	"pan/features/extfs"
	"pan/lib/bootstrap"
	"pan/lib/log"
	"pan/lib/web"

	"pan/lib/runtime"
)

//go:generate go run script/build_web.go
//go:embed web/dist
var embedFS embed.FS // Declare embedded file system for web assets.

func main() {

	// init logger
	logger := log.Default()

	logger.Debug("main", "begin")
	defer logger.Debug("main", "end")

	// add web assets
	assetsFS, err := fs.Sub(embedFS, "web/dist")
	if err != nil {
		logger.Error("main", "Web Assets Error:"+err.Error())
		return
	}

	// init runtime
	engine := runtime.New()

	// mount modules
	err = engine.Mount(app.New(), extfs.New(), web.NewWebAssets("/", assetsFS), app.Bootstrap())

	if err == nil {
		err = engine.Bootstrap()
	}

	if err != nil && !errors.Is(err, bootstrap.ErrBootstrapExit) {
		logger.Error("main", "Run Error:"+err.Error())
	}
}
