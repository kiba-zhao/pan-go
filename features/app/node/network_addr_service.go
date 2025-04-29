package node

import (
	"iter"
)

type NetworkAddrService struct {
	NetworkAddrRepo NetworkAddrRepository
}

func (s *NetworkAddrService) SeqByPeerId(peerId string) (iter.Seq[NetworkAddr], error) {
	return s.NetworkAddrRepo.SeqByPeerId(peerId)
}
