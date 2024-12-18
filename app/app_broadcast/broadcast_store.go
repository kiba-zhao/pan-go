package appbroadcast

import (
	"encoding/base64"
	"pan/app/broadcast"
	"pan/app/peer"
)

type BroadcastStore struct {
	Repo       AppBroadcastInfoRepository
	PeerModule peer.PeerModule
}

func (store *BroadcastStore) SelectOrCreate(info broadcast.BroadcastInfo) (broadcast.BroadcastInfo, error) {
	var modelInfo AppBroadcastInfo
	modelInfo.PeerID = base64.StdEncoding.EncodeToString(info.PeerID)
	modelInfo.Hightest = info.Heightest

	modelInfo, _, err := store.Repo.SelectOrCreate(modelInfo)
	if err != nil {
		return info, err
	}

	info.Heightest = modelInfo.Hightest
	info.UpdatedAt = modelInfo.UpdatedAt
	return info, nil
}

func (store *BroadcastStore) SaveHighest(peerId peer.PeerID, heightest uint64) error {
	var modelInfo AppBroadcastInfo
	modelInfo.PeerID = base64.StdEncoding.EncodeToString(peerId)
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

func (store *BroadcastStore) Delete(peerId peer.PeerID) error {
	return store.Repo.DeleteByPeerID(base64.StdEncoding.EncodeToString(peerId))
}

func (store *BroadcastStore) Init(peerId peer.PeerID) error {
	var modelInfo AppBroadcastInfo
	modelInfo.PeerID = base64.StdEncoding.EncodeToString(peerId)
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
