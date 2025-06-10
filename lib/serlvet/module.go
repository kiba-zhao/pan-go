package serlvet

import (
	"context"
	libApp "pan/lib/app"
	"pan/lib/bootstrap"
	"pan/lib/runtime"
	"sync"
)

type SerlvetModule interface {
	SetupToSerlvet(app SerlvetApp) error
}

func New() interface{} {
	module := &stdSerlvetModule{}

	serlvet := &stdSerlvet{}
	module.serlvet = serlvet

	return module
}

type stdSerlvetModule struct {
	registry runtime.Registry
	rw       sync.RWMutex
	already  bool

	serlvet *stdSerlvet
}

var _ = (runtime.InitializeModule)((*stdSerlvetModule)(nil))

func (module *stdSerlvetModule) Init(ctx context.Context, registry runtime.Registry) error {
	module.rw.Lock()
	module.registry = registry
	module.rw.Unlock()

	if module.already {
		return module.ReloadModules(ctx)
	}
	return nil
}

var _ = (bootstrap.DeferModule)((*stdSerlvetModule)(nil))

func (module *stdSerlvetModule) Defer(ctx context.Context) error {
	return module.ReloadModules(ctx)
}

func (module *stdSerlvetModule) ReloadModules(ctx context.Context) error {
	module.rw.RLock()
	registry := module.registry
	module.rw.RUnlock()

	serlvetApp := libApp.NewApp()

	err := runtime.TraverseRegistry(registry, func(module SerlvetModule) error {
		if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
			return ctxErr
		}
		return module.SetupToSerlvet(serlvetApp)
	})

	if err == nil {
		setApp(module.serlvet, serlvetApp)
	}
	return err
}
