package main

import (
	"embed"
	"io/fs"
	"pan/app"
	"pan/app/net"
	"pan/extfs"
	"pan/runtime"
)

//go:generate npm --prefix ./web install
//go:generate npm --prefix ./web run build -- -m production
//go:embed web/dist
var embedFS embed.FS

func main() {

	assetsFS, err := fs.Sub(embedFS, "web/dist")
	if err != nil {
		panic(err)
	}

	engine := runtime.New()

	err = engine.Mount(app.New(), extfs.New(), net.NewWebAssets("/", assetsFS))

	var ctx runtime.Context
	if err == nil {
		ctx, err = engine.Bootstrap()
	}

	if err == nil {
		err = ctx.Wait()
	}

	if err != nil {
		panic(err)
	}
}
