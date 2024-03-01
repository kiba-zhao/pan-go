package main

import (
	"pan/app"
	"pan/core"
)

func mountApp(coreApp *core.App) {
	m := app.NewModule()
	coreApp.Mount(m)
}

func main() {

	coreApp := core.New()
	err := coreApp.Init()

	if err == nil {
		mountApp(coreApp)
		err = coreApp.Run()
	}

	if err != nil {
		panic(err)
	}
}
