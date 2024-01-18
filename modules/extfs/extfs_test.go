package extfs_test

import (
	"errors"
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
		efs.PeerService = new(services.PeerService)
		return efs
	}

	t.Run("OnNodeAdded success", func(t *testing.T) {

		peerId := peer.PeerId(uuid.New())
		efs := setup()

		repo := new(mockedRepo.MockPeerRepository)
		efs.PeerService.PeerRepo = repo
		var peerRow models.Peer
		peerRow.RemoteHash = []byte{1, 2, 3, 4, 5, 6, 7}
		peerRow.Enabled = true
		repo.On("FindOne", peerId.String()).Once().Return(peerRow, nil)

		api := new(mockedPeer.MockAPI)
		efs.PeerService.API = api
		var peerInfo models.PeerInfo
		peerInfo.Hash = []byte{11, 12, 13, 14, 15, 16, 17}
		peerInfo.Time = 123
		api.On("GetPeerInfo", peerId).Once().Return(peerInfo, nil)

		peerRow.ID = peerId.String()
		peerRow.RemoteHash = peerInfo.Hash
		peerRow.RemoteTime = peerInfo.Time
		repo.On("Save", peerRow).Once().Return(nil)

		event := new(mockedEvent.MockRemotePeerEvent)
		efs.PeerService.PeerEvent = event
		event.On("OnRemotePeerUpdated", peerId).Once()

		efs.OnNodeAdded(peerId)

		repo.AssertExpectations(t)
		api.AssertExpectations(t)
		event.AssertExpectations(t)
	})

	t.Run("OnNodeAdded failed with Repo.FindOne", func(t *testing.T) {

		peerId := peer.PeerId(uuid.New())
		efs := setup()

		repo := new(mockedRepo.MockPeerRepository)
		efs.PeerService.PeerRepo = repo
		var peerRow models.Peer
		peerRow.Enabled = false
		repo.On("FindOne", peerId.String()).Once().Return(peerRow, errors.New("Test Error"))
		repo.On("FindOne", peerId.String()).Once().Return(peerRow, nil)

		efs.OnNodeAdded(peerId)
		efs.OnNodeAdded(peerId)

		repo.AssertExpectations(t)
	})

	t.Run("OnNodeAdded failed with API.GetPeerInfo", func(t *testing.T) {

		peerId := peer.PeerId(uuid.New())
		efs := setup()

		repo := new(mockedRepo.MockPeerRepository)
		efs.PeerService.PeerRepo = repo
		var peerRow models.Peer
		peerRow.RemoteHash = []byte{1, 2, 3, 4, 5, 6, 7}
		peerRow.Enabled = true
		repo.On("FindOne", peerId.String()).Once().Return(peerRow, nil)

		api := new(mockedPeer.MockAPI)
		efs.PeerService.API = api
		var peerInfo models.PeerInfo
		api.On("GetPeerInfo", peerId).Once().Return(peerInfo, errors.New("Test Error"))

		efs.OnNodeAdded(peerId)

		repo.AssertExpectations(t)
		api.AssertExpectations(t)
	})

	t.Run("OnNodeAdded ingore with same RemoteHash", func(t *testing.T) {

		peerId := peer.PeerId(uuid.New())
		efs := setup()

		repo := new(mockedRepo.MockPeerRepository)
		efs.PeerService.PeerRepo = repo
		var peerRow models.Peer
		peerRow.RemoteHash = []byte{1, 2, 3, 4, 5, 6, 7}
		peerRow.Enabled = true
		repo.On("FindOne", peerId.String()).Once().Return(peerRow, nil)

		api := new(mockedPeer.MockAPI)
		efs.PeerService.API = api
		var peerInfo models.PeerInfo
		peerInfo.Hash = peerRow.RemoteHash
		api.On("GetPeerInfo", peerId).Once().Return(peerInfo, nil)

		efs.OnNodeAdded(peerId)

		repo.AssertExpectations(t)
		api.AssertExpectations(t)
	})

	t.Run("OnNodeAdded failed with Repo.Save", func(t *testing.T) {

		peerId := peer.PeerId(uuid.New())
		efs := setup()

		repo := new(mockedRepo.MockPeerRepository)
		efs.PeerService.PeerRepo = repo
		var peerRow models.Peer
		peerRow.RemoteHash = []byte{1, 2, 3, 4, 5, 6, 7}
		peerRow.Enabled = true
		repo.On("FindOne", peerId.String()).Once().Return(peerRow, nil)

		api := new(mockedPeer.MockAPI)
		efs.PeerService.API = api
		var peerInfo models.PeerInfo
		peerInfo.Hash = []byte{11, 12, 13, 14, 15, 16, 17}
		peerInfo.Time = 123
		api.On("GetPeerInfo", peerId).Once().Return(peerInfo, nil)

		peerRow.ID = peerId.String()
		peerRow.RemoteHash = peerInfo.Hash
		peerRow.RemoteTime = peerInfo.Time
		repo.On("Save", peerRow).Once().Return(errors.New("Test Error"))

		efs.OnNodeAdded(peerId)

		repo.AssertExpectations(t)
		api.AssertExpectations(t)
	})

}
