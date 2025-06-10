// Define app module
package app

import (
	"context"
	appbroadcast "pan/features/app/broadcast"
	diskfile "pan/features/app/disk_file"
	appnode "pan/features/app/node"

	appsettings "pan/features/app/settings"
	"pan/lib/bootstrap"
	"pan/lib/broadcast"
	"pan/lib/config"
	"pan/lib/feature"
	"pan/lib/injection"
	"pan/lib/repository"
	"path/filepath"
	"sync"

	"pan/lib/peer"
	"pan/lib/quic"
	"pan/lib/runtime"
)

const ModuleName = "app"

func New(modules ...interface{}) interface{} {
	m := &module{}

	appSettingsService := appsettings.AppSettingsService{}
	m.appSettingsService = &appSettingsService

	// base modules
	modules_ := []interface{}{
		bootstrap.New(),
		config.NewWithDefaults(ModuleName+".toml", config.NewDefaultSettings()),
		repository.New(),
		peer.New(),
		broadcast.New(),
		quic.New(),
	}

	// specific modules
	if len(modules) > 0 {
		modules_ = append(modules_, modules...)
	}

	// feature modules
	modules_ = append(modules_, feature.New(ModuleName, m), m)
	return runtime.NewModule(modules_...)
}

// Bootstrap returns the bootstrap engine for the application.
//
// The bootstrap engine is used to initialize the application.
// It is a runtime.Engine, and can be used to mount components and modules into the runtime.
func Bootstrap() interface{} {
	return bootstrap.Bootstrap()
}

type module struct {
	injection.BaseComponentStoreProvider

	BroadcastModule broadcast.BroadcastModule
	BroadcastStore  broadcast.BroadcastStore

	PeerCluster peer.PeerCluster
	PeerGuard   peer.PeerGuard

	QuicExplorer      quic.QuicExplorer
	QuicExplorerGuide quic.QuicExplorerGuide

	AppConfig          config.AppConfig
	PeerConfig         peer.PeerConfig
	appSettingsService *appsettings.AppSettingsService

	controllers     []feature.WebController
	controllersOnce sync.Once

	metaList     []feature.RepositoryMeta
	metaListOnce sync.Once
}

var _ = (config.AppConfigListener)((*module)(nil))

func (m *module) OnConfigUpdated(settings config.AppSettings) {
	m.appSettingsService.SetConfigSettings(settings)
}

var _ = (peer.PeerConfigListener)((*module)(nil))

func (m *module) OnPeerConfigUpdated(settings *peer.PeerSettings) {
	peerId := peer.EncodePeerID(settings.PeerID())
	m.appSettingsService.SetPeerID(peerId)
}

var _ (bootstrap.DeferModule) = (*module)(nil)

func (m *module) Defer(ctx context.Context) error {
	m.BroadcastModule.SetStore(m.BroadcastStore)
	m.PeerCluster.RegisterPeerGuard(m.PeerGuard)
	m.QuicExplorer.AddGuide(m.QuicExplorerGuide)

	m.AppConfig.Subscribe(m)
	m.PeerConfig.Subscribe(m)

	rootPath := filepath.Dir(m.AppConfig.ConfigFilePath())
	m.appSettingsService.SetConfigPath(rootPath)
	return nil
}

var _ = (bootstrap.DestroyModule)((*module)(nil))

func (m *module) Destroy() {
	m.PeerCluster.UnregisterPeerGuard(m.PeerGuard)
	m.QuicExplorer.RemoveGuide(m.QuicExplorerGuide)

	m.AppConfig.Unsubscribe(m)
	m.PeerConfig.Unsubscribe(m)
}

var _ = (feature.WebControllerProvider)((*module)(nil))

// WebControllers initializes and returns a slice of WebControllers for the app module.
//
// It uses a sync.Once to ensure that the controllers are only initialized once.
// The method adds controllers for app node, disk file, and app settings to the
// controllers slice. This setup is essential for the web application to handle
// requests related to these components.
func (m *module) WebControllers() []feature.WebController {
	m.controllersOnce.Do(func() {
		m.controllers = []feature.WebController{
			&appnode.AppNodeController{},
			&diskfile.DiskFileController{},
			&appsettings.AppSettingsController{},
		}
	})
	return m.controllers
}

var _ = (repository.Repository)((*module)(nil))

func (m *module) SetupToRepository(db repository.RepositoryDB) error {
	return db.AutoMigrate(&appnode.AppNode{},
		&appnode.NetworkAddr{},
		&appbroadcast.AppBroadcastInfo{})
}

var _ = (feature.RepositoryMetaProvider)((*module)(nil))

func (m *module) RepositoryMetaList() []feature.RepositoryMeta {
	m.metaListOnce.Do(func() {
		m.metaList = []feature.RepositoryMeta{
			feature.NewRepositoryMeta[appnode.AppNodeRepository](appnode.NewAppNodeRepository()),
			feature.NewRepositoryMeta[appnode.NetworkAddrRepository](appnode.NewNetworkAddrRepository()),
			feature.NewRepositoryMeta[appbroadcast.AppBroadcastInfoRepository](appbroadcast.NewAppBroadcastInfoRepository()),
		}
	})

	return m.metaList
}

var _ = (injection.ComponentProvider)((*module)(nil))

// Components returns a slice of injection.Component representing the components
// provided by the app module. This includes the app module itself, the peer
// guard, services, repositories, controllers and the broadcast store.
func (m *module) Components() []injection.Component {
	// base
	components := []injection.Component{
		injection.NewComponent(m, injection.ComponentNoneScope),
	}

	// services
	components = feature.AppendComponent(components, &diskfile.DiskFileService{})
	components = feature.AppendExternalComponent[appsettings.AppSettingsExternalService](components, m.appSettingsService)
	components = feature.AppendExternalComponent[appnode.AppNodeExternalService](components, &appnode.AppNodeService{})
	components = feature.AppendComponent(components, &appnode.NetworkAddrService{})

	//  others
	components = feature.AppendInternalComponent[broadcast.BroadcastStore](components, &appbroadcast.BroadcastStore{})
	components = feature.AppendComponent[peer.PeerGuard](components, &appnode.PeerGuard{})
	components = feature.AppendComponent[quic.QuicExplorerGuide](components, &appnode.NetworkAddrGuide{})
	return components
}
