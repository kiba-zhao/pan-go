package impl

import (
	"pan/modules/extfs/services"
	"pan/peer"
)

type remoteFilesStateEventImpl struct {
	RemoteFileService *services.RemoteFileService
}

// OnRemotePeerUpdated ...
func (e *remoteFilesStateEventImpl) OnRemoteFilesStateUpdated(peerId peer.PeerId) {
	// TODO: to be implement
	e.RemoteFileService.Sync(peerId)
}
