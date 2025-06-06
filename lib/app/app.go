// Define peer application
//
// It provides a simple and easy-to-use p2p application framework for Pan.
// Its implementation and usage logic are similar to gin
package app

import (
	"bytes"
	"cmp"
	"slices"
	"sync"
)

type AppContext = *Context
type AppHandleFunc = HandleFunc[AppContext]
type AppHandleChain = HandleChain[AppContext]

// NewAppContext returns a new initialized AppContext.
//
// It allocates a new context and calls InitContext to initialize it with default values.
func NewAppContext() AppContext {
	ctx := &Context{}
	InitContext(ctx)
	return ctx
}

type AppHandleGroup interface {
	// Use registers a middleware function for the router.
	Use(...AppHandleFunc) AppHandleGroup
	// Handle registers a handle function for a specific route name.
	Handle(RequestName, ...AppHandleFunc) AppHandleGroup
	// Default registers a default handle function for the router.
	Default(...AppHandleFunc) AppHandleGroup
	// Group creates a new AppHandleGroup
	Group() AppHandleGroup
	// Route creates a new AppHandleGroup
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

// NewApp returns a new initialized App.
//
// It allocates a new App and sets up the default values.
// The AppRouter is initialized with seq=0.
// The routes and defaults_ are initialized as empty slices.
func NewApp() *App {
	app := &App{}
	app.AppRouter = &AppRouter{}
	app.AppRouter.app = app
	app.AppRouter.seq = 0

	app.routes = make([]*AppHandleChainItem[RequestName], 0)
	app.defaults_ = make([]*AppHandleChainItem[AppSequence], 0)

	return app
}

// compareRoute compares the request name associated with a route to a given request name.
// It uses a byte comparison to determine the order.
//
// Parameters:
// - route: The route item containing a request name code to compare.
// - name: The request name to compare against the route's code.
//
// Returns an integer:
// - a negative number if the route's code is less than the name.
// - zero if they are equal.
// - a positive number if the route's code is greater than the name.

func (app *App) compareRoute(route *AppHandleChainItem[RequestName], name RequestName) int {
	return bytes.Compare(route.code, name)
}

// route registers a handle function for a specific route name.
//
// The handle function is added to the slice of handles for the route name.
// If the route name is not present in the app.routes slice, it is added.
// The dispatchDefaults handle is appended to the slice of handles as the last handle.
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

// compareDefaults_ compares the sequence code of a default handle chain item with a given sequence.
// It returns an integer indicating the result of the comparison:
// a negative number if the item's code is less than the sequence,
// zero if they are equal, and a positive number if the item's code is greater than the sequence.

func (app *App) compareDefaults_(item *AppHandleChainItem[AppSequence], seq AppSequence) int {
	return cmp.Compare(item.code, seq)
}

// setDefaults sets a default handle chain for a given sequence code.
//
// It appends the dispatchDefaults handle to the end of the handles slice.
// If the sequence code is not present in the app.defaults_ slice, it is added.
// If the sequence code is present and the handles slice is empty, the item is removed from the slice.
// If the sequence code is present and the handles slice is not empty, the item is updated in the slice.
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

// newRoute returns a new AppRouter.
//
// It increments the sequence number of the App, and allocates a new AppRouter with the given name.
// If the name is nil, it is set to an empty string.
func (app *App) newRoute(name RequestName) *AppRouter {
	app.seq_++
	return NewAppRouter(app, name)
}

// newGroup returns a new AppHandleGroup.
//
// It allocates a new AppRouter with name=nil and returns it as an AppHandleGroup.
func (app *App) newGroup() AppHandleGroup {
	return app.newRoute(nil)
}

// Init initializes the App by concatenating default handle chains and optionally
// clears the AppRouter and default handle chains if extreme is true.
//
// Locks are used to ensure thread-safety when modifying defaults and defaults_.
//
// Parameters:
// - extreme: A boolean flag that, when true, clears the AppRouter and defaults_.

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

// Run executes a handle chain registered for the given context name.
//
// It performs a binary search on the app.routes slice to find the index of the
// route with code matching the context name. If found, it executes the associated
// handles chain and returns the result. If not found, it calls the dispatchDefaults
// function with the context and next function, and returns the result.
//
// The function takes a context and a next function as parameters. The context is
// used to get the name for the route to search for. The next function is called
// if the end of the handle chain is reached, and its result is returned.
//
// Parameters:
// - ctx: The context of type AppContext to get the name from.
// - next: The next function to call if the end of the handle chain is reached.
//
// Returns an error if any handle function in the chain returns an error.
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

// dispatchDefaults executes the default handle chain.
//
// It locks the app's read-write mutex, checks if the app has any default handles,
// and if so, clones the handles and unlocks the mutex. It then calls the Dispatch
// function with the cloned handles, the context, and the next function.
//
// If the app has no default handles, or if the end of the handle chain is reached,
// it calls the next function and returns the result.
//
// Parameters:
// - ctx: The context of type AppContext to execute the handle chain with.
// - next: The next function to call if the end of the handle chain is reached.
//
// Returns an error if any handle function in the chain returns an error.
func (app *App) dispatchDefaults(ctx AppContext, next Next) error {

	app.rw.RLock()
	if len(app.defaults) > 0 {
		handles := slices.Clone(app.defaults)
		app.rw.RUnlock()
		return Dispatch(ctx, handles, 0, next)
	}
	app.rw.RUnlock()
	if next != nil {
		return next()
	}
	return nil

}

type AppRouter struct {
	seq         AppSequence
	name        RequestName
	app         *App
	middlewares AppHandleChain
}

// NewAppRouter creates a new instance of AppRouter.
//
// It initializes the AppRouter with the given App and RequestName,
// setting the sequence number from the App's current sequence.
// The AppRouter is returned with its name and app fields set.
//
// Parameters:
// - app: The application instance to associate with the new router.
// - name: The request name to be used for the router.
//
// Returns the newly initialized AppRouter.

func NewAppRouter(app *App, name RequestName) *AppRouter {
	router := &AppRouter{}
	router.seq = app.seq_
	router.name = name
	router.app = app
	return router
}

// Use registers a middleware function for the router.
//
// It appends the given middleware functions to the current chain of middleware functions
// for the router, and returns the AppHandleGroup instance.
//
// Parameters:
// - handles: The middleware functions to be appended to the current chain.
//
// Returns the AppHandleGroup instance.
func (router *AppRouter) Use(handles ...AppHandleFunc) AppHandleGroup {
	router.middlewares = append(router.middlewares, handles...)
	return returnAppHandleGroup(router)
}

// Handle registers a handle function for the specified request name.
//
// It generates a route name by combining the router's name with the provided request name,
// and concatenates the router's middlewares with the provided handle functions.
// The resulting handle chain is then registered with the application.
//
// Parameters:
// - name: The request name for which the handle functions are registered.
// - handles: The handle functions to be executed for the specified request name.
//
// Returns the AppHandleGroup instance.

func (router *AppRouter) Handle(name RequestName, handles ...AppHandleFunc) AppHandleGroup {

	name_ := GenerateRouteName(router.name, name)
	handles_ := slices.Concat(router.middlewares, handles)
	app := router.app
	app.route(name_, handles_)

	return returnAppHandleGroup(router)
}

// Default registers a default handle chain for the router.
//
// It takes the given middleware functions, and if the router's name is not empty,
// it generates a route name by concatenating the router's name with the default route name,
// and registers a handle chain for the generated route name.
// The handle chain is generated by concatenating the router's middlewares with the given middleware functions.
// The generated handle chain is then registered with the application.
//
// Parameters:
// - handles: The middleware functions to be used to generate the default handle chain.
//
// Returns the AppHandleGroup instance.
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

// Group creates a new AppHandleGroup.
//
// It allocates a new AppRouter with name=nil and returns it as an AppHandleGroup.
func (router *AppRouter) Group() AppHandleGroup {
	return router.app.newGroup()
}

// Route creates a new AppHandleGroup for the specified request name.
//
// If the provided request name is not empty, it allocates a new AppRouter
// with the given name and returns it as an AppHandleGroup. If the name is empty,
// it returns the current router instance.
//
// Parameters:
// - name: The request name to be used for the new route.
//
// Returns the AppHandleGroup instance.

func (router *AppRouter) Route(name RequestName) AppHandleGroup {
	if len(name) > 0 {
		return router.app.newRoute(name)
	}
	return router
}

// returnAppHandleGroup returns the AppHandleGroup instance based on the provided router.
//
// If the sequence number of the router is greater than 0, it returns the router itself.
// Otherwise, it returns the AppHandleGroup instance of the App associated with the router.
func returnAppHandleGroup(router *AppRouter) AppHandleGroup {
	if router.seq > 0 {
		return router
	}
	return router.app
}

// routeHandle generates a handle function for a given route name and handle chain.
//
// The generated handle function checks if the request name of the context starts with the given route name,
// and if so, it executes the handle chain. If the handle chain is empty, it calls the next function.
// If the context's request name does not match the route name, it calls the next function.
//
// Parameters:
// - name: The route name to check against the context's request name.
// - handles: The handle chain to execute if the request name matches the route name.
//
// Returns a handle function that can be registered with the application.
func routeHandle(name RequestName, handles AppHandleChain) AppHandleFunc {
	return func(ctx AppContext, next Next) error {
		if bytes.HasPrefix(ctx.Name(), name) {
			return Dispatch(ctx, handles, 0, next)
		}
		return next()
	}
}

// GenerateRouteName combines a scope and a request name into a single route name.
//
// If the scope is non-empty, it joins the scope and name using a byte slice,
// returning the combined route name. If the scope is empty, it returns the name as-is.
//
// Parameters:
// - scope: The scope to be combined with the request name.
// - name: The request name to be used.
//
// Returns the combined route name if the scope is non-empty, otherwise returns the name.

func GenerateRouteName(scope RequestName, name RequestName) RequestName {
	if len(scope) > 0 {
		return bytes.Join([][]byte{scope, name}, nil)
	}
	return name
}
