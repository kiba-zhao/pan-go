package sample

import (
	"pan/app/bootstrap"
	"reflect"
)

func AppendSampleComponent[T any](components []bootstrap.Component, component T) []bootstrap.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, bootstrap.NewComponentByType(reflect.TypeOf(component), component, bootstrap.ComponentNoneScope))
	}
	components = append(components, bootstrap.NewComponent(component, bootstrap.ComponentInternalScope))
	return components
}

func AppendSampleExternalComponent[T any](components []bootstrap.Component, component T) []bootstrap.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, bootstrap.NewComponentByType(reflect.TypeOf(component), component, bootstrap.ComponentInternalScope))
	}
	components = append(components, bootstrap.NewComponent(component, bootstrap.ComponentExternalScope))
	return components
}

func AppendSampleInternalComponent[T any](components []bootstrap.Component, component T) []bootstrap.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, bootstrap.NewComponentByType(reflect.TypeOf(component), component, bootstrap.ComponentInternalScope))
	}
	components = append(components, bootstrap.NewComponent(component, bootstrap.ComponentInternalScope))
	return components
}
