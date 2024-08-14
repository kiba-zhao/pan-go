package net

import "pan/runtime"

func New() interface{} {
	return runtime.NewModule(&webServer{}, &broadcast{}, &quicModule{})
}
