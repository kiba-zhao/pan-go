// Package web provides a web application framework for Pan.
//
// It is based on the Gin web framework.
// It is used to define web functionality in the application

package web

import (
	"context"
	"errors"
	"net/http"
	"pan/app/config"
	"pan/logger"
	"pan/runtime"
	"reflect"
	"slices"
	"sync"

	"github.com/gin-gonic/gin"
)

var ErrWebModuleUnavailable = errors.New("web.WebModule Error: Unavailable")
var ErrWebModuleInvalidApp = errors.New("web.WebModule Error: Invalid App")

type WebApp = *gin.Engine
type WebRouter = gin.IRouter
type WebContext = *gin.Context

// NewWebApp creates a new web application.
//
// It creates a new Gin Engine which can be used as a WebApp.
func NewWebApp() WebApp {
	return gin.New()
}

// WebAppModule is a module for the web application.
//
// The module that implements this interface will be obtained by the application from the runtime and loaded into the gin engine
type WebAppModule interface {

	// SetupToWeb sets up the web application
	//
	// Example:
	// type MyModule struct {}
	//
	// func (m *MyModule) SetupToWeb(app WebApp) error {
	// app.GET("/my-module", m.MyMethod)
	// app.NoRoute(m.noRoute)
	// ...
	// return nil
	// }
	SetupToWeb(WebApp) error
}

// WebAppModuleProvider is a module that provides multiple WebAppModules
type WebAppModuleProvider interface {
	WebAppModules() []WebAppModule
}

// WebController is a controller for the web application
//
// The module that implements this interface will be obtained by the application from the runtime and loaded into the gin engine
// It is used to define web api in the application
type WebController interface {
	// SetupToWeb sets up the controller for the web application
	//
	// Example:
	// type MyController struct {}
	//
	// func (c *MyController) SetupToWeb(r WebRouter) error {
	// 	r.GET("/my-controller", c.MyMethod)
	//  ...
	// 	return nil
	// }
	//

	SetupToWeb(WebRouter) error
}

// WebScopeModule is a module that provides a scope for the web application
type WebScopeModule interface {
	WebScope() string
}

// WebControllerProvider is a module that provides multiple WebControllers
type WebControllerProvider interface {
	WebControllers() []WebController
}

// New creates a new runtime engine module.
//
// It returns a *webModule, which implements runtime.Module.
func New() interface{} {
	return &webModule{}
}

type webModule struct {
	app        WebApp
	appLocker  sync.RWMutex
	registry   runtime.Registry
	locker     sync.RWMutex
	addresses  []string
	reloadChan chan struct{}
	reloadOnce sync.Once
	needReload bool
}

// ReloadChan returns a channel that receives a signal when the web module needs to be reloaded.
//
// It is used by the http server to reload the web module.
func (w *webModule) ReloadChan() chan struct{} {

	w.reloadOnce.Do(func() {
		w.reloadChan = make(chan struct{}, 1)
	})

	return w.reloadChan
}

// OnConfigUpdated updates the web module configuration.
//
// It is called when the configuration of the engine is updated.
//
// It checks if the web address is changed, and if so, updates the address and
// triggers a reload by sending a signal to the reload channel.
func (w *webModule) OnConfigUpdated(settings config.AppSettings) {
	w.locker.Lock()
	defer w.locker.Unlock()

	if slices.Equal(w.addresses, settings.WebAddress) {
		return
	}

	w.addresses = settings.WebAddress

	// trigger to reload
	if w.needReload {
		return
	}
	w.needReload = true
	w.ReloadChan() <- struct{}{}
}

// Init initializes the web module with the provided registry.
// It sets the module's registry and then attempts to reload modules.
// Returns an error if the module reloading fails.

func (w *webModule) Init(registry runtime.Registry) error {
	w.locker.Lock()
	w.registry = registry
	w.locker.Unlock()

	return w.ReloadModules()
}

// Defer triggers the reloading of modules for the web module.
//
// It calls ReloadModules, which reloads the web module's components.
// Returns an error if reloading fails.

