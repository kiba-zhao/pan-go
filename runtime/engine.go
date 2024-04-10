package runtime

import (
	"errors"
	"reflect"
	"slices"

	"golang.org/x/sync/errgroup"
)

type EngineExtensionModule interface {
	EngineTypes() []reflect.Type
}

type Module interface {
	TypeOfModule() []reflect.Type
}

type ProviderModule interface {
	GetModules() []interface{}
}

type InitializeModule interface {
	Init(registry Registry) error
}

type ReadyModule interface {
	Ready() error
}

type Context = *errgroup.Group

var ErrDuplicateExtType = errors.New("[engine:Engine] Mount Error: duplicate extension type")

type Engine struct {
	Registry Registry
	extTypes []reflect.Type
}

func New() *Engine {
	engine := &Engine{}
	engine.Registry = NewRegistry()
	engine.extTypes = []reflect.Type{
		reflect.TypeFor[InitializeModule](),
		reflect.TypeFor[ReadyModule](),
	}
	return engine
}

func (engine *Engine) Mount(modules ...interface{}) error {
	var err error
	for _, module := range modules {

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

		if m, ok := module.(Module); ok {
			types = append(types, m.TypeOfModule()...)
		}

		if len(types) > 0 {
			err = engine.Registry.AddModule(module, types...)
		}

		if err != nil {
			break
		}

		if providerModule, ok := module.(ProviderModule); ok {
			modules := providerModule.GetModules()
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

func (engine *Engine) Bootstrap() (Context, error) {
	registry := engine.Registry
	err := TraverseModules(registry, func(module InitializeModule) error {
		return module.Init(registry)
	})

	var ctx errgroup.Group
	if err == nil {
		TraverseModules(registry, func(module ReadyModule) error {
			ctx.Go(module.Ready)
			return nil
		})
	}
	return &ctx, err
}
