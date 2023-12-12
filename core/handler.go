package core

import "bytes"

type Next = func()

type Handle[T Context] func(ctx T, next Next)

type Context interface {
	Method() []byte
}

type Handler[T Context] interface {
	Handle(ctx T, next Next)
}

type handlerStruct[T Context] struct {
	method []byte
	handle Handle[T]
}

// Handle ...
func (h *handlerStruct[T]) Handle(ctx T, next Next) {
	if h.method != nil && bytes.Equal(ctx.Method(), h.method) == false {
		next()
		return
	}

	h.handle(ctx, next)
}
