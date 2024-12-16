package sample

import (
	"pan/app/injection"
	"reflect"
)

func AppendSampleComponent[T any](components []injection.Component, component T) []injection.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, injection.NewComponentByType(reflect.TypeOf(component), component, injection.ComponentNoneScope))
	}
	components = append(components, injection.NewComponent(component, injection.ComponentInternalScope))
	return components
}

func AppendSampleExternalComponent[T any](components []injection.Component, component T) []injection.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, injection.NewComponentByType(reflect.TypeOf(component), component, injection.ComponentInternalScope))
	}
	components = append(components, injection.NewComponent(component, injection.ComponentExternalScope))
	return components
}

func AppendSampleInternalComponent[T any](components []injection.Component, component T) []injection.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, injection.NewComponentByType(reflect.TypeOf(component), component, injection.ComponentInternalScope))
	}
	components = append(components, injection.NewComponent(component, injection.ComponentInternalScope))
	return components
}
