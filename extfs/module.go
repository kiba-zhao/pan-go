package extfs

import (
	"pan/app/bootstrap"
	"pan/app/peer"
	"pan/app/sample"
	"pan/app/web"

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
	return sample.New(&module{vfs: &vfs.VFS{}})
}

const moduleName = "extfs"

type module struct {
	DB                 sample.RepositoryDB
	SamplePeer         sample.SamplePeer
	controllers        []web.WebController
	controllersOnce    sync.Once
	peerAppModules     []peer.PeerAppModule
	peerAppModulesOnce sync.Once
	vfs                *vfs.VFS
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

func (m *module) Components() []bootstrap.Component {

	// base
	components := []bootstrap.Component{
		bootstrap.NewComponent(m.SamplePeer, bootstrap.ComponentInternalScope),
	}

	// services
	components = sample.AppendSampleInternalComponent[nodeitem.NodeItemInternalService](components, &nodeitem.NodeItemService{})
	components = sample.AppendSampleInternalComponent[nodefile.NodeFileInternalService](components, &nodefile.NodeFileService{})

	components = sample.AppendSampleComponent(components, &remotenode.RemoteNodeService{})
	components = sample.AppendSampleComponent(components, &remoteitem.RemoteItemService{})
	components = sample.AppendSampleComponent(components, &remotefile.RemoteFileService{})
	components = sample.AppendSampleComponent(components, &remoteblock.RemoteBlockService{})

	// repositories
	components = sample.AppendSampleComponent(components, nodeitem.NewNodeItemRepository(m.DB))

	// controllers
	for _, ctrl := range m.WebControllers() {
		components = append(components, bootstrap.NewComponent(ctrl, bootstrap.ComponentNoneScope))
	}

	// peer app modules
	for _, peerAppModule := range m.PeerAppModules() {
		components = append(components, bootstrap.NewComponent(peerAppModule, bootstrap.ComponentNoneScope))
	}

	// vfs components
	vfsComponents := m.vfs.VFSComponents()
	components = append(components, vfsComponents...)
	return components
}

func (m *module) Modules() []interface{} {
	return []interface{}{
		m.vfs,
	}
}
