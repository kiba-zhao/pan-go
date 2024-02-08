package main

import "pan/app"

func main() {
	a := app.New()
	err := a.Init()
	if err != nil {
		panic(err)
	}
	a.Run()
}
