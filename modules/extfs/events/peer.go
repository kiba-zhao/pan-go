package events

import (
	"pan/peer"
)

type RemotePeerEvent interface {
	OnRemotePeerUpdated(peerId peer.PeerId)
}
