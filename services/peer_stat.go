package services

import (
	"pan/models"
	"pan/peer"

	"github.com/google/uuid"
)

type PeerStatService struct {
	Peer peer.Peer
}

// FindOne ...
func (s *PeerStatService) FindOne(id string) (stat models.PeerStat, err error) {
	peerId, err := uuid.Parse(id)
	if err != nil {
		return
	}

	stat.ID = id
	stat.Stat = s.Peer.Stat(peerId)

	return
}
