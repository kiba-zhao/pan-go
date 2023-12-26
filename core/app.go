package core

import (
	"sync"
)

type App[T Context] interface {
	Use(handlers ...Handler[T])
	UseFn(method []byte, handles ...Handle[T])
	Run(ctx T) error
}

type appStruct[T Context] struct {
	handlers []Handler[T]
	rw       *sync.RWMutex
}

// Use ...
func (app *appStruct[T]) Use(handlers ...Handler[T]) {
	app.rw.Lock()
	app.handlers = append(app.handlers, handlers...)
	app.rw.Unlock()
}

// UseFn ...
func (app *appStruct[T]) UseFn(method []byte, handles ...Handle[T]) {
	handlers := make([]Handler[T], 0)
	for _, handle := range handles {
		hanlder := new(handlerStruct[T])
		hanlder.method = method
		hanlder.handle = handle
		handlers = append(handlers, hanlder)
	}
	app.Use(handlers...)
}

// Run ...
func (app *appStruct[T]) Run(ctx T) error {

	app.rw.RLock()
	handlers := app.handlers[:]
	app.rw.RUnlock()

	return Dispatch(ctx, handlers, 0, nil)
}

// New ...
func New[T Context]() App[T] {
	app := new(appStruct[T])
	app.handlers = make([]Handler[T], 0)
	app.rw = new(sync.RWMutex)
	return app
}

func Dispatch[T Context, H Handler[T]](ctx T, handlers []H, index int, next Next) error {

	if index >= len(handlers) {
		if next != nil {
			return next()
		}
		return nil
	}

	handler := handlers[index]
	return handler.Handle(ctx, func() error {
		return Dispatch(ctx, handlers, index+1, next)
	})
}
