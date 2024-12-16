package injection

import "reflect"

const (
	ComponentNoneScope     = "none"
	ComponentInternalScope = "internal"
	ComponentExternalScope = "external"
)

type Component interface {
	Type() reflect.Type
	Target() interface{}
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

type componentImpl struct {
	componentBase
	target interface{}
}

func NewComponent[T any](target T, scope string) Component {
	base := componentBase{
		ty:    reflect.TypeFor[T](),
		scope: scope,
	}

	return &componentImpl{
		target:        target,
		componentBase: base,
	}
}

func NewComponentByType(ty reflect.Type, target interface{}, scope string) Component {

	base := componentBase{
		ty:    ty,
		scope: scope,
	}

	return &componentImpl{
		target:        target,
		componentBase: base,
	}
}

func (c *componentImpl) Target() interface{} {
	return c.target
}
