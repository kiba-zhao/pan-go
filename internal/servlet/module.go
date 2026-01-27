package servlet

import (
	"context"
	libApp "pan/internal/app"
	"pan/internal/bootstrap"
	"pan/internal/runtime"
	"sync"
)

type ServletModule interface {
	SetupToServlet(app ServletApp) error
}

func New() interface{} {
	module := &stdServletModule{}

	servlet := &stdServlet{}
	module.servlet = servlet

	return module
}

type stdServletModule struct {
	registry runtime.Registry
	rw       sync.RWMutex
	already  bool

	servlet *stdServlet
}

var _ = (runtime.InitializeModule)((*stdServletModule)(nil))

func (module *stdServletModule) Init(ctx context.Context, registry runtime.Registry) error {
	module.rw.Lock()
	module.registry = registry
	module.rw.Unlock()

	if module.already {
		return module.ReloadModules(ctx)
	}
	return nil
}

var _ = (bootstrap.DeferModule)((*stdServletModule)(nil))

func (module *stdServletModule) Defer(ctx context.Context) error {
	return module.ReloadModules(ctx)
}

func (module *stdServletModule) ReloadModules(ctx context.Context) error {
	module.rw.RLock()
	registry := module.registry
	module.rw.RUnlock()

	servletApp := libApp.NewApp()

	err := runtime.TraverseRegistry(registry, func(module ServletModule) error {
		if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
			return ctxErr
		}
		return module.SetupToServlet(servletApp)
	})

	if err == nil {
		setApp(module.servlet, servletApp)
	}
	return err
}
