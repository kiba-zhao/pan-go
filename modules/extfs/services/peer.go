package services

import "pan/peer"

type PeerService struct {
	PeerIdGenerator peer.PeerIdGenerator
}

func (s *PeerService) GetPeerId() peer.PeerId {
	return s.PeerIdGenerator.LocalPeerId()
}
