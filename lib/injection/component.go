// Define injection component
package injection

import "reflect"

const (
	ComponentNoneScope     = "none"
	ComponentInternalScope = "internal"
	ComponentExternalScope = "external"
)

type Component interface {
	// Type returns the reflect.Type of the component.
	Type() reflect.Type
	// Target returns the target object of the component.
	//
	// This can be used to get the value of the component.
	Target() interface{}
	// Scope returns the scope of the component.
	Scope() string
}

type componentBase struct {
	ty    reflect.Type
	scope string
}

func (c *componentBase) Type() reflect.Type {
	return c.ty
}

func (c *componentBase) Scope() string {
	return c.scope
}

type stdComponent struct {
	componentBase
	target interface{}
}

// NewComponent creates a new component with a specified target and scope.
// The target is the object the component represents, and the scope determines
// the component's visibility or lifecycle. The function returns a Component
// interface implementation that can be used in the injection system.

func NewComponent[T any](target T, scope string) Component {
	base := componentBase{
		ty:    reflect.TypeFor[T](),
		scope: scope,
	}

	return &stdComponent{
		target:        target,
		componentBase: base,
	}
}

// NewComponentByType creates a new component with a specified reflect.Type, target, and scope.
// The reflect.Type `ty` represents the type of the component,
// `target` is the object the component represents, and `scope` determines
// the component's visibility or lifecycle. This function returns an implementation
// of the Component interface, which can be used within the injection system.

func NewComponentByType(ty reflect.Type, target interface{}, scope string) Component {

	base := componentBase{
		ty:    ty,
		scope: scope,
	}

	return &stdComponent{
		target:        target,
		componentBase: base,
	}
}

func (c *stdComponent) Target() interface{} {
	return c.target
}
