package sample

import (
	"pan/app/injection"
	"reflect"
)

// AppendSampleComponent appends a new component to the provided slice of components.
// If the type of the component is an interface, it is added using NewComponentByType
// with ComponentNoneScope. The component is then always added using NewComponent
// with ComponentInternalScope. The function returns the updated slice of components.

func AppendSampleComponent[T any](components []injection.Component, component T) []injection.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, injection.NewComponentByType(reflect.TypeOf(component), component, injection.ComponentNoneScope))
	}
	components = append(components, injection.NewComponent(component, injection.ComponentInternalScope))
	return components
}

// AppendSampleExternalComponent appends a new component to the provided slice of components.
// If the type of the component is an interface, it is added using NewComponentByType
// with ComponentInternalScope. The component is then always added using NewComponent
// with ComponentExternalScope. The function returns the updated slice of components.

func AppendSampleExternalComponent[T any](components []injection.Component, component T) []injection.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, injection.NewComponentByType(reflect.TypeOf(component), component, injection.ComponentInternalScope))
	}
	components = append(components, injection.NewComponent(component, injection.ComponentExternalScope))
	return components
}

// AppendSampleInternalComponent appends a new component to the provided slice of components.
// If the type of the component is an interface, it is added using NewComponentByType
// with ComponentInternalScope. The component is then always added using NewComponent
// with ComponentInternalScope. The function returns the updated slice of components.

func AppendSampleInternalComponent[T any](components []injection.Component, component T) []injection.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, injection.NewComponentByType(reflect.TypeOf(component), component, injection.ComponentInternalScope))
	}
	components = append(components, injection.NewComponent(component, injection.ComponentInternalScope))
	return components
}
