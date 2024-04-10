package app

import (
	"errors"
	"pan/runtime"
	"reflect"
)

type DependencyProvider interface {
	GetDependencies() map[reflect.Type]interface{}
}

type ComponentProvider interface {
	GetComponents() []interface{}
}

var ErrComponentConflict = errors.New("[app:Component] Injector Error: dependency conflict")

type injector struct {
}

func (in *injector) Init(registry runtime.Registry) error {
	store := make(map[reflect.Type]interface{})
	pendings := make(map[reflect.Type][]reflect.Value)

	// inject dependencies to component
	inject := func(component interface{}) {
		et := reflect.TypeOf(component)
		if et.Kind() == reflect.Ptr {
			et = et.Elem()
		}
		fields := reflect.VisibleFields(et)
		v := reflect.ValueOf(component)
		iv := reflect.Indirect(v)
		for _, field := range fields {
			if !field.IsExported() {
				continue
			}
			if !(field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Interface) {
				continue
			}
			if tag, ok := field.Tag.Lookup("inject"); ok && tag == "-" {
				continue
			}
			fv := iv.FieldByName(field.Name)
			dependency_, ok := store[field.Type]
			if ok {
				fv.Set(reflect.ValueOf(dependency_))
				continue
			}

			if values, ok := pendings[field.Type]; ok {
				pendings[field.Type] = append(values, fv)
				continue
			}
			pendings[field.Type] = []reflect.Value{fv}
		}
	}

	// traverse component provider
	err := runtime.TraverseModules(registry, func(provider ComponentProvider) error {
		components := provider.GetComponents()
		for _, component := range components {
			inject(component)
		}
		return nil
	})

	if err == nil {
		// traverse dependency provider
		err = runtime.TraverseModules(registry, func(provider DependencyProvider) error {
			var dependencyErr error
			dependencies := provider.GetDependencies()
			for t, dependency := range dependencies {

				if _, ok := store[t]; ok {
					dependencyErr = ErrComponentConflict
					break
				}

				if values, ok := pendings[t]; ok {
					for _, v := range values {
						v.Set(reflect.ValueOf(dependency))
					}
					delete(pendings, t)
				}
				store[t] = dependency
				if t.Kind() == reflect.Interface {
					continue
				}

				inject(dependency)

			}
			return dependencyErr
		})
	}
	return err
}

func (in *injector) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[ComponentProvider](),
		reflect.TypeFor[DependencyProvider](),
	}
}
