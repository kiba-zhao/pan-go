package user

import (
	"context"
	"pan/pkg/bootstrap"
	"pan/pkg/injection"
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

	return m
}

type stdModule struct {
	*module.BaseModule

	RepositoryManager repository.RepositoryManager

	repositoryBase *repository.RepositoryBase

	peerTopics    []net.PeerTopic
	peerTopicOnce sync.Once
}

var _ = (repository.RepositoryDBModule)((*stdModule)(nil))

func (module *stdModule) SetupToRepository(db repository.RepositoryDB) error {
	err := db.AutoMigrate(
		&User{},
		&UserConsensus{},
		&UserDevice{},
		&UserSecret{},
		&Passport{},
		&PassportUser{},
	)

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
			&UserSecretTopic{},
		}
	})
	return module.peerTopics
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (module *stdModule) Components() []injection.Component {
	components := []injection.Component{
		// repository
		injection.NewComponent(module.repositoryBase, injection.ComponentInternalScope),
		injection.NewComponent[UserDataRepository](&stdUserDataRepository{}, injection.ComponentInternalScope),
		injection.NewComponent[UserRepository](&stdUserRepository{}, injection.ComponentInternalScope),
		injection.NewComponent[UserConsensusRepository](&stdUserConsensusRepository{}, injection.ComponentInternalScope),
		injection.NewComponent[UserDeviceRepository](&stdUserDeviceRepository{}, injection.ComponentInternalScope),
		injection.NewComponent[UserSecretRepository](&stdUserSecretRepository{}, injection.ComponentInternalScope),

		// service
		injection.NewComponent(&UserDataService{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserConsensusService{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserDeviceService{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserSecretService{}, injection.ComponentInternalScope),

		// broker
		injection.NewComponent(&net.PeerBroker{PeerRouteModule: module}, injection.ComponentInternalScope),
		injection.NewComponent(&UserDataBroker{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserConsensusBroker{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserDeviceBroker{}, injection.ComponentInternalScope),
		injection.NewComponent(&UserSecretBroker{}, injection.ComponentInternalScope),
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
		defer module.RepositoryManager.DetachBaseModule(module)
	}
	return nil
}
