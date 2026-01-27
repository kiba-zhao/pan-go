package user

import (
	"context"
	"pan/internal/bootstrap"
	"pan/internal/feature"
	"pan/internal/injection"
	"pan/internal/peer"
	"pan/internal/repository"
	sync "sync"
)

var subModuleNewFuncArray []feature.SubModuleNewFunc[*stdModule]

const (
	ModuleName = "user"
)

func New() interface{} {

	module := &stdModule{}
	module.featureRepository = &feature.Repository{}
	module.FeatureModule = feature.New(module, subModuleNewFuncArray...)

	return module
}

type stdModule struct {
	*feature.FeatureModule

	RepositoryManager repository.RepositoryManager

	featureRepository *feature.Repository

	peerTopics    []feature.PeerTopic
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
		err = module.featureRepository.SetupToRepository(db)
	}
	return err
}

func (module *stdModule) DBName() string {
	return ModuleName + ".db"
}

var _ = (feature.PeerAppModule)((*stdModule)(nil))

func (module *stdModule) PeerScope() []byte {
	return []byte(ModuleName)
}

func (module *stdModule) PeerTopics() []feature.PeerTopic {
	module.peerTopicOnce.Do(func() {
		module.peerTopics = []feature.PeerTopic{
			&UserDataTopic{},
			&UserConsensusTopic{},
			&UserDeviceTopic{},
			&UserSecretTopic{},
		}
	})
	return module.peerTopics
}

var _ = (peer.PeerAppModule)((*stdModule)(nil))

func (module *stdModule) SetupToPeer(router peer.PeerRouter) error {
	router_ := router.Route(module.PeerScope())

	topics := module.PeerTopics()
	for _, topic := range topics {
		if err := topic.SetupToPeer(router_); err != nil {
			return err
		}
	}
	return nil
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (module *stdModule) Components() []injection.Component {
	components := []injection.Component{
		// repository
		injection.NewComponent(module.featureRepository, injection.ComponentInternalScope),
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
		injection.NewComponent(&feature.PeerBroker{PeerAppModule: module}, injection.ComponentInternalScope),
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
