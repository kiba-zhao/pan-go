package events

import (
	"pan/peer"
)

type RemoteFilesStateEvent interface {
	OnRemoteFilesStateUpdated(peerId peer.PeerId)
}
