package services

import (
	extfs "pan/modules/extfs/peer"
	"pan/modules/extfs/repositories"
	"pan/peer"
)

type RemoteFileService struct {
	API                     extfs.API
	RemoteFilesStateService *RemoteFilesStateService
	RemoteFileRepo          repositories.RemoteFileRepository
}

func (s *RemoteFileService) Sync(peerId peer.PeerId) error {
	_, err := s.RemoteFilesStateService.FindOne(peerId)
	if err != nil {
		return err
	}

	return nil
}
