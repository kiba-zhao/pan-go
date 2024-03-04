package main

import (
	"pan/app"
	"pan/core"
	"pan/extfs"
)

func mountApp(coreApp *core.App) {
	m := app.NewModule()
	coreApp.Mount(m)
}

func mountExtFS(coreApp *core.App) {
	m := extfs.NewModule()
	coreApp.Mount(m)
}

func main() {

	coreApp := core.New()
	err := coreApp.Init()

	if err == nil {
		mountApp(coreApp)
		mountExtFS(coreApp)
		err = coreApp.Run()
	}

	if err != nil {
		panic(err)
	}
}
