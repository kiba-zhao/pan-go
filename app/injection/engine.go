// Define injection engine
package injection

import (
	"errors"
	"pan/runtime"
	"reflect"
	"strings"
)

type ComponentPendings = map[reflect.Type][]reflect.Value

var ErrComponentConflict = errors.New("[app.injection] engine Error: dependency conflict")
var ErrComponentScope = errors.New("[app.injection] engine Error: invalid component scope")

type engine struct {
}

// Init initializes the engine with the given registry.
//
// It will traverse all modules and try to inject all components.
// If any error occurs during injection, it will return the error.
func (ie *engine) Init(registry runtime.Registry) error {
	store := NewComponentStore()
	pendings := make(ComponentPendings)

	// traverse component provider
	err := runtime.TraverseRegistry(registry, func(provider ComponentProvider) error {
		var internalStore ComponentStore
		if storeProvider, ok := provider.(ComponentStoreProvider); ok {
			internalStore = storeProvider.ComponentStore()
		}
		if internalStore == nil {
			internalStore = NewComponentStore()
		}
		var componentErr error
		components := provider.Components()
		for _, component := range components {
			// inject component
			componentErr = inject(component, pendings, internalStore, store)
			if componentErr != nil {
				break
			}
		}
		return componentErr
	})
	return err
}

// EngineTypes returns a slice of reflect.Type representing the various engine types
// associated with the injection module. These types include:
//
//   - ComponentProvider: Provides components for the injection module.
func (in *engine) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[ComponentProvider](),
	}
}

// inject injects the given component into the given internal and external stores.
//
// It also handles pendings and injects the component into the fields of the target.
// If any error occurs during injection, it will return the error.
//
// The rules of injection are as follows:
//
//  2. If a component is registered, but the target is not a pointer or a struct, it will be ignored.
//  3. If a component is registered, but the target is a pointer or a struct, it will be injected into the target.
//  4. If a component is registered and the target is a pointer or a struct, but the field is not exported, it will be ignored.
//  5. If a component is registered and the target is a pointer or a struct, but the field is exported and is a pointer or a struct, it will be injected into the field recursively.
//  6. If a component is registered and the target is a pointer or a struct, but the field is exported and is not a pointer or a struct, it will be ignored.
//  7. If a component is registered and the target is a pointer or a struct, but the field is exported and has a tag "inject" with value "-", it will be ignored.
//  8. If a component is registered and the target is a pointer or a struct, but the field is exported and has a tag "inject" with value "volatile", it will be injected into the field, but the value will be set to zero after injection.
//  9. If a component is registered and the target is a pointer or a struct, but the field is exported and has a tag "inject" with value "volatile;", it will be injected into the field, but the value will be set to zero after injection.
//
// 10. If a component is registered and the target is a pointer or a struct, but the field is exported and has a tag "inject" with value ";volatile", it will be injected into the field, but the value will be set to zero after injection.
func inject(component Component, pendings ComponentPendings, internalStore ComponentStore, store ComponentStore) error {

	t := component.Type()
	target := component.Target()
	scope := component.Scope()

	var store_ ComponentStore
	switch scope {
	case ComponentInternalScope:
		store_ = internalStore
	case ComponentExternalScope:
		store_ = store
	case ComponentNoneScope:
	default:
		return ErrComponentScope
	}

	if store_ != nil {
		// store component and handle pendings
		if _, ok := store_[t]; ok {
			return ErrComponentConflict
		}

		if values, ok := pendings[t]; ok {
			for _, v := range values {
				v.Set(reflect.ValueOf(target))
			}
			delete(pendings, t)
		}

		store_[t] = target
		if t.Kind() == reflect.Interface {
			return nil
		}
	}

	// inject component fields and add  pendings fields
	et := reflect.TypeOf(target)
	if et.Kind() == reflect.Ptr {
		et = et.Elem()
	}
	fields := reflect.VisibleFields(et)
	v := reflect.ValueOf(target)
	iv := reflect.Indirect(v)
	for _, field := range fields {
		if !field.IsExported() {
			continue
		}
		if !(field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Interface) {
			continue
		}

		volatile := false
		if tag, ok := field.Tag.Lookup("inject"); ok {
			if tag == "-" {
				continue
			}
			volatileIdx := strings.Index(tag, "volatile")
			if volatileIdx == 0 || (volatileIdx > 0 && tag[volatileIdx-1] == ';') {
				volatile = true
			}
		}
		fv := iv.FieldByName(field.Name)
		if !fv.IsZero() && !volatile {
			continue
		}
		field_target, ok := internalStore[field.Type]
		if !ok {
			field_target, ok = store[field.Type]
		}
		if ok {
			fv.Set(reflect.ValueOf(field_target))
			continue
		}

		if values, ok := pendings[field.Type]; ok {
			pendings[field.Type] = append(values, fv)
			continue
		}
		pendings[field.Type] = []reflect.Value{fv}
	}
	return nil
}
