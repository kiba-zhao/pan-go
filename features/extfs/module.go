package extfs

import (
	"pan/lib/feature"
	"pan/lib/injection"
	"pan/lib/repository"
	"pan/lib/runtime"

	remotesearchfile "pan/features/extfs/remote_search_file"
	"pan/features/extfs/vfs"
	"sync"

	nodeitem "pan/features/extfs/node_item"
	nodesearchfile "pan/features/extfs/node_search_file"
	remoteitem "pan/features/extfs/remote_item"
	remotenode "pan/features/extfs/remote_node"
	searchitem "pan/features/extfs/search_item"
)

const ModuleName = "extfs"

func New() interface{} {
	m := &module{}
	nodesearchfileModule := nodesearchfile.New(m)
	return runtime.NewModule(
		vfs.New(m),
		feature.New(ModuleName, nodesearchfileModule), nodesearchfileModule,
		feature.New(ModuleName, m), m,
	)
}

type module struct {
	injection.BaseComponentStoreProvider

	controllers     []feature.WebController
	controllersOnce sync.Once

	topics     []feature.PeerTopic
	topicsOnce sync.Once

	metaList     []feature.RepositoryMeta
	metaListOnce sync.Once

	brokerMetaList     []feature.BrokerMeta
	brokerMetaListOnce sync.Once
}

var _ = (feature.WebControllerProvider)((*module)(nil))

func (m *module) WebControllers() []feature.WebController {
	m.controllersOnce.Do(func() {

		m.controllers = []feature.WebController{
			// nodeitem controllers
			&nodeitem.NodeItemController{},
			&nodeitem.NodeFileInfoController{},
			&nodeitem.NodeFileStreamController{},

			// remotenode controllers
			&remotenode.RemoteNodeController{},
			&remoteitem.RemoteItemController{},
			&remoteitem.RemoteFileInfoController{},
			&remoteitem.RemoteFileStreamController{},

			// search controllers
			&searchitem.SearchItemController{},
			&remotesearchfile.RemoteSearchFileController{},
		}
	})
	return m.controllers
}

var _ = (repository.Repository)((*module)(nil))

func (m *module) SetupToRepository(db repository.RepositoryDB) error {
	return db.AutoMigrate(
		&nodeitem.NodeItem{},
		&searchitem.SearchItem{},
	)
}

var _ = (feature.RepositoryMetaProvider)((*module)(nil))

func (m *module) RepositoryMetaList() []feature.RepositoryMeta {
	m.metaListOnce.Do(func() {
		m.metaList = []feature.RepositoryMeta{
			feature.NewRepositoryMeta[nodeitem.NodeItemRepository](nodeitem.NewNodeItemRepository()),
			feature.NewRepositoryMeta[searchitem.SearchItemRepository](searchitem.NewSearchItemRepository()),
		}
	})
	return m.metaList
}

var _ = (feature.BrokerMetaProvider)((*module)(nil))

func (m *module) BrokerMetaList() []feature.BrokerMeta {
	m.brokerMetaListOnce.Do(func() {
		m.brokerMetaList = []feature.BrokerMeta{
			feature.NewBrokerMeta[*remoteitem.RemoteItemBroker](&remoteitem.RemoteItemBroker{}),
			feature.NewBrokerMeta[*remoteitem.RemoteFileInfoBroker](&remoteitem.RemoteFileInfoBroker{}),
			feature.NewBrokerMeta[*remoteitem.RemoteFileStreamBroker](&remoteitem.RemoteFileStreamBroker{}),
			feature.NewBrokerMeta[*remotesearchfile.RemoteSearchFileBroker](&remotesearchfile.RemoteSearchFileBroker{}),
		}
	})
	return m.brokerMetaList
}

var _ = (feature.PeerTopicProvider)((*module)(nil))

func (m *module) PeerTopics() []feature.PeerTopic {
	m.topicsOnce.Do(func() {
		m.topics = []feature.PeerTopic{
			&remoteitem.RemoteItemTopic{},
			&remoteitem.RemoteFileInfoTopic{},
			&remoteitem.RemoteFileStreamTopic{},
			&remotesearchfile.RemoteSearchFileTopic{},
		}
	})
	return m.topics
}

func (m *module) Components() []injection.Component {

	// base
	components := []injection.Component{}

	// node-items services
	components = feature.AppendInternalComponent[nodeitem.NodeItemInternalService](components, &nodeitem.NodeItemService{})
	components = feature.AppendInternalComponent[nodeitem.NodeFilePathInternalService](components, &nodeitem.NodeFilePathService{})
	components = feature.AppendInternalComponent[nodeitem.NodeFileInfoInternalService](components, &nodeitem.NodeFileInfoService{})

	// remote-nodes services
	components = feature.AppendComponent(components, &remotenode.RemoteNodeService{})
	components = feature.AppendComponent(components, &remoteitem.RemoteItemService{})
	components = feature.AppendInternalComponent[remoteitem.RemoteFileInfoInternalService](components, &remoteitem.RemoteFileInfoService{})
	components = feature.AppendComponent(components, &remoteitem.RemoteFileStreamService{})

	// search file services
	components = feature.AppendComponent(components, &searchitem.SearchItemService{})
	components = feature.AppendComponent(components, &remotesearchfile.RemoteSearchFileService{})

	return components
}
