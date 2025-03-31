package vfs

import (
	nodeitem "pan/features/extfs/node_item"
	remoteitem "pan/features/extfs/remote_item"
	remotenode "pan/features/extfs/remote_node"
)

type VFSFSRuntime struct {
	RemoteFileStreamServiceImpl *remoteitem.RemoteFileStreamService
	RemoteFileInfoServiceImpl   *remoteitem.RemoteFileInfoService
	RemoteItemServiceImpl       *remoteitem.RemoteItemService
	RemoteNodeServiceImpl       *remotenode.RemoteNodeService
	NodeItemServiceImpl         *nodeitem.NodeItemService
}

func (runtime *VFSFSRuntime) RemoteFileStreamService() *remoteitem.RemoteFileStreamService {
	return runtime.RemoteFileStreamServiceImpl
}

func (runtime *VFSFSRuntime) RemoteFileInfoService() *remoteitem.RemoteFileInfoService {
	return runtime.RemoteFileInfoServiceImpl
}

func (runtime *VFSFSRuntime) RemoteItemService() *remoteitem.RemoteItemService {
	return runtime.RemoteItemServiceImpl
}

func (runtime *VFSFSRuntime) RemoteNodeService() *remotenode.RemoteNodeService {
	return runtime.RemoteNodeServiceImpl
}

func (runtime *VFSFSRuntime) NodeItemService() *nodeitem.NodeItemService {
	return runtime.NodeItemServiceImpl
}
