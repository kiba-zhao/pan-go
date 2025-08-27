package gobase

func New() interface{} {
	module := &stdModule{}
	return module
}

type stdModule struct{}
