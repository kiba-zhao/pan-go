package services

import (
	"pan/modules/extfs/repositories"
	"pan/peer"
)

type RemotePeerService struct {
	RemotePeerRepo repositories.RemotePeerRepository
}

// HasRemotePeer ...
func (s *RemotePeerService) HasEnabled(peerId peer.PeerId) bool {
	remote, err := s.RemotePeerRepo.FindOne(peerId.String())
	if err != nil {
		return false
	}
	return remote.Enabled
}
