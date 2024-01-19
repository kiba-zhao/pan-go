package controllers_test

import (
	"errors"
	"net/http"
	"pan/modules/extfs/models"
	"pan/modules/extfs/peer/controllers"
	"pan/modules/extfs/services"
	"pan/peer"
	"testing"

	mockedRepo "pan/mocks/pan/modules/extfs/repositories"
	mockedPeer "pan/mocks/pan/peer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestAuth ...
func TestAuth(t *testing.T) {

	// setup function return a pointer to controllers.AuthController
	setup := func() *controllers.AuthController {
		ctrl := new(controllers.AuthController)
		ctrl.PeerService = new(services.PeerService)
		return ctrl
	}

	t.Run("Auth with pass next", func(t *testing.T) {
		ctrl := setup()
		repo := new(mockedRepo.MockPeerRepository)
		ctrl.PeerService.PeerRepo = repo

		peerId := peer.PeerId(uuid.New())
		repo.On("FindOne", peerId.String()).Once().Return(models.Peer{ID: peerId.String(), Enabled: true}, nil)

		called := false
		nextErr := errors.New("next error")
		next := func() error {
			called = true
			return nextErr
		}
		ctx := new(mockedPeer.MockContext)
		ctx.On("PeerId").Once().Return(peerId)
		err := ctrl.Auth(ctx, next)

		assert.Equal(t, err, nextErr)
		assert.True(t, called)

		repo.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("Auth with forbidden", func(t *testing.T) {
		ctrl := setup()
		repo := new(mockedRepo.MockPeerRepository)
		ctrl.PeerService.PeerRepo = repo

		peerId := peer.PeerId(uuid.New())
		repoErr := errors.New("repository error")
		repo.On("FindOne", peerId.String()).Once().Return(models.Peer{ID: peerId.String()}, repoErr)

		called := false
		nextErr := errors.New("next error")
		next := func() error {
			called = true
			return nextErr
		}
		ctx := new(mockedPeer.MockContext)
		ctx.On("PeerId").Once().Return(peerId)
		throwErr := errors.New("throw error")
		ctx.On("ThrowError", http.StatusForbidden, "Forbidden").Once().Return(throwErr)

		err := ctrl.Auth(ctx, next)

		assert.Equal(t, throwErr, err)
		assert.False(t, called)

		repo.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})

	t.Run("Auth with disabled", func(t *testing.T) {
		ctrl := setup()
		repo := new(mockedRepo.MockPeerRepository)
		ctrl.PeerService.PeerRepo = repo

		peerId := peer.PeerId(uuid.New())
		repo.On("FindOne", peerId.String()).Once().Return(models.Peer{ID: peerId.String(), Enabled: false}, nil)

		called := false
		nextErr := errors.New("next error")
		next := func() error {
			called = true
			return nextErr
		}
		ctx := new(mockedPeer.MockContext)
		ctx.On("PeerId").Once().Return(peerId)
		throwErr := errors.New("throw error")
		ctx.On("ThrowError", http.StatusForbidden, "Forbidden").Once().Return(throwErr)

		err := ctrl.Auth(ctx, next)

		assert.Equal(t, throwErr, err)
		assert.False(t, called)

		repo.AssertExpectations(t)
		ctx.AssertExpectations(t)
	})
}
