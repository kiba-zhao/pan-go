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
	"pan/lib/peer"
	"pan/lib/quic"
	"pan/lib/repository"
	"pan/lib/runtime"
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

	settingsSvc *SettingsService
	configurer  SettingsConfigurer

	viper        *viper.Viper
	peerSecurity peer.PeerSecurity

	componentStore     injection.ComponentStore
	componentStoreOnce sync.Once

	locker      sync.Mutex
	homePath    string
	settingsCfg SettingsConfig

	configListeners []SettingsConfigListener

	isMobileMode bool
}

func New() interface{} {

	viper := viper.New()

	settingsSvc := &SettingsService{}

	logger := log.Default()
	configurer := config.NewConfigurer[SettingsConfig](logger)

	module := &stdModule{}
	module.logger = logger
	module.configurer = configurer
	module.viper = viper
	module.settingsSvc = settingsSvc

	return module
}

var _ = (SettingsConfigListener)((*stdModule)(nil))

func (m *stdModule) OnConfigUpdated(cfg SettingsConfig) {
	m.locker.Lock()
	defer m.locker.Unlock()

	// init viper and peerSecurity if homePath changed
	homePath := cfg.HomePath()
	if homePath != m.homePath {
		m.homePath = homePath
		err := m.initViper()
		if err == nil {
			err = m.initPeerSecurity()
		}
		if err != nil {
			return
		}
	}
	//
	m.settingsSvc.Setup(cfg)

	// init Configurer if not mobile mode
	m.settingsCfg = cfg
	if !m.isMobileMode {
		err := m.configure(nil, false)
		if err != nil {
			return
		}
	}
	//

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

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (m *stdModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentInternalScope),
		injection.NewComponent(m.viper, injection.ComponentInternalScope),
		// configurer
		injection.NewComponent(m.configurer, injection.ComponentExternalScope),
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

func (m *stdModule) initViper() error {
	err := initViper(m.viper, m.homePath)
	if err != nil {
		m.logger.Error("SettingsModule", "initViper Error: "+err.Error())
	}
	return err
}

func (m *stdModule) initPeerSecurity() error {
	peerSecurity, err := peer.NewPeerSecurity(m.homePath)
	if err != nil {
		m.logger.Error("SettingsModule", "initPeerSecurity Error: "+err.Error())
	} else {
		m.peerSecurity = peerSecurity
	}
	return err
}

func (m *stdModule) configure(netIfaces []NetInterface, isSubNet bool) error {
	settings, err := m.settingsSvc.Load()
	if err != nil {
		m.logger.Error("SettingsModule", "configure Error: load failed"+err.Error())
		return err
	}

	err = m.configureQuic(&settings, netIfaces)
	if err == nil {
		err = m.configureBroadcast(&settings, netIfaces, isSubNet)
	}
	if err == nil {
		err = m.configureRepository()
	}
	return err
}

func (m *stdModule) configureQuic(settings *Settings, netIfaces []NetInterface) error {
	quicConfig := newQuicConfig(settings, m.peerSecurity, netIfaces)
	return m.QuicConfigurer.Configure(quicConfig)
}

func (m *stdModule) configureBroadcast(settings *Settings, netIfaces []NetInterface, isSubNet bool) error {
	broadcastConfig := newBroadcastConfig(settings, m.peerSecurity, netIfaces, isSubNet)
	return m.BroadcastConfigurer.Configure(broadcastConfig)
}

func (m *stdModule) configureRepository() error {
	repositoryConfig := newRepositoryConfig(m.settingsCfg)
	return m.RepositoryConfigurer.Configure(repositoryConfig)
}
