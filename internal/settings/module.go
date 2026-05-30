package settings

import (
	"context"
	"errors"
	"pan/internal/bootstrap"
	"pan/internal/config"
	"pan/internal/injection"
	"pan/internal/log"
	"pan/internal/module"
	"pan/internal/net"
	"pan/internal/repository"
	"sync"

	"github.com/spf13/viper"
)

var subModuleNewFuncArray []module.SubModuleNewFunc[*stdModule]

var ErrSettingsModuleUavailable = errors.New("settings.Module Error: Unavailable")

const (
	SettingsModuleName = "settings"
)

type stdModule struct {
	QuicConfigurer       net.QuicConfigurer
	BroadcastConfigurer  net.BroadcastConfigurer
	RepositoryConfigurer repository.RepositoryConfigurer

	*module.BaseModule

	logger log.Logger

	settingsSvc        *SettingsService
	configurer         SettingsConfigurer
	securityConfigurer SecurityConfigurer

	viper *viper.Viper

	locker sync.Mutex

	ignoreConfigured bool

	securityConfigListener SecurityConfigurerListener
}

func New() interface{} {

	viper := viper.New()

	settingsSvc := &SettingsService{}

	logger := log.Default()
	configurer := config.NewConfigurer[SettingsConfig](logger)
	securityConfigurer := config.NewConfigurer[SecurityConfig](logger)

	m := &stdModule{}
	m.logger = logger
	m.configurer = configurer
	m.securityConfigurer = securityConfigurer
	m.viper = viper
	m.settingsSvc = settingsSvc

	m.BaseModule = module.New(m, subModuleNewFuncArray...)

	securityConfigListener := &stdModuleSecurityConfigProxy{}
	securityConfigListener.module = m
	m.securityConfigListener = securityConfigListener

	return m
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
		injection.NewComponent[SettingsExternalService](m.settingsSvc, injection.ComponentExternalScope),
	}
}

var _ = (SettingsChangedTrigger)((*stdModule)(nil))

func (m *stdModule) OnSettingsChanged(settings Settings) {
	if !m.ignoreConfigured {
		securityCfg := m.securityConfigurer.Config()
		m.configure(securityCfg, settings, nil, false)
		return
	}

	subModules := m.Modules()
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
		err = m.configureBroadcast(&settings, netIfaces, isSubNet)
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

func (m *stdModule) configureBroadcast(settings *Settings, netIfaces []NetInterface, isSubNet bool) error {
	broadcastConfig := newBroadcastConfig(settings, netIfaces, isSubNet)
	return m.BroadcastConfigurer.Configure(broadcastConfig)
}

func (m *stdModule) configureRepository() error {
	repositoryConfig := newRepositoryConfig(m.configurer.Config())
	return m.RepositoryConfigurer.Configure(repositoryConfig)
}

func (m *stdModule) loadSettings() (Settings, error) {
	settings, err := m.settingsSvc.Load()
	if err != nil {
		m.logger.Error("settings.SettingsModule", "loadSettings Error: "+err.Error())
	}
	return settings, err
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
