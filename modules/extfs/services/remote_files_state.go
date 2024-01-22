package services

import (
	"bytes"
	"pan/modules/extfs/events"
	extfs "pan/modules/extfs/peer"
	"pan/modules/extfs/repositories"
	"pan/peer"
)

type RemoteFilesStateService struct {
	API                   extfs.API
	RemoteFilesStateRepo  repositories.RemoteFilesStateRepository
	RemoteFilesStateEvent events.RemoteFilesStateEvent
}

// Sync ...
func (s *RemoteFilesStateService) Sync(peerId peer.PeerId) (err error) {

	stateRow, err := s.RemoteFilesStateRepo.FindOne(peerId.String())
	if err != nil {
		return
	}
	info, err := s.API.GetRemoteFilesState(peerId)
	if err != nil {
		return
	}

	if bytes.Equal(stateRow.RemoteHash, info.Hash) {
		return
	}

	stateRow.ID = peerId.String()
	stateRow.RemoteHash = info.Hash
	stateRow.RemoteTime = info.Time

	err = s.RemoteFilesStateRepo.Save(stateRow)
	if err != nil {
		return
	}

	s.RemoteFilesStateEvent.OnRemoteFilesStateUpdated(peerId)

	return
}
