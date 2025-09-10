//go:build android || ios

package gomobile

import (
	"context"
	"pan/lib/bootstrap"
	"pan/lib/injection"
	libServlet "pan/lib/servlet"
)

type stdModule struct {
	Servlet libServlet.Servlet

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
	setupAgentServlet(module.agent, module.Servlet)
	return nil
}
