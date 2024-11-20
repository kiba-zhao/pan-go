package bootstrap

import (
	"context"
	"errors"
	"pan/runtime"
	"reflect"
	"sync"
)

var ErrBootstrapDeferModuleUnavailable = errors.New("bootstrap.DeferModule Error: Unavailable")

type DeferModule interface {
	Defer() error
}

type deferEngine struct {
	registry runtime.Registry
	locker   sync.RWMutex
}

func (de *deferEngine) Init(registry runtime.Registry) error {
	de.locker.Lock()
	de.registry = registry
	de.locker.Unlock()
	return nil
}

func (de *deferEngine) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[DeferModule](),
	}
}

func (de *deferEngine) Components() []Component {
	return []Component{
		NewComponent(de, ComponentExternalScope),
	}
}

func (de *deferEngine) bootstrap(ctx context.Context) error {
	de.locker.RLock()
	registry := de.registry
	de.locker.RUnlock()
	if registry == nil {
		return ErrBootstrapDeferModuleUnavailable
	}

	return runtime.TraverseRegistry(registry, func(module DeferModule) error {
		err := module.Defer()
		if err == nil {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			default:
			}
		}
		return err
	})
}
