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
	"pan/pkg/ptp"
	"pan/pkg/repository"
	"sync"

	"github.com/spf13/viper"
)

var subModuleNewFuncArray []module.SubModuleNewFunc[*stdModule]

var ErrSettingsModuleUavailable = errors.New("settings.Module Error: Unavailable")
var ErrSettingsModuleHomePathConflict = errors.New("settings.Module Error: Home Path Conflict")

type ModuleNetProvider interface {
	NewNetConfig(deviceNetwork *DeviceNetwork) (ptp.QuicConfig, ptp.BroadcastConfig)
}

const (
	SettingsModuleName = "settings"
)

type stdModule struct {
	QuicConfigurer       ptp.QuicConfigurer
	BroadcastConfigurer  ptp.BroadcastConfigurer
	RepositoryConfigurer repository.RepositoryConfigurer

	*module.BaseModule

	logger log.Logger

	deviceInfoSvc      *DeviceInfoService
	deviceNetworkSvc   *DeviceNetworkService
	configurer         SettingsConfigurer
	securityConfigurer SecurityConfigurer

	viper *viper.Viper

	locker sync.Mutex

	securityConfigListener SecurityConfigurerListener
	netProvider            ModuleNetProvider
}

func New() interface{} {

	viper := viper.New()

	deviceInfoSvc := &DeviceInfoService{}
	deviceNetworkSvc := &DeviceNetworkService{}

	logger := log.Default()
	configurer := config.NewConfigurer[SettingsConfig](logger)
	securityConfigurer := config.NewConfigurer[SecurityConfig](logger)

	m := &stdModule{}
	m.logger = logger
	m.configurer = configurer
	m.securityConfigurer = securityConfigurer
	m.viper = viper
	m.deviceInfoSvc = deviceInfoSvc
	m.deviceNetworkSvc = deviceNetworkSvc

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

	m.deviceInfoSvc.setup(cfg)

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
		injection.NewComponent[DeviceNetworkChangedTrigger](m, injection.ComponentInternalScope),
		injection.NewComponent(m.viper, injection.ComponentInternalScope),
		// configurer
		injection.NewComponent(m.configurer, injection.ComponentExternalScope),
		injection.NewComponent(m.securityConfigurer, injection.ComponentExternalScope),
		// service
		injection.NewComponent(m.deviceInfoSvc, injection.ComponentInternalScope),
		injection.NewComponent(m.deviceNetworkSvc, injection.ComponentInternalScope),
		// injection.NewComponent[SettingsExternalService](m.settingsSvc, injection.ComponentExternalScope),
	}
}

var _ = (DeviceNetworkChangedTrigger)((*stdModule)(nil))

func (m *stdModule) OnDeviceNetworkChanged(deviceNetwork DeviceNetwork) {
	m.configure(&deviceNetwork)
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

func (m *stdModule) configure(deviceNetwork *DeviceNetwork) error {

	var err error
	var deviceNetwork_ *DeviceNetwork
	if deviceNetwork != nil {
		deviceNetwork_ = deviceNetwork
	} else {
		deviceNetwork_, err = m.loadDeviceNetwork()
		if err != nil {
			return err
		}
	}

	var quicConfig ptp.QuicConfig
	var broadcastConfig ptp.BroadcastConfig
	if m.netProvider != nil {
		quicConfig, broadcastConfig = m.netProvider.NewNetConfig(deviceNetwork_)
	}

	if quicConfig == nil {
		quicConfig = newQuicConfig(deviceNetwork_, m.securityConfigurer.Config(), []string{""})
	}

	if broadcastConfig == nil {
		broadcastConfig = newBroadcastConfig(deviceNetwork_, m.configurer.Config())
	}

	err = m.QuicConfigurer.Configure(quicConfig)
	if err == nil {
		err = m.BroadcastConfigurer.Configure(broadcastConfig)
	} else {
		m.BroadcastConfigurer.Configure(broadcastConfig)
	}

	return err
}

func (m *stdModule) loadDeviceNetwork() (*DeviceNetwork, error) {
	deviceNetwork, err := m.deviceNetworkSvc.Load()
	if err != nil {
		m.logger.Error("settings.SettingsModule", "loadDeviceNetwork Error: "+err.Error())
	}
	return &deviceNetwork, err
}

type stdModuleSecurityConfigProxy struct {
	module *stdModule
}

var _ = (SecurityConfigurerListener)((*stdModuleSecurityConfigProxy)(nil))

func (p *stdModuleSecurityConfigProxy) OnConfigUpdated(cfg SecurityConfig) {
	if p.module != nil {
		deviceNetwork, err := p.module.loadDeviceNetwork()
		if err != nil {
			return
		}
		p.module.configure(deviceNetwork)
	}
}
