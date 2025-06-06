package bootstrap

import (
	"context"
	"errors"
	"pan/lib/injection"
	"pan/lib/runtime"
	"reflect"
	"sync"
)

var ErrBootstrapDestroyModuleUnavailable = errors.New("bootstrap.DestroyModule Error: Unavailable")

type DestroyModule interface {
	Destroy()
}

type destroyEngine struct {
	registry runtime.Registry
	locker   sync.RWMutex
}

var _ = (runtime.InitializeModule)((*destroyEngine)(nil))

func (de *destroyEngine) Init(registry runtime.Registry) error {
	de.locker.Lock()
	de.registry = registry
	de.locker.Unlock()
	return nil
}

var _ = (runtime.EngineExtensionModule)((*destroyEngine)(nil))

func (de *destroyEngine) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[DestroyModule](),
	}
}

var _ = (injection.ComponentProvider)((*destroyEngine)(nil))

func (de *destroyEngine) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(de, injection.ComponentExternalScope),
	}
}

func (de *destroyEngine) bootstrap(ctx context.Context) error {
	de.locker.RLock()
	registry := de.registry
	de.locker.RUnlock()
	if registry == nil {
		return ErrBootstrapDestroyModuleUnavailable
	}

	return runtime.TraverseRegistry(registry, func(module DestroyModule) error {
		module.Destroy()
		return nil
	})
}
