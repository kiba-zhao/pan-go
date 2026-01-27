// Define injection engine
package injection

import (
	"errors"
	"reflect"
	"strings"
)

type ComponentPendings = map[reflect.Type][]reflect.Value

var ErrComponentConflict = errors.New("[app.injection] engine Error: dependency conflict")
var ErrComponentScope = errors.New("[app.injection] engine Error: invalid component scope")

type stdInjectEngine struct {
	store    ComponentStore
	pendings ComponentPendings
}

func (engine *stdInjectEngine) InjectComponents(store ComponentStore, components ...Component) error {
	internalStore := store
	if internalStore == nil {
		internalStore = NewComponentStore()
	}

	var err error
	for _, component := range components {
		// inject component
		err = inject(engine, internalStore, component)
		if err != nil {
			break
		}
	}

	return err
}

func inject(engine *stdInjectEngine, internalStore ComponentStore, component Component) error {

	t := component.Type()
	target := component.Target()
	scope := component.Scope()

	var store_ ComponentStore
	switch scope {
	case ComponentInternalScope:
		store_ = internalStore
	case ComponentExternalScope:
		store_ = engine.store
	case ComponentNoneScope:
	default:
		return ErrComponentScope
	}

	if store_ != nil {
		// store component and handle pendings
		if _, ok := store_[t]; ok {
			return ErrComponentConflict
		}

		if values, ok := engine.pendings[t]; ok {
			for _, v := range values {
				v.Set(reflect.ValueOf(target))
			}
			delete(engine.pendings, t)
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
			field_target, ok = engine.store[field.Type]
		}
		if ok {
			fv.Set(reflect.ValueOf(field_target))
			continue
		}

		if values, ok := engine.pendings[field.Type]; ok {
			engine.pendings[field.Type] = append(values, fv)
			continue
		}
		engine.pendings[field.Type] = []reflect.Value{fv}
	}
	return nil
}

func newStdInjectEngine() *stdInjectEngine {
	engine := &stdInjectEngine{}
	engine.pendings = make(ComponentPendings)
	engine.store = NewComponentStore()

	return engine
}
