// Package web provides a web application framework for Pan.
//
// It is based on the Gin web framework.
// It is used to define web functionality in the application

package web

import (
	"context"
	"errors"
	"pan/internal/bootstrap"
	"pan/internal/config"
	"pan/internal/injection"
	"pan/internal/log"
	"pan/internal/runtime"
	"reflect"
	"sync"
)

var ErrWebModuleUnavailable = errors.New("web.WebModule Error: Unavailable")
var ErrWebModuleInvalidApp = errors.New("web.WebModule Error: Invalid App")

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

type WebController interface {
	SetupToWeb(WebRouter) error
}

type WebControllerProvider interface {
	WebControllers() []WebController
}

type WebRouteModule interface {
	WebRouteName() string
}

func New() interface{} {

	server := &stdWebServer{}
	server.reloadChan = make(chan struct{}, 1)

	logger := log.Default()
	server.logger = logger

	configurer := config.NewConfigurer[WebConfig](logger)

	wm := &stdWebModule{}
	wm.server = server
	wm.configurer = configurer

	return wm
}

type stdWebModule struct {
	server     *stdWebServer
	configurer WebConfigurer

	registry runtime.Registry
	rw       sync.RWMutex
	already  bool
}

var _ = (WebConfigListener)((*stdWebModule)(nil))

// OnConfigUpdated updates the web module configuration.
//
// It is called when the configuration of the engine is updated.
//
// It checks if the web address is changed, and if so, updates the address and
// triggers a reload by sending a signal to the reload channel.
func (w *stdWebModule) OnConfigUpdated(cfg WebConfig) {
	w.server.Setup(cfg)
}

var _ = (runtime.InitializeModule)((*stdWebModule)(nil))

// Init initializes the web module with the provided registry.
// It sets the module's registry and then attempts to reload modules.
// Returns an error if the module reloading fails.

func (w *stdWebModule) Init(ctx context.Context, registry runtime.Registry) error {
	w.rw.Lock()
	w.registry = registry
	w.rw.Unlock()

	if w.already {
		return w.ReloadModules(ctx, registry)
	}
	return nil
}

var _ = (bootstrap.DeferModule)((*stdWebModule)(nil))

// Defer triggers the reloading of modules for the web module.
//
// It calls ReloadModules, which reloads the web module's components.
// Returns an error if reloading fails.

func (w *stdWebModule) Defer(ctx context.Context) error {
	if w.already {
		return nil
	}

	w.already = true
	w.rw.RLock()
	registry := w.registry
	w.rw.RUnlock()

	return w.ReloadModules(ctx, registry)
}

var _ = (bootstrap.ReadyModule)((*stdWebModule)(nil))

func (w *stdWebModule) Ready(ctx context.Context) error {
	w.configurer.Subscribe(w)
	defer w.configurer.Unsubscribe(w)

	return w.server.ListenAndServe(ctx)
}

var _ = (runtime.EngineExtensionModule)((*stdWebModule)(nil))

// EngineTypes returns a slice of types for the web module's engine types.
//
// These are:
//
//   - WebAppModule
func (w *stdWebModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[WebAppModule](),
		reflect.TypeFor[WebControllerProvider](),
	}
}

func (w *stdWebModule) ReloadModules(ctx context.Context, registry runtime.Registry) error {
	app := NewWebApp()

	err := runtime.TraverseRegistry(registry, func(module WebAppModule) error {
		if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
			return ctxErr
		}
		return module.SetupToWeb(app)
	})

	if err == nil {
		err = runtime.TraverseRegistry(registry, func(module WebControllerProvider) error {
			if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
				return ctxErr
			}
			controllers := module.WebControllers()
			if len(controllers) <= 0 {
				return nil
			}
			var router WebRouter
			router = app
			if module, ok := module.(WebRouteModule); ok {
				routeName := module.WebRouteName()
				if len(routeName) > 0 {
					router = app.Group(routeName)
				}
			}

			for _, controller := range controllers {
				err := controller.SetupToWeb(router)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err == nil {
		w.server.SetupApp(app)
	}

	return err
}

var _ = (injection.ComponentProvider)((*stdWebModule)(nil))

func (w *stdWebModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(w, injection.ComponentNoneScope),
	}
}
