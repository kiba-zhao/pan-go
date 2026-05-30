//go:build android || ios

package gomobile

import (
	"context"
	"pan/internal/app"
	"pan/internal/bootstrap"
	"pan/internal/injection"
)

type stdModule struct {
	Applet app.Applet

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
	module.agent.SetupApplet(module.Applet)
	return nil
}
