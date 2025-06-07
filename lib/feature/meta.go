package feature

import (
	"reflect"
)

type Metadata[T any] struct {
	metaType reflect.Type
	target   T
}

func newMetadata[T any, V any](target V) Metadata[V] {
	return Metadata[V]{
		metaType: reflect.TypeFor[T](),
		target:   target,
	}
}
