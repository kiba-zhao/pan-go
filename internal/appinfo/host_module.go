//go:build !(android || ios)

package appinfo

import (
	"pan/internal/injection"
	"pan/internal/module"
	"pan/internal/web"
	"sync"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newHostModule)
}

func newHostModule(m *stdModule) interface{} {
	hostModule := &stdHostModule{}
	hostModule.SubModule = module.NewSubModule(hostModule, m)
	return hostModule
}

type stdHostModule struct {
	*module.SubModule[*stdHostModule, *stdModule]

	controllers     []web.WebController
	controllersOnce sync.Once
}

var _ = (web.WebRouteModule)((*stdHostModule)(nil))

func (m *stdHostModule) WebRouteName() string {
	return ModuleName
}

var _ = (web.WebControllerProvider)((*stdHostModule)(nil))

func (m *stdHostModule) WebControllers() []web.WebController {
	m.controllersOnce.Do(func() {
		m.controllers = []web.WebController{
			&QRCodeController{},
		}
	})
	return m.controllers
}

var _ = (injection.ComponentProvider)((*stdHostModule)(nil))

func (m *stdHostModule) Components() []injection.Component {
	components := make([]injection.Component, 0)
	for _, ctrl := range m.WebControllers() {
		components = append(components, injection.NewComponent(ctrl, injection.ComponentNoneScope))
	}
	return components
}
