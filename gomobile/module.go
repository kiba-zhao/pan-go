package gomobile

import (
	"context"
	"pan/lib/bootstrap"
	"pan/lib/injection"
	libSerlvet "pan/lib/serlvet"
)

type stdModule struct {
	Serlvet *libSerlvet.Serlvet

	agent *GoMobileAgent
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (module *stdModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(module, injection.ComponentNoneScope),
	}
}

var _ = (bootstrap.DeferModule)((*stdModule)(nil))

func (module *stdModule) Defer(ctx context.Context) error {
	setAgentSerlvet(module.agent, module.Serlvet)

	return nil
}
