package app

import (
	"context"
	"pan/pkg/bootstrap"
	"pan/pkg/runtime"
	"pan/pkg/servlet"
	"reflect"
	"sync"
)

type AppletModule interface {
	SetupToApplet(router AppServletRouter) error
}

type stdModule struct {
	registry runtime.Registry
	rw       sync.RWMutex
	already  bool

	applet *stdApplet
}

func New() interface{} {

	applet := &stdApplet{}

	module := &stdModule{}
	module.applet = applet

	return module
}

var _ = (runtime.EngineExtensionModule)((*stdModule)(nil))

func (module *stdModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[AppletModule](),
	}
}

var _ = (runtime.InitializeModule)((*stdModule)(nil))

func (module *stdModule) Init(ctx context.Context, registry runtime.Registry) error {
	module.rw.Lock()
	module.registry = registry
	module.rw.Unlock()

	if module.already {
		return module.ReloadModules(ctx, registry)
	}
	return nil
}

var _ = (bootstrap.DeferModule)((*stdModule)(nil))

func (module *stdModule) Defer(ctx context.Context) error {
	if module.already {
		return nil
	}
	module.already = true

	module.rw.RLock()
	registry := module.registry
	module.rw.RUnlock()

	return module.ReloadModules(ctx, registry)
}

func (module *stdModule) ReloadModules(ctx context.Context, registry runtime.Registry) error {

	applet := module.applet
	appServlet := servlet.NewServlet[AppServletContext]()

	err := runtime.TraverseRegistry(registry, func(module AppletModule) error {
		if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
			return ctxErr
		}
		return module.SetupToApplet(appServlet.Router)
	})
	if err == nil {
		applet.SetupAppServlet(appServlet)
	}
	return err
}
