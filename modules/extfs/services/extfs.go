package services

import (
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
	"pan/peer"
)

type ExtFSService struct {
	ExtFSRepo       repositories.ExtFSRepository
	PeerIdGenerator peer.PeerIdGenerator
}

// GetPeerId ...
func (s *ExtFSService) GetPeerId() peer.PeerId {
	return s.PeerIdGenerator.LocalPeerId()
}

// GetLatestOne ...
func (s *ExtFSService) GetLatestOneToRemote() (info models.PeerInfo, err error) {

	extfsInfo, err := s.ExtFSRepo.GetLatestOne()
	if err != nil {
		return
	}

	peerId := s.GetPeerId()
	info.PeerId = peerId[:]
	info.Hash = extfsInfo.Hash
	info.Time = extfsInfo.CreatedAt.Unix()
	return
}
