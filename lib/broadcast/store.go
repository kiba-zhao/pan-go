// Define broadcast store for broadcast
package broadcast

import (
	"pan/lib/peer"
	"time"
)

type BroadcastInfo struct {
	PeerID    peer.PeerID
	Heightest uint64
	UpdatedAt time.Time
}

// BroadcastStore
type BroadcastStore interface {
	// SelectOrCreate selects or creates a BroadcastInfo from the store.
	//
	// If the peerID exists in the store, the function returns the associated BroadcastInfo.
	// If the peerID does not exist, the function creates a new BroadcastInfo with the
	// given peerID and heightest, and returns it. If there is an error with the store,
	// the function returns an error.
	//
	SelectOrCreate(info BroadcastInfo) (BroadcastInfo, error)
	// SaveHighest saves the highest heightest for the given peerId.
	//
	// If the peerId exists in the store and the given heightest is higher than
	// the current heightest, the function updates the heightest for the peerId.
	// If the peerId does not exist, the function creates a new BroadcastInfo
	// with the given peerId and heightest, and saves it. If there is an error
	// with the store, the function returns an error.
	//
	SaveHighest(peerId peer.PeerID, heightest uint64) error
	// Delete deletes the BroadcastInfo associated with the given peerId.
	//
	// If the peerId does not exist in the store, the function returns an error.
	// If there is an error with the store, the function returns an error.
	//
	Delete(peerId peer.PeerID) error
	// Init initializes the BroadcastStore for the given peerId.
	//
	// If the peerId does not exist in the store, the function creates a new BroadcastInfo
	// with the given peerId and heightest 0, and saves it. If the peerId exists,
	// the function does nothing. If there is an error with the store, the function
	// returns an error.
	//
	Init(peerId peer.PeerID) error
}
