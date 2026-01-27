package injection

import "reflect"

type LazyComponentFunc[T any] func() T

type stdLazyComponent[T any] struct {
	componentBase
	lazyFunc LazyComponentFunc[T]
}

// NewLazyComponent returns a new Component with a lazy loading function.
// The function lazyFunc is called when the component is first accessed.
// The scope of the component is determined by the scope parameter.
// If the scope is empty, the component is external.
// If the scope is "none", the component is not registered in the
// injection system.
// If the scope is "internal", the component is registered in the
// injection system and can be injected into other components.
func NewLazyComponent[T any](lazyFunc LazyComponentFunc[T], scope string) Component {
	base := componentBase{
		ty:    reflect.TypeFor[T](),
		scope: scope,
	}
	return &stdLazyComponent[T]{
		lazyFunc:      lazyFunc,
		componentBase: base,
	}
}

func (c *stdLazyComponent[T]) Target() interface{} {
	return c.lazyFunc()
}