func (w *webModule) Defer() error {
	return w.ReloadModules()
}

// EngineTypes returns a slice of types for the web module's engine types.
//
// These are:
//
//   - WebAppModule
//   - WebAppModuleProvider
//   - WebControllerProvider
func (w *webModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[WebAppModule](),
		reflect.TypeFor[WebAppModuleProvider](),
		reflect.TypeFor[WebControllerProvider](),
	}
}

// ServeHTTP implements the http.Handler interface.
//
// It serves the web application's http requests by calling the ServeHTTP method
// on the underlying Gin Engine.
//
// If the web application is not available (i.e., the app is nil), it returns an
// HTTP error with the status code 500.
func (w *webModule) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	w.appLocker.RLock()
	app := w.app
	w.appLocker.RUnlock()

	if app != nil {
		app.ServeHTTP(rw, req)
	} else {
		http.Error(rw, ErrWebModuleInvalidApp.Error(), http.StatusInternalServerError)
	}
}

// ReloadModules reloads the web application modules from the registry.
//
// It traverses the registry, looking for modules that implement the WebAppModule,
// WebAppModuleProvider, and WebControllerProvider interfaces. For each module,
// it calls the SetupToWeb method to set up the web application.
//
// If any error occurs during the reloading process, it will be returned.
//
// If the reloading process is successful, the web application will be updated
// with the new modules.
func (w *webModule) ReloadModules() error {
	w.locker.Lock()
	registry := w.registry
	w.locker.Unlock()
	if registry == nil {
		return ErrWebModuleUnavailable
	}

	w.appLocker.Lock()
	defer w.appLocker.Unlock()

	app := NewWebApp()

	err := runtime.TraverseRegistry(registry, func(module WebAppModule) error {
		return module.SetupToWeb(app)
	})
	if err == nil {
		err = runtime.TraverseRegistry(registry, func(module WebAppModuleProvider) error {
			for _, wc := range module.WebAppModules() {
				perr := wc.SetupToWeb(app)
				if perr != nil {
					return perr
				}
			}
			return nil
		})
	}
	if err == nil {
		err = runtime.TraverseRegistry(registry, func(module WebControllerProvider) error {
			var router WebRouter
			if scopeModule, ok := module.(WebScopeModule); ok {
				scope := scopeModule.WebScope()
				router = app.Group(scope)
			} else {
				router = app
			}
			for _, wc := range module.WebControllers() {
				perr := wc.SetupToWeb(router)
				if perr != nil {
					return perr
				}
			}
			return nil
		})
	}

	if err == nil {
		w.app = app
	}

	return err
}

// Ready starts the web server and waits for it to exit.
//
// It takes a context, and listens for the Done signal. When the context is done,
// it shuts down the web server and waits for it to exit.
//
// It also listens for the reload signal, which triggers a reload of the web
// application modules from the registry.
//
// If the reloading process fails, it returns the error. Otherwise, it starts the
// web server and waits for it to exit.
//
// If the web server exits with an error, it logs the error and returns it.
func (w *webModule) Ready(ctx context.Context) error {

	var wg sync.WaitGroup
	var servers []*http.Server
	var err error
	closed := false
	for {

		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-w.ReloadChan():
		}
		w.locker.Lock()
		w.needReload = false
		addresses := w.addresses
		w.locker.Unlock()

		if len(servers) > 0 {
			for _, server := range servers {
				server.Shutdown(context.Background())
			}
			wg.Wait()
		}

		if closed {
			break
		}

		for _, address := range addresses {
			httpServer := &http.Server{
				Addr:    address,
				Handler: w,
			}

			servers = append(servers, httpServer)
			wg.Add(1)
			go func(s *http.Server) {
				defer wg.Done()
				err = s.ListenAndServe()
				if err != nil {
					logger.Default().Log(context.Background(), logger.LevelError, "app.web.webModule.Ready Error: %s", err.Error())
				}
				// TODO: echo error into log
			}(httpServer)
		}

	}

	return err
}
