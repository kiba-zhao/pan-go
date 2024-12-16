package broadcast

import (
	"pan/app/peer"
	"time"
)

type BroadcastInfo struct {
	PeerID    peer.PeerID
	Heightest uint64
	UpdatedAt time.Time
}

type BroadcastStore interface {
	SelectOrCreate(info BroadcastInfo) (BroadcastInfo, error)
	SaveHighest(peerId peer.PeerID, heightest uint64) error
	Delete(peerId peer.PeerID) error
	Init(peerId peer.PeerID) error
}
