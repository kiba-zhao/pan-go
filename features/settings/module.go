package settings

import (
	"context"
	"errors"
	"pan/lib/bootstrap"
	"pan/lib/broadcast"
	"pan/lib/config"
	"pan/lib/feature"
	"pan/lib/injection"
	"pan/lib/log"
	"pan/lib/quic"
	"pan/lib/repository"
	"pan/lib/runtime"
	"slices"
	"sync"

	"github.com/spf13/viper"
)

var subModuleNewFuncArray []feature.SubModuleNewFunc[*stdModule]

var ErrSettingsModuleUavailable = errors.New("settings.Module Error: Unavailable")

const (
	SettingsModuleName = "settings"
)

type stdModule struct {
	QuicConfigurer       quic.QuicConfigurer
	BroadcastConfigurer  broadcast.BroadcastConfigurer
	RepositoryConfigurer repository.RepositoryConfigurer

	logger log.Logger

	settingsSvc        *SettingsService
	configurer         SettingsConfigurer
	securityConfigurer SecurityConfigurer

	viper *viper.Viper

	componentStore     injection.ComponentStore
	componentStoreOnce sync.Once

	locker      sync.Mutex
	settingsCfg SettingsConfig

	subModules       []interface{}
	subModulesRW     sync.RWMutex
	ignoreConfigured bool

	securityConfigListener SecurityConfigurerListener
}

func New() interface{} {

	viper := viper.New()

	settingsSvc := &SettingsService{}

	logger := log.Default()
	configurer := config.NewConfigurer[SettingsConfig](logger)
	securityConfigurer := config.NewConfigurer[SecurityConfig](logger)

	module := &stdModule{}
	module.logger = logger
	module.configurer = configurer
	module.securityConfigurer = securityConfigurer
	module.viper = viper
	module.settingsSvc = settingsSvc

	securityConfigListener := &stdModuleSecurityConfigProxy{}
	securityConfigListener.module = module
	module.securityConfigListener = securityConfigListener

	return module
}

var _ = (SettingsConfigListener)((*stdModule)(nil))

func (m *stdModule) OnConfigUpdated(cfg SettingsConfig) {
	m.locker.Lock()
	defer m.locker.Unlock()

	m.settingsSvc.Setup(cfg)

	// init viper and securityConfig
	homePath := cfg.HomePath()
	err := m.initViper(homePath)
	if err == nil {
		err = m.initSecurityConfig(homePath)
	}
	if err != nil {
		return
	}
	//

}

var _ = (bootstrap.DeferModule)((*stdModule)(nil))

func (m *stdModule) Defer(ctx context.Context) error {
	m.configurer.Subscribe(m)
	m.securityConfigurer.Subscribe(m.securityConfigListener)
	return nil
}

var _ = (bootstrap.DestroyModule)((*stdModule)(nil))

func (m *stdModule) Destroy() {
	m.configurer.Unsubscribe(m)
	m.securityConfigurer.Unsubscribe(m.securityConfigListener)
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (m *stdModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentInternalScope),
		injection.NewComponent[SettingsChangedTrigger](m, injection.ComponentInternalScope),
		injection.NewComponent(m.viper, injection.ComponentInternalScope),
		// configurer
		injection.NewComponent(m.configurer, injection.ComponentExternalScope),
		injection.NewComponent(m.securityConfigurer, injection.ComponentExternalScope),
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
	m.subModulesRW.Lock()
	defer m.subModulesRW.Unlock()
	m.subModules = feature.NewSubModules(m, subModuleNewFuncArray...)
	return m.subModules
}

var _ = (SettingsChangedTrigger)((*stdModule)(nil))

func (m *stdModule) OnSettingsChanged(settings Settings) {
	if !m.ignoreConfigured {
		securityCfg := m.securityConfigurer.Config()
		m.configure(securityCfg, settings, nil, false)
		return
	}

	subModules := m.loadSubModules()
	if len(subModules) <= 0 {
		return
	}

	for _, subModule := range subModules {
		if settingsTrigger, ok := subModule.(SettingsChangedTrigger); ok {
			settingsTrigger.OnSettingsChanged(settings)
		}
	}

}

func (m *stdModule) initViper(homePath string) error {
	err := initViper(m.viper, homePath)
	if err != nil {
		m.logger.Error("SettingsModule", "initViper Error: "+err.Error())
	}
	return err
}

func (m *stdModule) initSecurityConfig(homePath string) error {
	securityConfig, err := newSecurityConfig(homePath)
	if err != nil {
		m.logger.Error("SettingsModule", "initSecurityConfig Error: "+err.Error())
	} else {
		err = m.securityConfigurer.Configure(securityConfig)
	}
	return err
}

func (m *stdModule) configure(securityCfg SecurityConfig, settings Settings, netIfaces []NetInterface, isSubNet bool) error {
	err := m.configureQuic(securityCfg, &settings, netIfaces)
	if err == nil {
		err = m.configureBroadcast(securityCfg, &settings, netIfaces, isSubNet)
	}
	if err == nil {
		err = m.configureRepository()
	}
	return err
}

func (m *stdModule) configureQuic(securityCfg SecurityConfig, settings *Settings, netIfaces []NetInterface) error {
	quicConfig := newQuicConfig(settings, securityCfg, netIfaces)
	return m.QuicConfigurer.Configure(quicConfig)
}

func (m *stdModule) configureBroadcast(securityCfg SecurityConfig, settings *Settings, netIfaces []NetInterface, isSubNet bool) error {
	broadcastConfig := newBroadcastConfig(settings, securityCfg, netIfaces, isSubNet)
	return m.BroadcastConfigurer.Configure(broadcastConfig)
}

func (m *stdModule) configureRepository() error {
	repositoryConfig := newRepositoryConfig(m.configurer.Config())
	return m.RepositoryConfigurer.Configure(repositoryConfig)
}

func (m *stdModule) loadSettings() (Settings, error) {
	settings, err := m.settingsSvc.Load()
	if err != nil {
		m.logger.Error("SettingsModule", "loadSettings Error: "+err.Error())
	}
	return settings, err
}

func (m *stdModule) loadSubModules() []interface{} {
	m.subModulesRW.RLock()
	defer m.subModulesRW.RUnlock()
	return slices.Clone(m.subModules)
}

type stdModuleSecurityConfigProxy struct {
	module *stdModule
}

var _ = (SecurityConfigurerListener)((*stdModuleSecurityConfigProxy)(nil))

func (p *stdModuleSecurityConfigProxy) OnConfigUpdated(cfg SecurityConfig) {
	module := p.module
	if module.ignoreConfigured {
		return
	}

	settings, err := module.loadSettings()
	if err != nil {
		return
	}

	module.configure(cfg, settings, nil, false)
}
