package bootstrap

import (
	"context"
	"pan/app/constant"
	"pan/runtime"
	"reflect"
	"sync"
)

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

func (re *readyEngine) Components() []Component {
	return []Component{
		NewComponent(re, ComponentExternalScope),
	}
}

func (re *readyEngine) bootstrap(ctx context.Context) error {
	re.locker.RLock()
	registry := re.registry
	re.locker.RUnlock()
	if registry == nil {
		return constant.ErrUnavailable
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
