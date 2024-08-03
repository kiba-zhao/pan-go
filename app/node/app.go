package node

import (
	"bytes"
	"cmp"
	"slices"
	"sync"
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

type AppSequence = uint16

type App struct {
	*AppRouter
	routes          []*AppHandleChainItem[RequestName]
	routesRW        sync.RWMutex
	defaults        AppHandleChain
	defaults_       []*AppHandleChainItem[AppSequence]
	defaultsLocker_ sync.Mutex
	rw              sync.RWMutex
	seq_            AppSequence
}

func NewApp() *App {
	app := &App{}
	app.AppRouter = &AppRouter{}
	app.AppRouter.app = app
	app.AppRouter.seq = 0

	app.routes = make([]*AppHandleChainItem[RequestName], 0)
	app.defaults_ = make([]*AppHandleChainItem[AppSequence], 0)

	return app
}

func (app *App) compareRoute(route *AppHandleChainItem[RequestName], name RequestName) int {
	return bytes.Compare(route.code, name)
}

func (app *App) route(name RequestName, handles AppHandleChain) {
	route := &AppHandleChainItem[RequestName]{}
	route.code = name
	route.handles = append(handles, app.dispatchDefaults)

	app.routesRW.Lock()
	defer app.routesRW.Unlock()
	idx, ok := slices.BinarySearchFunc(app.routes, name, app.compareRoute)
	if !ok {
		app.routes = slices.Insert(app.routes, idx, route)
	}
}

func (app *App) compareDefaults_(item *AppHandleChainItem[AppSequence], seq AppSequence) int {
	return cmp.Compare(item.code, seq)
}

func (app *App) setDefaults(seq AppSequence, handles AppHandleChain) {
	app.defaultsLocker_.Lock()
	defer app.defaultsLocker_.Unlock()
	idx, ok := slices.BinarySearchFunc(app.defaults_, seq, app.compareDefaults_)
	if len(handles) > 0 {
		defaultItem := &AppHandleChainItem[AppSequence]{}
		defaultItem.handles = handles
		defaultItem.code = seq
		if ok {
			app.defaults_[idx] = defaultItem
		} else {
			app.defaults_ = slices.Insert(app.defaults_, idx, defaultItem)
		}
		return
	}

	if ok {
		app.defaults_ = slices.Delete(app.defaults_, idx, idx+1)
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

	app.defaultsLocker_.Lock()
	defer app.defaultsLocker_.Unlock()

	app.rw.Lock()
	defer app.rw.Unlock()

	for _, item := range app.defaults_ {
		app.defaults = slices.Concat(app.defaults, item.handles)
	}

	if extreme {
		app.AppRouter = nil
		app.defaults_ = nil
	}
}

func (app *App) Run(ctx AppContext, next Next) error {

	name := ctx.Name()
	app.routesRW.RLock()
	idx, ok := slices.BinarySearchFunc(app.routes, name, app.compareRoute)
	if ok {
		defer app.routesRW.RUnlock()
		route := app.routes[idx]
		return Dispatch(ctx, route.handles, 0, next)
	}
	app.routesRW.RUnlock()

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
