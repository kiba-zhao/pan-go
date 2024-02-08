package main

import "pan/desktop"

func init() {
	println("Hello, Init!")
}
func main() {
	println("Hello, World!")
	dt := desktop.New()
	defer dt.Destroy()
	dt.Show()
}
