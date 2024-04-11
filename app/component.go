package app

import (
	"errors"
	"pan/runtime"
	"reflect"
)

type ComponentStore = map[reflect.Type]interface{}
type ComponentPendings = map[reflect.Type][]reflect.Value

type Component interface {
	Type() reflect.Type
	Target() interface{}
	Scope() string
}

const (
	ComponentNoneScope     = ""
	ComponentInternalScope = "internal"
	ComponentExternalScope = "external"
)

type componentBase[T any] struct {
	scope string
}

func (c *componentBase[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

func (c *componentBase[T]) Scope() string {
	return c.scope
}

type componentImpl[T any] struct {
	componentBase[T]
	target T
}

func NewComponent[T any](target T, scope string) Component {

	base := componentBase[T]{
		scope: scope,
	}

	return &componentImpl[T]{
		target:        target,
		componentBase: base,
	}
}

func (c *componentImpl[T]) Target() interface{} {
	return c.target
}

type LazyComponentFunc[T any] func() T

type lazyComponentImpl[T any] struct {
	componentBase[T]
	lazyFunc LazyComponentFunc[T]
}

func NewLazyComponent[T any](lazyFunc LazyComponentFunc[T], scope string) Component {
	base := componentBase[T]{
		scope: scope,
	}
	return &lazyComponentImpl[T]{
		lazyFunc:      lazyFunc,
		componentBase: base,
	}
}

func (c *lazyComponentImpl[T]) Target() interface{} {
	return c.lazyFunc()
}

type ComponentProvider interface {
	Components() []Component
}

var ErrComponentConflict = errors.New("[app:Component] Injector Error: dependency conflict")
var ErrComponentScope = errors.New("[app:Component] Injector Error: invalid component scope")

type injector struct {
}

func (in *injector) Init(registry runtime.Registry) error {
	store := make(ComponentStore)
	pendings := make(ComponentPendings)

	// traverse component provider
	err := runtime.TraverseRegistry(registry, func(provider ComponentProvider) error {
		internalStore := make(ComponentStore)
		var componentErr error
		components := provider.Components()
		for _, component := range components {
			// inject component
			componentErr = injectComponent(component, pendings, internalStore, store)
			if componentErr != nil {
				break
			}
		}
		return componentErr
	})
	return err
}

func (in *injector) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[ComponentProvider](),
	}
}

func injectComponent(component Component, pendings ComponentPendings, internalStore ComponentStore, store ComponentStore) error {

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
		if tag, ok := field.Tag.Lookup("inject"); ok && tag == "-" {
			continue
		}
		fv := iv.FieldByName(field.Name)
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
