package extfs

import (
	"pan/modules/extfs/services"
	"pan/peer"
)

type ExtFS struct {
	PeerService *services.PeerService
}

// OnNodeAdded ...
func (extfs *ExtFS) OnNodeAdded(peerId peer.PeerId) {
	extfs.PeerService.SyncRemotePeer(peerId)
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
