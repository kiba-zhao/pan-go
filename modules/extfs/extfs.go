package extfs

import (
	"pan/modules/extfs/services"
	"pan/peer"
)

type ExtFS struct {
	RemotePeerService       *services.RemotePeerService
	RemoteFilesStateService *services.RemoteFilesStateService
}

// OnNodeAdded ...
func (extfs *ExtFS) OnNodeAdded(peerId peer.PeerId) {
	hasEnabled := extfs.RemotePeerService.HasEnabled(peerId)
	if !hasEnabled {
		return
	}
	extfs.RemoteFilesStateService.Sync(peerId)
	// TODO: log error
}

// OnNodeRemoved ...
func (extfs *ExtFS) OnNodeRemoved(peerId peer.PeerId) {
	// TODO: to be implement
}

// OnNodeRemoved ...
func (extfs *ExtFS) OnRouteAdded(peerId peer.PeerId) {
	// TODO: to be implement
}

// OnNodeRemoved ...
func (extfs *ExtFS) OnRouteRemoved(peerId peer.PeerId) {
	// TODO: to be implement
}
