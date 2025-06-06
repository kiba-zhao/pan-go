package vfs

import (
	nodeitem "pan/features/extfs/node_item"
	remoteitem "pan/features/extfs/remote_item"
	remotenode "pan/features/extfs/remote_node"
)

type stdVFSFSRuntime struct {
	RemoteFileStreamServiceImpl *remoteitem.RemoteFileStreamService
	RemoteFileInfoServiceImpl   *remoteitem.RemoteFileInfoService
	RemoteItemServiceImpl       *remoteitem.RemoteItemService
	RemoteNodeServiceImpl       *remotenode.RemoteNodeService
	NodeItemServiceImpl         *nodeitem.NodeItemService
}

func (runtime *stdVFSFSRuntime) RemoteFileStreamService() *remoteitem.RemoteFileStreamService {
	return runtime.RemoteFileStreamServiceImpl
}

func (runtime *stdVFSFSRuntime) RemoteFileInfoService() *remoteitem.RemoteFileInfoService {
	return runtime.RemoteFileInfoServiceImpl
}

func (runtime *stdVFSFSRuntime) RemoteItemService() *remoteitem.RemoteItemService {
	return runtime.RemoteItemServiceImpl
}

func (runtime *stdVFSFSRuntime) RemoteNodeService() *remotenode.RemoteNodeService {
	return runtime.RemoteNodeServiceImpl
}

func (runtime *stdVFSFSRuntime) NodeItemService() *nodeitem.NodeItemService {
	return runtime.NodeItemServiceImpl
}
