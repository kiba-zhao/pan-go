// Define fuse for remote item
package remoteitem

import (
	remotenode "pan/extfs/remote_node"
)

type VFSFUSERemoteItemRuntime interface {

	// It extends the VFSFUSERemoteFileInfoServiceFUSEruntime interface.
	VFSFUSERemoteFileRuntime

	RemoteNodeService() *remotenode.RemoteNodeService
	// Returns the remote item service.
	RemoteItemService() *RemoteItemService
}
