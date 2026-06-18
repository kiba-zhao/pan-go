//go:build !(android || ios)

package settings

import (
	"context"
	"pan/pkg/bootstrap"
	"pan/pkg/injection"
	"pan/pkg/module"
	"pan/pkg/web"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newHostModule)
}

func newHostModule(m *stdModule) interface{} {
	baseSettingsCtrl := &BaseSettingsController{}
	hostCtrl := &HostSettingsController{}
	hostSettingsSvc := &HostSettingsService{}

	hostModule := &stdHostSettingsModule{}
	hostModule.SubModule = module.NewSubModule(hostModule, m)
	hostModule.baseSettingsCtrl = baseSettingsCtrl
	hostModule.hostSettingsCtrl = hostCtrl
	hostModule.hostSettingsSvc = hostSettingsSvc

	return hostModule
}

const (
	HostSettingsModuleName = "host-settings"
)

type stdHostSettingsModule struct {
	WebConfigurer web.WebConfigurer

	*module.SubModule[*stdHostSettingsModule, *stdModule]

	baseSettingsCtrl *BaseSettingsController
	hostSettingsCtrl *HostSettingsController
	hostSettingsSvc  *HostSettingsService
}

var _ = (web.WebAppModule)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) SetupToWeb(app web.WebApp) error {
	route := app.Group(web.WEB_API_PATH)
	err := m.baseSettingsCtrl.SetupToWeb(route.Group(SettingsModuleName))
	if err == nil {
		err = m.hostSettingsCtrl.SetupToWeb(route.Group(HostSettingsModuleName))
	}
	return err
}

var _ = (bootstrap.DeferModule)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) Defer(ctx context.Context) error {
	hostSettings, err := m.loadHostSettings()
	if err == nil {
		err = m.configure(hostSettings)
	}
	return err
}

var _ = (injection.ComponentProvider)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentInternalScope),
		// controller
		injection.NewComponent(m.baseSettingsCtrl, injection.ComponentNoneScope),
		injection.NewComponent(m.hostSettingsCtrl, injection.ComponentNoneScope),
		// service
		injection.NewComponent(m.hostSettingsSvc, injection.ComponentInternalScope),
	}
}

var _ = (HostSettingsChangedTrigger)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) OnHostSettingsChanged(hostSettings HostSettings) {
	m.configure(hostSettings)
}

func (m *stdHostSettingsModule) configure(hostSettings HostSettings) error {
	return m.configureWeb(&hostSettings)
}

func (m *stdHostSettingsModule) configureWeb(hostSettings *HostSettings) error {
	webConfig := newWebConfig(hostSettings)
	return m.WebConfigurer.Configure(webConfig)
}

func (m *stdHostSettingsModule) loadHostSettings() (HostSettings, error) {
	hostSettings, err := m.hostSettingsSvc.Load()
	if err != nil {
		m.ParentModule().logger.Error("settings.HostSettingsModule", "loadHostSettings Error: "+err.Error())
	}
	return hostSettings, err
}
