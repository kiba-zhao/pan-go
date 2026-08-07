package settings

import (
	"context"
	"errors"
	"os"
	"pan/pkg/bootstrap"
	"pan/pkg/config"
	"pan/pkg/injection"
	"pan/pkg/log"
	"pan/pkg/module"
	"pan/pkg/net"
	"pan/pkg/repository"
	"sync"

	"github.com/spf13/viper"
)

var subModuleNewFuncArray []module.SubModuleNewFunc[*stdModule]

var ErrSettingsModuleUavailable = errors.New("settings.Module Error: Unavailable")
var ErrSettingsModuleHomePathConflict = errors.New("settings.Module Error: Home Path Conflict")

type ModuleNetProvider interface {
	NewNetConfig(settings *Settings) (net.QuicConfig, net.BroadcastConfig)
}

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

	securityConfigListener SecurityConfigurerListener
	netProvider            ModuleNetProvider
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

	homePath := cfg.HomePath()
	stat, err := os.Stat(homePath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(homePath, 0755)
	} else if err == nil && !stat.IsDir() {
		err = ErrSettingsModuleHomePathConflict
	}

	if err != nil {
		m.logger.Error("SettingsModule", "OnConfigUpdated Error: "+err.Error())
		panic(err)
	}

	m.settingsSvc.setup(cfg)

	// init viper and securityConfig

	err = m.initViper(homePath)
	if err == nil {
		err = m.initSecurityConfig(homePath)
	}
	if err == nil {
		err = m.initRepository(cfg)
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
	m.configure(&settings)
}

func (m *stdModule) initViper(homePath string) error {
	err := initViper(m.viper, homePath)

	if err != nil {
		if os.IsNotExist(err) {
			m.logger.Warn("SettingsModule", "initViper Warning: "+err.Error())
			err = nil
		} else {
			m.logger.Error("SettingsModule", "initViper Error: "+err.Error())
		}
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

func (m *stdModule) initRepository(settingsConfig SettingsConfig) error {
	repositoryConfig := newRepositoryConfig(settingsConfig)
	return m.RepositoryConfigurer.Configure(repositoryConfig)
}

func (m *stdModule) configure(settings *Settings) error {

	var err error
	var settings_ *Settings
	if settings != nil {
		settings_ = settings
	} else {
		settings_, err = m.loadSettings()
		if err != nil {
			return err
		}
	}

	var quicConfig net.QuicConfig
	var broadcastConfig net.BroadcastConfig
	if m.netProvider != nil {
		quicConfig, broadcastConfig = m.netProvider.NewNetConfig(settings_)
	}

	if quicConfig == nil {
		quicConfig = newQuicConfig(settings_, m.securityConfigurer.Config(), []string{""})
	}

	if broadcastConfig == nil {
		broadcastConfig = newBroadcastConfig(settings_, m.configurer.Config())
	}

	err = m.QuicConfigurer.Configure(quicConfig)
	if err == nil {
		err = m.BroadcastConfigurer.Configure(broadcastConfig)
	} else {
		m.BroadcastConfigurer.Configure(broadcastConfig)
	}

	return err
}

func (m *stdModule) loadSettings() (*Settings, error) {
	settings, err := m.settingsSvc.Load()
	if err != nil {
		m.logger.Error("settings.SettingsModule", "loadSettings Error: "+err.Error())
	}
	return &settings, err
}

type stdModuleSecurityConfigProxy struct {
	module *stdModule
}

var _ = (SecurityConfigurerListener)((*stdModuleSecurityConfigProxy)(nil))

func (p *stdModuleSecurityConfigProxy) OnConfigUpdated(cfg SecurityConfig) {
	if p.module != nil {
		settings, err := p.module.loadSettings()
		if err != nil {
			return
		}
		p.module.configure(settings)
	}
}
