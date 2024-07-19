package node

import (
	"bytes"
	"cmp"
	"slices"
	"sync"

	"pan/app/cache"
)

type AppContext = *Context
type AppHandleFunc = HandleFunc[AppContext]
type AppHandleChain = HandleChain[AppContext]

func NewAppContext() AppContext {
	ctx := &Context{}
	InitContext(ctx)
	return ctx
}

type AppHandleGroup interface {
	Use(...AppHandleFunc) AppHandleGroup
	Handle(RequestName, ...AppHandleFunc) AppHandleGroup
	Default(...AppHandleFunc) AppHandleGroup
	Group() AppHandleGroup
	Route(RequestName) AppHandleGroup
}

type AppHandleChainItem[T any] struct {
	handles AppHandleChain
	code    T
}

func (item *AppHandleChainItem[T]) HashCode() T {
	return item.code
}

type AppSequence = uint16

type App struct {
	*AppRouter
	routes    cache.Bucket[RequestName, *AppHandleChainItem[RequestName]]
	defaults  AppHandleChain
	defaults_ cache.Bucket[AppSequence, *AppHandleChainItem[AppSequence]]
	rw        sync.RWMutex
	seq_      AppSequence
}

func NewApp() *App {
	app := &App{}
	app.AppRouter = &AppRouter{}
	app.AppRouter.app = app
	app.AppRouter.seq = 0

	routes := cache.NewBucket[RequestName, *AppHandleChainItem[RequestName]](bytes.Compare)
	app.routes = cache.WrapSyncBucket(routes)

	defaults_ := cache.NewBucket[AppSequence, *AppHandleChainItem[AppSequence]](cmp.Compare[AppSequence])
	app.defaults_ = cache.WrapSyncBucket(defaults_)

	return app
}

func (app *App) route(name RequestName, handles AppHandleChain) {
	route := &AppHandleChainItem[RequestName]{}
	route.code = name
	route.handles = append(handles, app.dispatchDefaults)
	err := app.routes.Store(route)
	if err != nil {
		panic(err)
	}
}

func (app *App) setDefaults(seq AppSequence, handles AppHandleChain) {
	if len(handles) > 0 {
		defaultItem := &AppHandleChainItem[AppSequence]{}
		defaultItem.handles = handles
		defaultItem.code = seq
		app.defaults_.Swap(defaultItem)
		return
	}

	defaultItem, ok := app.defaults_.Search(seq)
	if ok {
		app.defaults_.Delete(defaultItem)
	}
}

func (app *App) newRoute(name RequestName) *AppRouter {
	app.seq_++
	return NewAppRouter(app, name)
}

func (app *App) newGroup() AppHandleGroup {
	return app.newRoute(nil)
}

func (app *App) Init(extreme bool) {
	app.rw.Lock()
	defer app.rw.Unlock()

	items := app.defaults_.Items()
	for _, item := range items {
		app.defaults = slices.Concat(app.defaults, item.handles)
	}

	if extreme {
		app.AppRouter = nil
		app.defaults_ = nil
	}
}

func (app *App) Run(ctx AppContext, next Next) error {

	name := ctx.Name()
	route, ok := app.routes.Search(name)
	if ok {
		return Dispatch(ctx, route.handles, 0, next)
	}

	return app.dispatchDefaults(ctx, next)
}

func (app *App) dispatchDefaults(ctx AppContext, next Next) error {

	app.rw.RLock()
	if len(app.defaults) <= 0 {
		app.rw.RUnlock()
		return next()
	}
	handles := slices.Clone(app.defaults)
	app.rw.RUnlock()

	return Dispatch(ctx, handles, 0, next)
}

type AppRouter struct {
	seq         AppSequence
	name        RequestName
	app         *App
	middlewares AppHandleChain
}

func NewAppRouter(app *App, name RequestName) *AppRouter {
	router := &AppRouter{}
	router.seq = app.seq_
	router.name = name
	router.app = app
	return router
}

func (router *AppRouter) Use(handles ...AppHandleFunc) AppHandleGroup {
	router.middlewares = append(router.middlewares, handles...)
	return returnAppHandleGroup(router)
}

func (router *AppRouter) Handle(name RequestName, handles ...AppHandleFunc) AppHandleGroup {

	name_ := name
	if len(router.name) > 0 {
		name_ = slices.Concat(router.name, name)
	}
	handles_ := slices.Concat(router.middlewares, handles)
	app := router.app
	app.route(name_, handles_)

	return returnAppHandleGroup(router)
}

func (router *AppRouter) Default(handles ...AppHandleFunc) AppHandleGroup {

	app := router.app
	var defaults AppHandleChain

	if len(handles) > 0 {
		handles_ := slices.Concat(router.middlewares, handles)
		defaults = handles_
		if len(router.name) > 0 {
			handle := routeHandle(router.name, handles_)
			defaults = AppHandleChain{handle}
		}
	}
	app.setDefaults(router.seq, defaults)

	return returnAppHandleGroup(router)
}

func (router *AppRouter) Group() AppHandleGroup {
	return router.app.newGroup()
}

func (router *AppRouter) Route(name RequestName) AppHandleGroup {
	return router.app.newRoute(name)
}

func returnAppHandleGroup(router *AppRouter) AppHandleGroup {
	if router.seq > 0 {
		return router
	}
	return router.app
}

func routeHandle(name RequestName, handles AppHandleChain) AppHandleFunc {
	return func(ctx AppContext, next Next) error {
		if bytes.HasPrefix(ctx.Name(), name) {
			return Dispatch(ctx, handles, 0, next)
		}
		return next()
	}
}
