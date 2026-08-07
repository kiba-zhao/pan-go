package user

import (
	"context"
	"pan/internal/settings"
	"pan/pkg/bootstrap"
	"pan/pkg/injection"
	"pan/pkg/log"
	"pan/pkg/module"
	"pan/pkg/net"
	"pan/pkg/repository"
	sync "sync"
)

var subModuleNewFuncArray []module.SubModuleNewFunc[*stdModule]

const (
	ModuleName = "user"
)

func New() interface{} {

	m := &stdModule{}
	m.repositoryBase = &repository.RepositoryBase{}
	m.BaseModule = module.New(m, subModuleNewFuncArray...)

	userDeviceSyncAgent := &UserDeviceSyncAgent{}
	userDeviceSyncAgent.syncCh = make(chan struct{}, 1)
	userDeviceSyncAgent.syncErrs = make([]*UserDeviceSyncErr, 0)
	m.userDeviceSyncAgent = userDeviceSyncAgent

	logger := log.Default()
	userDeviceSyncAgent.logger = logger

	return m
}

type stdModule struct {
	*module.BaseModule

	SecurityConfigurer settings.SecurityConfigurer
	RepositoryManager  repository.RepositoryManager

	repositoryBase *repository.RepositoryBase

	peerTopics    []net.PeerTopic
	peerTopicOnce sync.Once

	userDeviceSyncAgent *UserDeviceSyncAgent
}

var _ = (repository.RepositoryDBModule)((*stdModule)(nil))

func (module *stdModule) SetupToRepository(db repository.RepositoryDB) error {
	var err error
	if db != nil {
		err = db.AutoMigrate(
			&User{},
			&UserConsensus{},
			&UserDevice{},
			&UserExtra{},
			&UserSecret{},
			&Passport{},
			&PassportUser{},
		)
	}

	if err == nil {
		err = module.repositoryBase.SetupToRepository(db)
	}
	return err
}

func (module *stdModule) DBName() string {
	return ModuleName + ".db"
}

var _ = (net.PeerRouteModule)((*stdModule)(nil))

func (module *stdModule) PeerRouteScope() net.PeerServletScope {
	return []byte(ModuleName)
}

var _ = (net.PeerTopicProvider)((*stdModule)(nil))

func (module *stdModule) PeerTopics() []net.PeerTopic {
	module.peerTopicOnce.Do(func() {
		module.peerTopics = []net.PeerTopic{
			&UserDataTopic{},
			&UserConsensusTopic{},
			&UserDeviceTopic{},
			&UserExtraTopic{},
			&UserSecretTopic{},
		}
	})
	return module.peerTopics
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (module *stdModule) Components() []injection.Component {
	components := []injection.Component{
		injection.NewComponent(module, injection.ComponentInternalScope),
		// repository
		injection.NewComponent(module.repositoryBase, injection.ComponentInternalScope),
		injection.NewComponent[UserRepository](&stdUserRepository{}, injection.ComponentInternalScope),
		injection.NewComponent[UserDataRepository](&stdUserDataRepository{}, injection.ComponentInternalScope),
		injection.NewComponent[UserConsensusRepository](&stdUserConsensusRepository{}, injection.ComponentInternalScope),
		injection.NewComponent[UserDeviceRepository](&stdUserDeviceRepository{}, injection.ComponentInternalScope),
		injection.NewComponent[UserExtraRepository](&stdUserExtraRepository{}, injection.ComponentInternalScope),
		injection.NewComponent[UserSecretRepository](&stdUserSecretRepository{}, injection.ComponentInternalScope),

		// service
		injection.NewComponent(&UserService{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserDataService{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserConsensusService{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserDeviceService{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserExtraService{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserSecretService{}, injection.ComponentInternalScope),

		// broker
		injection.NewComponent(&net.PeerBroker{PeerRouteModule: module}, injection.ComponentInternalScope),
		injection.NewComponent(&UserDataBroker{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserConsensusBroker{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserDeviceBroker{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserExtraBroker{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserSecretBroker{}, injection.ComponentInternalScope),

		// others
		injection.NewComponent(module.userDeviceSyncAgent, injection.ComponentInternalScope),
	}

	// topics
	topics := module.PeerTopics()
	for _, topic := range topics {
		components = append(components, injection.NewComponent(topic, injection.ComponentNoneScope))
	}

	return components
}

var _ = (bootstrap.ReadyModule)((*stdModule)(nil))

func (module *stdModule) Ready(ctx context.Context) error {
	if module.RepositoryManager != nil {
		module.RepositoryManager.AttachBaseModule(module)
		module.RepositoryManager.DetachBaseModule(module)
	}

	if module.SecurityConfigurer != nil {
		module.SecurityConfigurer.Subscribe(module)
		defer module.SecurityConfigurer.Unsubscribe(module)
	}

	return module.userDeviceSyncAgent.doSync(ctx)
}

var _ = (net.PeerServerListener)((*stdModule)(nil))

func (module *stdModule) OnServePeerConn(ctx context.Context, conn net.PeerConn) error {
	return module.userDeviceSyncAgent.update(ctx, conn.PeerID())
}

func (module *stdModule) OnClosePeerConn(ctx context.Context, conn net.PeerConn) error {
	return module.userDeviceSyncAgent.purge(ctx, conn.PeerID())
}

var _ = (settings.SecurityConfigurerListener)((*stdModule)(nil))

func (module *stdModule) OnConfigUpdated(config settings.SecurityConfig) {
	userConfig := NewUserConfig(config)
	module.userDeviceSyncAgent.setup(userConfig)
}
