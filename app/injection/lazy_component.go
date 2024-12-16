package injection

import "reflect"

type LazyComponentFunc[T any] func() T

type lazyComponentImpl[T any] struct {
	componentBase
	lazyFunc LazyComponentFunc[T]
}

func NewLazyComponent[T any](lazyFunc LazyComponentFunc[T], scope string) Component {
	base := componentBase{
		ty:    reflect.TypeFor[T](),
		scope: scope,
	}
	return &lazyComponentImpl[T]{
		lazyFunc:      lazyFunc,
		componentBase: base,
	}
}

func (c *lazyComponentImpl[T]) Target() interface{} {
	return c.lazyFunc()
}
