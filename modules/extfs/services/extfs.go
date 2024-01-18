package services

import (
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
)

type ExtFSService struct {
	PeerRepo  repositories.PeerRepository
	ExtFSRepo repositories.ExtFSRepository
}

// GetLatestOne ...
func (s *ExtFSService) GetLatestOneToRemote() (info models.PeerInfo, err error) {

	// p, err := s.PeerRepo.FindOne(peerId.String())
	// if err != nil {
	// 	return
	// }
	// if p.Enabled != true {
	// 	err = errors.New("Forbidden")
	// 	return

	// }

	return
}
