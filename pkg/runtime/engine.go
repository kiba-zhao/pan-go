/*
# Engine

The engine is a runtime environment that provides a set of modules that
can be used to extend and customize the behavior of the application.
It allows modules to be registered and accessed by other modules, and it
provides a simple and efficient way to manage and access modules in the
application.
*/
package runtime

import (
	"context"
	"errors"
	"reflect"
	"slices"

	"golang.org/x/sync/errgroup"
)

// Engine extension module interface
//
// The EngineTypes method will be available in Engine After Mount,
// add the module types supported by the Engine extension type returned.
type EngineExtensionModule interface {

	// EngineTypes returns the types supported by the Engine extension
	EngineTypes() []reflect.Type
}

// Engine module interface
//
// Module interface for specifying engine type
type Module interface {

	// Method for specifying module specific engine types
	TypeOfModule() []reflect.Type
}

// Engine module interface for provide sub-modules
type ProviderModule interface {

	// Modules returns a list of sub-modules provided by the module.
	// These sub-modules are typically used to extend or enhance
	// the functionality of the parent module or engine.
	Modules() []interface{}
}

// Engine module interface for Module that need to be initialized
type InitializeModule interface {

	// Init initializes the module with the provided registry.
	// It sets up necessary configurations and dependencies required
	// for the module to function correctly within the engine.
	// Returns an error if initialization fails.
	Init(ctx context.Context, registry Registry) error
}

type Context = *errgroup.Group

// When attempting to add an existing engine type, this error will be returned
var ErrDuplicateExtType = errors.New("[engine:Engine] Mount Error: duplicate extension type")

// Runtime Engine
type Engine struct {
	Registry Registry       // Registration of engine modules
	extTypes []reflect.Type // Types supported by the engine
}

// Create a new Engine instance.
func New(modules ...interface{}) (*Engine, error) {
	engine := &Engine{}
	engine.Registry = NewRegistry()
	engine.extTypes = []reflect.Type{
		reflect.TypeFor[InitializeModule](),
	}

	var err error
	if len(modules) > 0 {
		err = engine.Mount(modules...)
	}
	return engine, err
}

// Mount module to engine.
//
// The module can be either an InitializeModule, a ProviderModule, or a Module.
// If the module is an InitializeModule, it will be initialized by calling the
// Init method with the engine's registry.
// If the module is a ProviderModule, its Modules method will be called and the
// returned modules will be mounted recursively.
// If the module is a Module, its TypeOfModule method will be called and the
// returned types will be used to register the module with the engine's registry.
//
// The Mount method will return an error if any of the modules fail to mount.
// If an error is returned, the engine's registry will not be modified.
func (engine *Engine) Mount(modules ...interface{}) error {
	var err error
	for _, module := range modules {

		// add extension types if module is Implemented EngineExtensionModule interface
		if extModule, ok := module.(EngineExtensionModule); ok {
			extTypes := extModule.EngineTypes()
			for _, extType := range extTypes {
				if slices.Contains(engine.extTypes, extType) {
					err = ErrDuplicateExtType
					break
				}
				engine.extTypes = append(engine.extTypes, extType)
			}
		}

		if err != nil {
			break
		}

		// Extract the engine type implemented by the module
		t := reflect.TypeOf(module)
		types := []reflect.Type{}
		for _, extType := range engine.extTypes {

			if extType.Kind() == reflect.Interface {
				if reflect.TypeOf(module).Implements(extType) {
					types = append(types, extType)
				}
				continue
			}

			if extType == t {
				types = append(types, extType)
			}
		}

		// Add engine type specified by Module
		if m, ok := module.(Module); ok {
			types = append(types, m.TypeOfModule()...)
		}

		// Register according to the engine types supported by the module
		if len(types) > 0 {
			err = engine.Registry.Append(module, types...)
		}

		if err != nil {
			break
		}

		// Mount the submodules provided by the ProviderModule interface
		if providerModule, ok := module.(ProviderModule); ok {
			modules := providerModule.Modules()
			if len(modules) > 0 {
				err = engine.Mount(modules...)
			}
			if err != nil {
				break
			}
		}

	}
	return err
}

// Bootstrap initializes all modules that implement the InitializeModule interface.
// It calls the Init method with the engine's registry on each of the modules.
// If any of the modules return an error during initialization, the error will be
// returned and initialization will halt.
func (engine *Engine) Bootstrap(ctx context.Context) error {
	registry := engine.Registry
	return TraverseRegistry(registry, func(module InitializeModule) error {
		if err := EnsureContext(ctx); err != nil {
			return err
		}
		return module.Init(ctx, registry)
	})
}
