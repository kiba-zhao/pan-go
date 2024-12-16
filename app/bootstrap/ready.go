package bootstrap

import (
	"context"
	"errors"
	"pan/app/injection"
	"pan/runtime"
	"reflect"
	"sync"
)

var ErrBootstrapReadyModuleUnavailable = errors.New("bootstrap.ReadyModule Error: Unavailable")

type ReadyModule interface {
	Ready(context.Context) error
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

func (re *readyEngine) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(re, injection.ComponentExternalScope),
	}
}

func (re *readyEngine) bootstrap(ctx context.Context) error {
	re.locker.RLock()
	registry := re.registry
	re.locker.RUnlock()
	if registry == nil {
		return ErrBootstrapReadyModuleUnavailable
	}

	var wg sync.WaitGroup
	causeCtx, causeCancel := context.WithCancelCause(ctx)
	err := runtime.TraverseRegistry(registry, func(module ReadyModule) error {
		wg.Add(1)
		go func(readyModule ReadyModule) {
			defer wg.Done()
			err := readyModule.Ready(causeCtx)
			if err != nil {
				causeCancel(err)
			}
		}(module)

		select {
		case <-causeCtx.Done():
			return causeCtx.Err()
		default:
			return nil
		}
	})

	if err == nil {
		wg.Wait()
		<-causeCtx.Done()
		err = causeCtx.Err()
	}

	return err
}
