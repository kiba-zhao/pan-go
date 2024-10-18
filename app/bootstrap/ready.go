package bootstrap

import (
	"pan/app/constant"
	"pan/runtime"
	"reflect"
	"sync"

	"golang.org/x/sync/errgroup"
)

type ReadyModule interface {
	Ready() error
}

type readyEngine struct {
	registry runtime.Registry
	locker   sync.RWMutex
}

func (re *readyEngine) Init(registry runtime.Registry) error {
	re.locker.Lock()
	re.registry = registry
	re.locker.Unlock()
	return nil
}

func (re *readyEngine) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[ReadyModule](),
	}
}

func (re *readyEngine) Components() []Component {
	return []Component{
		NewComponent(re, ComponentExternalScope),
	}
}

func (re *readyEngine) bootstrap() error {
	re.locker.RLock()
	registry := re.registry
	re.locker.RUnlock()
	if registry == nil {
		return constant.ErrUnavailable
	}

	var ctx errgroup.Group
	runtime.TraverseRegistry(registry, func(module ReadyModule) error {
		ctx.Go(module.Ready)
		return nil
	})
	return ctx.Wait()
}
