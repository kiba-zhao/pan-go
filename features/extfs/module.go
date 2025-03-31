package extfs

import (
	"pan/lib/injection"
	"pan/lib/peer"
	"pan/lib/runtime"
	"pan/lib/sample"
	"pan/lib/web"

	remotesearchfile "pan/features/extfs/remote_search_file"
	"pan/features/extfs/vfs"
	"sync"

	nodeitem "pan/features/extfs/node_item"
	nodesearchfile "pan/features/extfs/node_search_file"
	remoteitem "pan/features/extfs/remote_item"
	remotenode "pan/features/extfs/remote_node"
	searchitem "pan/features/extfs/search_item"
)

func New() interface{} {
	m := &module{}
	m.store = injection.NewComponentStore()

	sampleModule := sample.New(m)
	m.sample = sampleModule

	return runtime.NewModule(vfs.New(m), nodesearchfile.New(m), sampleModule)
}

const moduleName = "extfs"

type module struct {
	store              injection.ComponentStore
	sample             sample.Sample
	controllers        []web.WebController
	controllersOnce    sync.Once
	peerAppModules     []peer.PeerAppModule
	peerAppModulesOnce sync.Once
}

func (m *module) Name() string {
	return moduleName
}

func (m *module) WebControllers() []web.WebController {
	m.controllersOnce.Do(func() {

		m.controllers = []web.WebController{
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
			&nodesearchfile.NodeSearchFileController{},
			&remotesearchfile.RemoteSearchFileController{},
		}
	})
	return m.controllers
}

func (m *module) PeerAppModules() []peer.PeerAppModule {
	m.peerAppModulesOnce.Do(func() {
		m.peerAppModules = []peer.PeerAppModule{
			&remoteitem.RemoteItemTopic{},
			&remoteitem.RemoteFileInfoTopic{},
			&remoteitem.RemoteFileStreamTopic{},
			&remotesearchfile.RemoteSearchFileTopic{},
		}
	})
	return m.peerAppModules
}

func (m *module) Models() []interface{} {
	return []interface{}{
		&nodeitem.NodeItem{},
		&searchitem.SearchItem{},
	}
}

func (m *module) ComponentStore() injection.ComponentStore {
	return m.store
}

func (m *module) Components() []injection.Component {

	// base
	components := []injection.Component{}

	// node-items services
	components = sample.AppendSampleInternalComponent[nodeitem.NodeItemInternalService](components, &nodeitem.NodeItemService{})
	components = sample.AppendSampleInternalComponent[nodeitem.NodeFilePathInternalService](components, &nodeitem.NodeFilePathService{})
	components = sample.AppendSampleInternalComponent[nodeitem.NodeFileInfoInternalService](components, &nodeitem.NodeFileInfoService{})

	// remote-nodes services
	components = sample.AppendSampleComponent(components, &remotenode.RemoteNodeService{})
	components = sample.AppendSampleComponent(components, &remoteitem.RemoteItemService{})
	components = sample.AppendSampleInternalComponent[remoteitem.RemoteFileInfoInternalService](components, &remoteitem.RemoteFileInfoService{})
	components = sample.AppendSampleComponent(components, &remoteitem.RemoteFileStreamService{})

	// search file services
	components = sample.AppendSampleInternalComponent[nodesearchfile.NodeSearchFileInternalService](components, &nodesearchfile.NodeSearchFileService{})
	components = sample.AppendSampleComponent(components, &searchitem.SearchItemService{})
	components = sample.AppendSampleComponent(components, &remotesearchfile.RemoteSearchFileService{})

	// brokers
	components = sample.AppendSampleComponent(components, &remoteitem.RemoteItemBroker{SamplePeer: m.sample})
	components = sample.AppendSampleComponent(components, &remoteitem.RemoteFileInfoBroker{SamplePeer: m.sample})
	components = sample.AppendSampleComponent(components, &remoteitem.RemoteFileStreamBroker{SamplePeer: m.sample})
	components = sample.AppendSampleComponent(components, &remotesearchfile.RemoteSearchFileBroker{SamplePeer: m.sample})

	// repositories
	components = sample.AppendSampleComponent(components, nodeitem.NewNodeItemRepository(m.sample.DB()))
	components = sample.AppendSampleComponent(components, searchitem.NewSearchItemRepository(m.sample.DB()))

	// controllers
	for _, ctrl := range m.WebControllers() {
		components = append(components, injection.NewComponent(ctrl, injection.ComponentNoneScope))
	}

	// peer app modules
	for _, peerAppModule := range m.PeerAppModules() {
		components = append(components, injection.NewComponent(peerAppModule, injection.ComponentNoneScope))
	}

	return components
}
