package servlet

import (
	"bytes"
	"slices"
)

type HandleContext interface {
	Name() RequestName
}

type HandleGroup[T HandleContext] interface {
	// Use registers a middleware function for the router.
	Use(...HandleFunc[T]) HandleGroup[T]
	// Handle registers a handle function for a specific route name.
	Handle(RequestName, ...HandleFunc[T]) HandleGroup[T]
	// Default registers a default handle function for the router.
	Default(...HandleFunc[T]) HandleGroup[T]
	// Group creates a new HandleGroup
	Group() HandleGroup[T]
	// Route creates a new HandleGroup
	Route(RequestName) HandleGroup[T]
}

type Router[T HandleContext] struct {
	seq         Sequence
	name        RequestName
	servlet     *Servlet[T]
	middlewares HandleChain[T]
}

func NewRouter[T HandleContext](servlet *Servlet[T], name RequestName) *Router[T] {
	router := &Router[T]{}
	router.seq = servlet.seq_
	router.name = name
	router.servlet = servlet
	return router
}

func (router *Router[T]) Use(handles ...HandleFunc[T]) HandleGroup[T] {
	router.middlewares = append(router.middlewares, handles...)
	return returnHandleGroup(router)
}

func (router *Router[T]) Handle(name RequestName, handles ...HandleFunc[T]) HandleGroup[T] {

	name_ := GenerateRouteName(router.name, name)
	handles_ := slices.Concat(router.middlewares, handles)
	servlet := router.servlet
	servlet.route(name_, handles_)

	return returnHandleGroup(router)
}

func (router *Router[T]) Default(handles ...HandleFunc[T]) HandleGroup[T] {

	servlet := router.servlet
	var defaults HandleChain[T]

	if len(handles) > 0 {
		handles_ := slices.Concat(router.middlewares, handles)
		defaults = handles_
		if len(router.name) > 0 {
			handle := routeHandle(router.name, handles_)
			defaults = HandleChain[T]{handle}
		}
	}
	servlet.setDefaults(router.seq, defaults)

	return returnHandleGroup(router)
}

func (router *Router[T]) Group() HandleGroup[T] {
	return router.servlet.newGroup()
}

func (router *Router[T]) Route(name RequestName) HandleGroup[T] {
	if len(name) > 0 {
		return router.servlet.newRoute(name)
	}
	return router
}

func returnHandleGroup[T HandleContext](router *Router[T]) HandleGroup[T] {
	if router.seq > 0 {
		return router
	}
	return router.servlet
}

func routeHandle[T HandleContext](name RequestName, handles HandleChain[T]) HandleFunc[T] {
	return func(ctx T, next Next) error {
		if bytes.HasPrefix(ctx.Name(), name) {
			return Dispatch(ctx, handles, 0, next)
		}
		return next()
	}
}

var AppRouteSeparator = []byte(".")

func GenerateRouteName(scope RequestName, name RequestName) RequestName {
	if len(scope) > 0 {
		return bytes.Join([][]byte{scope, name}, AppRouteSeparator)
	}
	return name
}
