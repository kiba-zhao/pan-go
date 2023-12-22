package core

import "bytes"

type Next = func() error

type Handle[T Context] func(ctx T, next Next) error

type Context interface {
	Method() []byte
}

type Handler[T Context] interface {
	Handle(ctx T, next Next) error
}

type handlerStruct[T Context] struct {
	method []byte
	handle Handle[T]
}

// Handle ...
func (h *handlerStruct[T]) Handle(ctx T, next Next) error {
	if h.method != nil && bytes.Equal(ctx.Method(), h.method) == false {
		return next()
	}

	return h.handle(ctx, next)
}
