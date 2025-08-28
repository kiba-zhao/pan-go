package settings

import (
	"context"
	"errors"
	"pan/lib/bootstrap"
	"pan/lib/config"
	"pan/lib/feature"
	"pan/lib/injection"
	"pan/lib/log"
	"pan/lib/runtime"
	"pan/lib/serlvet"
	"pan/lib/web"
	"sync"

	"github.com/spf13/viper"
)

var subModuleNewFuncArray []feature.SubModuleNewFunc[*stdModule]

var ErrSettingsModuleUavailable = errors.New("settings.Module Error: Unavailable")

const (
	SettingsModuleName = "settings"
)

type stdModule struct {
	settingsCtrl *SettingsController
	settingsSrv  *SettingsServlet
	settingsSvc  *SettingsService
	configurer   SettingsConfigurer

	viper *viper.Viper

	componentStore     injection.ComponentStore
	componentStoreOnce sync.Once

	locker   sync.Mutex
	homePath string

	configListeners []SettingsConfigListener
}

func New() interface{} {

	viper := viper.New()
	settingsCtrl := &SettingsController{}
	settingsSrv := &SettingsServlet{}
	settingsSvc := &SettingsService{}

	logger := log.Default()
	configurer := config.NewConfigurer[SettingsConfig](logger)

	module := &stdModule{}
	module.configurer = configurer
	module.viper = viper
	module.settingsCtrl = settingsCtrl
	module.settingsSrv = settingsSrv
	module.settingsSvc = settingsSvc

	return module
}

var _ = (SettingsConfigListener)((*stdModule)(nil))

func (m *stdModule) OnConfigUpdated(cfg SettingsConfig) {
	m.locker.Lock()
	defer m.locker.Unlock()
	viper := m.viper

	homePath := cfg.HomePath()
	if homePath != m.homePath {
		m.homePath = homePath
		initViper(viper, homePath)
	}
	m.settingsSvc.Setup(cfg)

	configListeners := m.configListeners
	if len(configListeners) <= 0 {
		return
	}

	for _, configListener := range configListeners {
		configListener.OnConfigUpdated(cfg)
	}
}

var _ = (bootstrap.DeferModule)((*stdModule)(nil))

func (m *stdModule) Defer(ctx context.Context) error {
	m.configurer.Subscribe(m)
	return nil
}

var _ = (bootstrap.DestroyModule)((*stdModule)(nil))

func (m *stdModule) Destroy() {
	m.configurer.Unsubscribe(m)
}

var _ = (web.WebAppModule)((*stdModule)(nil))

func (m *stdModule) SetupToWeb(app web.WebApp) error {
	err := m.settingsCtrl.SetupToWeb(app.Group(SettingsModuleName))
	return err
}

var _ = (serlvet.SerlvetModule)((*stdModule)(nil))

func (m *stdModule) SetupToSerlvet(app serlvet.SerlvetApp) error {
	err := m.settingsSrv.SetupToSerlvet(app.Route([]byte(SettingsModuleName)))
	return err
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (m *stdModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m.viper, injection.ComponentInternalScope),
		// configurer
		injection.NewComponent(m.configurer, injection.ComponentExternalScope),
		// controller
		injection.NewComponent(m.settingsCtrl, injection.ComponentNoneScope),
		// servlet
		injection.NewComponent(m.settingsSrv, injection.ComponentNoneScope),
		// service
		injection.NewComponent(m.settingsSvc, injection.ComponentInternalScope),
	}
}

var _ = (injection.ComponentStoreProvider)((*stdModule)(nil))

func (m *stdModule) ComponentStore() injection.ComponentStore {
	m.componentStoreOnce.Do(func() {
		m.componentStore = injection.NewComponentStore()
	})
	return m.componentStore
}

var _ = (runtime.ProviderModule)((*stdModule)(nil))

func (m *stdModule) Modules() []interface{} {
	m.configListeners = make([]SettingsConfigListener, 0)
	return feature.NewSubModules(m, subModuleNewFuncArray...)
}
