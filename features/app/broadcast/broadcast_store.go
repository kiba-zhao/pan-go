// Define app broadcast store
package broadcast

import (
	appnode "pan/features/app/node"
	"pan/lib/broadcast"
	"pan/lib/peer"
)

type BroadcastStore struct {
	Repo       AppBroadcastInfoRepository
	PeerModule peer.PeerModule
}

// SelectOrCreate selects or creates a BroadcastInfo from the store.
//
// If the peerID exists in the store, the function returns the associated BroadcastInfo.
// If the peerID does not exist, the function creates a new BroadcastInfo with the
// given peerID and heightest, and returns it. If there is an error with the store,
// the function returns an error.
func (store *BroadcastStore) SelectOrCreate(info broadcast.BroadcastInfo) (broadcast.BroadcastInfo, error) {
	var modelInfo AppBroadcastInfo
	modelInfo.PeerID = appnode.EncodePeerID(info.PeerID)
	modelInfo.Hightest = info.Heightest

	modelInfo, _, err := store.Repo.SelectOrCreate(modelInfo)
	if err != nil {
		return info, err
	}

	info.Heightest = modelInfo.Hightest
	info.UpdatedAt = modelInfo.UpdatedAt
	return info, nil
}

// SaveHighest saves the highest heightest for the given peerId.
//
// If the peerId exists in the store and the given heightest is higher than
// the current heightest, the function updates the heightest for the peerId.
// If the peerId does not exist, the function creates a new BroadcastInfo
// with the given peerId and heightest, and saves it. If there is an error
// with the store, the function returns an error.
func (store *BroadcastStore) SaveHighest(peerId peer.PeerID, heightest uint64) error {
	var modelInfo AppBroadcastInfo
	modelInfo.PeerID = appnode.EncodePeerID(peerId)
	modelInfo.Hightest = heightest

	modelInfo, ok, err := store.Repo.SelectOrCreate(modelInfo)
	if err != nil || modelInfo.Hightest >= heightest {
		return err
	}

	if ok && modelInfo.Hightest <= heightest {
		modelInfo.Hightest = heightest
		modelInfo, err = store.Repo.UpdateHeightest(modelInfo)
	}
	return err
}

// Delete removes the BroadcastInfo associated with the given peerId from the store.
//
// If the peerId does not exist in the store, the function returns an error.
// If there is an error with the store, the function returns an error.

func (store *BroadcastStore) Delete(peerId peer.PeerID) error {
	return store.Repo.DeleteByPeerID(appnode.EncodePeerID(peerId))
}

// Init initializes the BroadcastStore for the given peerId.
//
// If the peerId does not exist in the store, the function creates a new BroadcastInfo
// with the given peerId and heightest 0, and saves it. If the peerId exists,
// the function does nothing. If there is an error with the store, the function
// returns an error.
func (store *BroadcastStore) Init(peerId peer.PeerID) error {
	var modelInfo AppBroadcastInfo
	modelInfo.PeerID = appnode.EncodePeerID(peerId)
	modelInfo.Hightest = 0

	modelInfo, ok, err := store.Repo.SelectOrCreate(modelInfo)
	if err != nil {
		return err
	}

	if ok {
		modelInfo.Hightest++
		modelInfo, err = store.Repo.UpdateHeightest(modelInfo)
	}

	return err
}
