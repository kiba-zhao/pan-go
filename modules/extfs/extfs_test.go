package extfs_test

import (
	"pan/modules/extfs"
	"pan/modules/extfs/models"
	"pan/modules/extfs/services"
	"pan/peer"
	"testing"

	mockedEvent "pan/mocks/pan/modules/extfs/events"
	mockedPeer "pan/mocks/pan/modules/extfs/peer"
	mockedRepo "pan/mocks/pan/modules/extfs/repositories"

	"github.com/google/uuid"
)

// TestExtFS ...
func TestExtFS(t *testing.T) {

	// setup ...
	setup := func() *extfs.ExtFS {
		efs := new(extfs.ExtFS)
		efs.RemotePeerService = new(services.RemotePeerService)
		efs.RemoteFilesStateService = new(services.RemoteFilesStateService)
		return efs
	}

	t.Run("OnNodeAdded with enabled", func(t *testing.T) {

		peerId := peer.PeerId(uuid.New())
		efs := setup()

		remotePeerRepo := new(mockedRepo.MockRemotePeerRepository)
		defer remotePeerRepo.AssertExpectations(t)
		efs.RemotePeerService.RemotePeerRepo = remotePeerRepo
		var remotePeerRow models.RemotePeer
		remotePeerRow.Enabled = true
		remotePeerRepo.On("FindOne", peerId.String()).Once().Return(remotePeerRow, nil)

		filesStateRepo := new(mockedRepo.MockRemoteFilesStateRepository)
		defer filesStateRepo.AssertExpectations(t)
		efs.RemoteFilesStateService.RemoteFilesStateRepo = filesStateRepo
		var stateRow models.RemoteFilesState
		stateRow.RemoteHash = []byte{1, 2, 3, 4, 5, 6, 7}
		filesStateRepo.On("FindOne", peerId.String()).Once().Return(stateRow, nil)

		api := new(mockedPeer.MockAPI)
		efs.RemoteFilesStateService.API = api
		defer api.AssertExpectations(t)
		var stateInfo models.RemoteStateInfo
		stateInfo.Hash = []byte{11, 12, 13, 14, 15, 16, 17}
		stateInfo.Time = 123
		api.On("GetRemoteFilesState", peerId).Once().Return(stateInfo, nil)

		stateRow.ID = peerId.String()
		stateRow.RemoteHash = stateInfo.Hash
		stateRow.RemoteTime = stateInfo.Time
		filesStateRepo.On("Save", stateRow).Once().Return(nil)

		event := new(mockedEvent.MockRemoteFilesStateEvent)
		defer event.AssertExpectations(t)
		efs.RemoteFilesStateService.RemoteFilesStateEvent = event
		event.On("OnRemoteFilesStateUpdated", peerId).Once()

		efs.OnNodeAdded(peerId)

	})

	t.Run("OnNodeAdded with disabled", func(t *testing.T) {
		peerId := peer.PeerId(uuid.New())
		efs := setup()

		remotePeerRepo := new(mockedRepo.MockRemotePeerRepository)
		defer remotePeerRepo.AssertExpectations(t)
		efs.RemotePeerService.RemotePeerRepo = remotePeerRepo
		var remotePeerRow models.RemotePeer
		remotePeerRow.Enabled = false
		remotePeerRepo.On("FindOne", peerId.String()).Once().Return(remotePeerRow, nil)

		efs.OnNodeAdded(peerId)
	})
}
