// Package injection provides the injection engine
//
// The injection engine is used to manage the dependencies of the application
package injection

import (
	"pan/lib/runtime"
	"reflect"
)

type stdInjectionModule struct {
}

func New() interface{} {
	return &stdInjectionModule{}
}

var _ = (runtime.InitializeModule)((*stdInjectionModule)(nil))

func (in *stdInjectionModule) Init(registry runtime.Registry) error {
	engine := newStdInjectEngine()

	err := runtime.TraverseRegistry(registry, func(provider ComponentProvider) error {
		var internalStore ComponentStore
		if storeProvider, ok := provider.(ComponentStoreProvider); ok {
			internalStore = storeProvider.ComponentStore()
		}

		components := provider.Components()
		return engine.InjectComponents(internalStore, components...)
	})
	return err
}

var _ = (runtime.EngineExtensionModule)((*stdInjectionModule)(nil))

func (in *stdInjectionModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[ComponentProvider](),
	}
}
