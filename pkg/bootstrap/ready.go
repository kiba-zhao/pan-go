package bootstrap

import (
	"context"
	"errors"
	"pan/pkg/injection"
	"pan/pkg/runtime"
	"reflect"
	"sync"
)

var ErrBootstrapReadyModuleUnavailable = errors.New("bootstrap.ReadyModule Error: Unavailable")

type ReadyModule interface {
	// Ready is called by the runtime to start the module after the application has finished initializing.
	//
	// Ready is called after the application has finished initializing and is used to start the module.
	// It is used to start the module.
	//
	// The function will be called with a context that is canceled when the application exits.
	// If the function returns an error, the error is logged and the application exits.
	Ready(context.Context) error
}

type readyEngine struct {
	registry runtime.Registry
	locker   sync.RWMutex
}

var _ = (runtime.InitializeModule)((*readyEngine)(nil))

// Init initializes the ready engine with the given registry.
//
// It sets the registry and does not return an error.
func (re *readyEngine) Init(ctx context.Context, registry runtime.Registry) error {
	re.locker.Lock()
	re.registry = registry
	re.locker.Unlock()
	return nil
}

var _ = (runtime.EngineExtensionModule)((*readyEngine)(nil))

// EngineTypes returns a slice of reflect.Type representing the various engine types
// associated with the ready engine. These types include:
//
//   - ReadyModule: Represents a module that can be started by the ready engine.
func (re *readyEngine) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[ReadyModule](),
	}
}

var _ = (injection.ComponentProvider)((*readyEngine)(nil))

// Components returns a slice of injection.Component representing the components
// provided by the ready engine. Currently, the only component provided is the
// readyEngine itself, which is scoped externally.
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
		if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
			return ctxErr
		}

		wg.Add(1)
		go func(readyModule ReadyModule) {
			defer wg.Done()
			err := readyModule.Ready(causeCtx)
			if err != nil {
				causeCancel(err)
			}
		}(module)

		return nil
	})

	if err == nil {
		wg.Wait()
		<-causeCtx.Done()
		err = causeCtx.Err()
	}

	return err
}
