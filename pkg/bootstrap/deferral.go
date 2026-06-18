// Define bootstrap deferral engine
//
// The deferral engine is used to defer the initialization of the application
package bootstrap

import (
	"context"
	"errors"
	"pan/pkg/injection"
	"pan/pkg/runtime"
	"reflect"
	"sync"
)

var ErrBootstrapDeferModuleUnavailable = errors.New("bootstrap.DeferModule Error: Unavailable")

type DeferModule interface {
	// Defer is called by the runtime to defer the initialization of the application
	//
	// Defer is called after the application has finished initializing.
	// It is used to defer the initialization of the application.
	//
	// The function will be called with a context that is canceled when the application exits.
	Defer(ctx context.Context) error
}

type deferEngine struct {
	registry runtime.Registry
	locker   sync.RWMutex
}

var _ = (runtime.InitializeModule)((*deferEngine)(nil))

// Init initializes the defer engine with the provided registry.
//
// It sets the registry and does not return an error.
func (de *deferEngine) Init(ctx context.Context, registry runtime.Registry) error {
	de.locker.Lock()
	de.registry = registry
	de.locker.Unlock()
	return nil
}

var _ = (runtime.EngineExtensionModule)((*deferEngine)(nil))

// EngineTypes returns a slice of reflect.Type representing the various engine types
// associated with the defer engine. These types include:
//
//   - DeferModule: Represents a module that can defer the initialization of the application.
func (de *deferEngine) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[DeferModule](),
	}
}

var _ = (injection.ComponentProvider)((*deferEngine)(nil))

// Components returns a slice of injection.Component representing the components
// provided by the defer engine. Currently, the only component provided is the
// deferEngine itself, which is scoped externally.
func (de *deferEngine) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(de, injection.ComponentExternalScope),
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
		if ctxErr := runtime.EnsureContext(ctx); ctxErr != nil {
			return ctxErr
		}
		return module.Defer(ctx)
	})
}
