package extfs

import (
	"pan/app/injection"
	"pan/app/peer"
	"pan/app/sample"
	"pan/app/web"
	"pan/runtime"

	"pan/extfs/vfs"
	"sync"

	nodefile "pan/extfs/node_file"
	nodeitem "pan/extfs/node_item"
	remoteblock "pan/extfs/remote_block"
	remotefile "pan/extfs/remote_file"
	remoteitem "pan/extfs/remote_item"
	remotenode "pan/extfs/remote_node"
)

func New() interface{} {
	m := &module{}
	m.store = injection.NewComponentStore()

	sampleModule := sample.New(m)
	m.sample = sampleModule

	return runtime.NewModule(vfs.New(m.store), sampleModule)
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
			&nodeitem.NodeItemController{},
			&nodefile.NodeFileController{},
			&remotenode.RemoteNodeController{},
			&remoteitem.RemoteItemController{},
			&remotefile.RemoteFileController{},
		}
	})
	return m.controllers
}

func (m *module) PeerAppModules() []peer.PeerAppModule {
	m.peerAppModulesOnce.Do(func() {
		m.peerAppModules = []peer.PeerAppModule{
			&remoteitem.RemoteItemTopic{},
			&remotefile.RemoteFileTopic{},
			&remoteblock.RemoteBlockTopic{},
		}
	})
	return m.peerAppModules
}

func (m *module) Models() []interface{} {
	return []interface{}{
		&nodeitem.NodeItem{},
	}
}

func (m *module) ComponentStore() injection.ComponentStore {
	return m.store
}

func (m *module) Components() []injection.Component {

	// base
	components := []injection.Component{}

	// services
	components = sample.AppendSampleInternalComponent[nodeitem.NodeItemInternalService](components, &nodeitem.NodeItemService{})
	components = sample.AppendSampleInternalComponent[nodefile.NodeFileInternalService](components, &nodefile.NodeFileService{})

	components = sample.AppendSampleComponent(components, &remotenode.RemoteNodeService{})
	components = sample.AppendSampleComponent(components, &remoteitem.RemoteItemService{})
	components = sample.AppendSampleComponent(components, &remotefile.RemoteFileService{})
	components = sample.AppendSampleComponent(components, &remoteblock.RemoteBlockService{})

	// brokers
	components = sample.AppendSampleComponent(components, &remoteitem.RemoteItemBroker{SamplePeer: m.sample})
	components = sample.AppendSampleComponent(components, &remotefile.RemoteFileBroker{SamplePeer: m.sample})
	components = sample.AppendSampleComponent(components, &remoteblock.RemoteBlockBroker{SamplePeer: m.sample})

	// repositories
	components = sample.AppendSampleComponent(components, nodeitem.NewNodeItemRepository(m.sample.DB()))

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
