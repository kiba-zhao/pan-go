// Package web provides a web application framework for Pan.
//
// It is based on the Gin web framework.
// It is used to define web functionality in the application

package web

import (
	"context"
	"errors"
	"pan/lib/bootstrap"
	"pan/lib/config"
	"pan/lib/injection"
	"pan/lib/log"
	"pan/lib/runtime"
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

func New() interface{} {
	wm := &stdWebModule{}

	server := &stdWebServer{}
	wm.server = server
	server.logger = log.Default()
	server.reloadChan = make(chan struct{}, 1)

	return wm
}

type stdWebModule struct {
	AppConfig config.AppConfig
	server    *stdWebServer
	registry  runtime.Registry
	rw        sync.RWMutex
	already   bool
}

var _ = (config.AppConfigListener)((*stdWebModule)(nil))

// OnConfigUpdated updates the web module configuration.
//
// It is called when the configuration of the engine is updated.
//
// It checks if the web address is changed, and if so, updates the address and
// triggers a reload by sending a signal to the reload channel.
func (w *stdWebModule) OnConfigUpdated(settings config.AppSettings) {
	w.server.SetAddrs(settings.WebAddress)
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
		return w.ReloadModules(ctx)
	}
	return nil
}

var _ = (bootstrap.DeferModule)((*stdWebModule)(nil))

// Defer triggers the reloading of modules for the web module.
//
// It calls ReloadModules, which reloads the web module's components.
// Returns an error if reloading fails.

func (w *stdWebModule) Defer(ctx context.Context) error {
	w.already = true
	return w.ReloadModules(ctx)
}

var _ = (bootstrap.ReadyModule)((*stdWebModule)(nil))

func (w *stdWebModule) Ready(ctx context.Context) error {
	w.AppConfig.Subscribe(w)
	defer w.AppConfig.Unsubscribe(w)

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
	}
}

func (w *stdWebModule) ReloadModules(ctx context.Context) error {
	w.rw.RLock()
	registry := w.registry
	w.rw.RUnlock()
	if registry == nil {
		return ErrWebModuleUnavailable
	}

	app := NewWebApp()

	err := runtime.TraverseRegistry(registry, func(module WebAppModule) error {
		if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
			return ctxErr
		}
		return module.SetupToWeb(app)
	})

	if err == nil {
		w.server.SetWebApp(app)
	}

	return err
}

var _ = (injection.ComponentProvider)((*stdWebModule)(nil))

func (w *stdWebModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(w, injection.ComponentNoneScope),
	}
}
