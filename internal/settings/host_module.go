//go:build !(android || ios)

package settings

import (
	"context"
	"pan/internal/bootstrap"
	"pan/internal/feature"
	"pan/internal/injection"
	"pan/internal/web"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newHostModule)
}

func newHostModule(module *stdModule) interface{} {
	settingsCtrl := &SettingsController{}
	hostCtrl := &HostSettingsController{}
	hostSettingsSvc := &HostSettingsService{}

	hostModule := &stdHostSettingsModule{}
	hostModule.SubModule = feature.NewSubModule(hostModule, module)
	hostModule.settingsCtrl = settingsCtrl
	hostModule.hostSettingsCtrl = hostCtrl
	hostModule.hostSettingsSvc = hostSettingsSvc

	return hostModule
}

const (
	HostSettingsModuleName = "host-settings"
)

type stdHostSettingsModule struct {
	WebConfigurer web.WebConfigurer

	*feature.SubModule[*stdHostSettingsModule, *stdModule]

	settingsCtrl     *SettingsController
	hostSettingsCtrl *HostSettingsController
	hostSettingsSvc  *HostSettingsService
}

var _ = (web.WebAppModule)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) SetupToWeb(app web.WebApp) error {
	err := m.settingsCtrl.SetupToWeb(app.Group(feature.WebAPIPrefix + SettingsModuleName))
	if err == nil {
		err = m.hostSettingsCtrl.SetupToWeb(app.Group(feature.WebAPIPrefix + HostSettingsModuleName))
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
		injection.NewComponent(m.settingsCtrl, injection.ComponentNoneScope),
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
		m.ParentModule().logger.Error("HostSettingsModule", "loadHostSettings Error: "+err.Error())
	}
	return hostSettings, err
}
