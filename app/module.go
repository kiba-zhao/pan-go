// Define app module
package app

import (
	appbroadcast "pan/app/app_broadcast"
	appnode "pan/app/app_node"
	appsettings "pan/app/app_settings"
	"pan/app/bootstrap"
	"pan/app/broadcast"
	"pan/app/config"
	diskfile "pan/app/disk_file"
	"pan/app/injection"
	"path/filepath"

	"pan/app/guard"
	"pan/app/peer"
	"pan/app/quic"
	"pan/app/sample"
	"pan/app/web"
	"pan/runtime"
	"sync"
)

// New returns the app module.
//
// The app module is the core module for Pan.
// It contains all the components and modules for the application.
// The module is also a runtime.Module, and can be used to load components into the runtime.
//
// The module is initialized with the given ComponentStoreProvider, which is used to get the ComponentStore.
// The ComponentStore is used to store components that are injected into other components.
func New() interface{} {
	m := &module{}
	m.peerGuard = &guard.PeerGuard{}
	m.store = injection.NewComponentStore()

	sampleModule := sample.New(m)
	m.sample = sampleModule

	return runtime.NewModule(bootstrap.New(), config.New(), peer.New(), broadcast.New(m), quic.New(), web.New(), sampleModule)
}

// Bootstrap returns the bootstrap engine for the application.
//
// The bootstrap engine is used to initialize the application.
// It is a runtime.Engine, and can be used to mount components and modules into the runtime.
func Bootstrap() interface{} {
	return bootstrap.Bootstrap()
}

const moduleName = "app"

type module struct {
	PeerModule      peer.PeerModule
	Config          config.AppConfig
	store           injection.ComponentStore
	sample          sample.Sample
	settings        config.AppSettings
	settingsRW      sync.RWMutex
	controllers     []web.WebController
	controllersOnce sync.Once
	peerGuard       *guard.PeerGuard
}

// Name returns the name of the module, which is "app".
// It implements simple.SampleProvider
func (m *module) Name() string {
	return moduleName
}

// WebControllers initializes and returns a slice of WebControllers for the app module.
//
// It uses a sync.Once to ensure that the controllers are only initialized once.
// The method adds controllers for app node, disk file, and app settings to the
// controllers slice. This setup is essential for the web application to handle
// requests related to these components.

func (m *module) WebControllers() []web.WebController {
	m.controllersOnce.Do(func() {
		// TODO: add web and node controllers
		m.controllers = []web.WebController{
			&appnode.AppNodeController{},
			&diskfile.DiskFileController{},
			&appsettings.AppSettingsController{},
		}
	})
	return m.controllers
}

// Models returns a slice of interfaces representing the models used by the app module.
// It includes the AppNode and AppBroadcastInfo models, which are essential
// for the node and broadcast functionalities within the application.

func (m *module) Models() []interface{} {
	return []interface{}{
		&appnode.AppNode{},
		&appbroadcast.AppBroadcastInfo{},
	}
}

// ComponentStore returns the ComponentStore associated with the app module.
//
// The ComponentStore is a map of components that are injected into other components.
// It is used to store components that are managed by the app module.
func (m *module) ComponentStore() injection.ComponentStore {
	return m.store
}

// Components returns a slice of injection.Component representing the components
// provided by the app module. This includes the app module itself, the peer
// guard, services, repositories, controllers and the broadcast store.
func (m *module) Components() []injection.Component {
	// base
	components := []injection.Component{
		injection.NewComponent(m, injection.ComponentNoneScope),
		// submodules
		injection.NewComponent(m.peerGuard, injection.ComponentNoneScope),
	}

	// services
	components = sample.AppendSampleComponent(components, &diskfile.DiskFileService{})
	components = sample.AppendSampleExternalComponent[appsettings.AppSettingsExternalService](components, &appsettings.AppSettingsService{Provider: m})
	components = sample.AppendSampleExternalComponent[appnode.AppNodeExternalService](components, &appnode.AppNodeService{})

	// repositories
	components = sample.AppendSampleComponent(components, appnode.NewAppNodeRepository(m.sample.DB()))
	components = sample.AppendSampleComponent(components, appbroadcast.NewAppBroadcastInfoRepository(m.sample.DB()))

	// controllers
	for _, ctrl := range m.WebControllers() {
		components = append(components, injection.NewComponent(ctrl, injection.ComponentNoneScope))
	}

	//  store
	components = sample.AppendSampleInternalComponent[broadcast.BroadcastStore](components, &appbroadcast.BroadcastStore{})
	return components
}

// OnConfigUpdated is called when the configuration is updated.
//
// It updates the configuration used by the app module.
func (m *module) OnConfigUpdated(settings config.AppSettings) {
	m.settingsRW.Lock()
	defer m.settingsRW.Unlock()
	m.settings = settings
}

// Settings returns a copy of the current configuration settings.
//
// The method returns a copy of the configuration settings used by the app module.
// The returned settings are a snapshot of the current configuration and may be
// outdated if the configuration is changed after calling this method.
// The method is thread-safe and may be called concurrently.
func (m *module) Settings() config.Settings {
	m.settingsRW.RLock()
	defer m.settingsRW.RUnlock()

	settings := *m.settings
	return settings
}

// SetSettings sets the given configuration settings to the underlying storage.
//
// It saves the given settings to the file specified by ConfigFilePath and
// notifies all registered config listeners about the updates.
// An error is returned if the saving process fails.
func (m *module) SetSettings(settings config.Settings) error {
	return m.Config.Save(&settings)
}

// RootPath returns the root path of the application.
//
// The root path is the path containing the configuration file. It is determined by
// the environment variable "rootPath" if set, otherwise it defaults to the user's
// home directory with the package name as a suffix.
func (m *module) RootPath() string {
	return filepath.Dir(m.Config.ConfigFilePath())
}

// PeerID returns the peer ID of the peer module, or an empty string if the peer
// module is not available or the peer settings are not available.
//
// The peer ID is the marshaled bytes of the public key of the peer's certificate.
func (m *module) PeerID() string {
	if m.PeerModule == nil {
		return ""
	}

	settings := m.PeerModule.PeerSettings()
	if settings == nil || !settings.Available() {
		return ""
	}

	return appnode.EncodePeerID(settings.PeerID())
}

// Modules returns a slice of interfaces representing the sub-modules
// of the app module. It currently includes the peer guard, which
// provides access control and security features within the p2p module.

func (m *module) Modules() []interface{} {
	return []interface{}{m.peerGuard}
}
