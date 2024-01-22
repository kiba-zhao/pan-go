package services

import (
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
)

type FilesStateService struct {
	FilesStateRepo repositories.FilesStateRepository
	PeerService    *PeerService
}

// GetLastOneToRemote ...
func (s *FilesStateService) GetLastOneToRemote() (info models.RemoteStateInfo, err error) {

	stateRow, err := s.FilesStateRepo.GetLastOne()
	if err != nil {
		return
	}

	peerId := s.PeerService.GetPeerId()
	info.PeerId = peerId[:]
	info.Hash = stateRow.Hash
	info.Time = stateRow.CreatedAt.Unix()
	return
}
