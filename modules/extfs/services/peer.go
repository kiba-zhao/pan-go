package services

import (
	"bytes"
	"pan/modules/extfs/events"
	extfs "pan/modules/extfs/peer"
	"pan/modules/extfs/repositories"
	"pan/peer"
)

type PeerService struct {
	API       extfs.API
	PeerRepo  repositories.PeerRepository
	PeerEvent events.RemotePeerEvent
}

// SyncRemotePeer ...
func (s *PeerService) SyncRemotePeer(peerId peer.PeerId) (err error) {

	peer, err := s.PeerRepo.FindOne(peerId.String())
	if err != nil || peer.Enabled != true {
		return
	}
	info, err := s.API.GetPeerInfo(peerId)
	if err != nil {
		return
	}

	if bytes.Equal(peer.RemoteHash, info.Hash) {
		return
	}

	peer.ID = peerId.String()
	peer.RemoteHash = info.Hash
	peer.RemoteTime = info.Time

	err = s.PeerRepo.Save(peer)
	if err != nil {
		return
	}

	s.PeerEvent.OnRemotePeerUpdated(peerId)

	return
}
