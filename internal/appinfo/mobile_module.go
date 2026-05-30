//go:build android || ios

package appinfo

import (
	"pan/internal/app"
	"pan/internal/injection"
	"pan/internal/module"
	"sync"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newMobileModule)
}

func newMobileModule(m *stdModule) interface{} {
	mobileModule := &stdMobileModule{}
	mobileModule.SubModule = module.NewSubModule(mobileModule, m)

	return mobileModule
}

type stdMobileModule struct {
	*module.SubModule[*stdMobileModule, *stdModule]

	modules     []app.AppletModule
	modulesOnce sync.Once
}

var _ = (app.AppletModule)((*stdMobileModule)(nil))

func (m *stdMobileModule) SetupToApplet(router app.AppServletRouter) error {
	router_ := router.Route([]byte(ModuleName))
	var err error
	for _, module := range m.AppModules() {
		err = module.SetupToApplet(router_)
		if err != nil {
			break
		}
	}
	return err
}

func (m *stdMobileModule) AppModules() []app.AppletModule {
	m.modulesOnce.Do(func() {
		m.modules = []app.AppletModule{
			&QRCodeAppletModule{},
		}
	})
	return m.modules
}

var _ = (injection.ComponentProvider)((*stdMobileModule)(nil))

func (m *stdMobileModule) Components() []injection.Component {
	components := make([]injection.Component, 0)
	for _, module := range m.AppModules() {
		components = append(components, injection.NewComponent(module, injection.ComponentNoneScope))
	}
	return components
}
