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
	setup := func() *controllers.GuardController {
		ctrl := new(controllers.GuardController)
		ctrl.RemotePeerService = new(services.RemotePeerService)
		return ctrl
	}

	t.Run("Auth with pass next", func(t *testing.T) {
		ctrl := setup()
		repo := new(mockedRepo.MockRemotePeerRepository)
		defer repo.AssertExpectations(t)
		ctrl.RemotePeerService.RemotePeerRepo = repo

		peerId := peer.PeerId(uuid.New())
		var remotePeerRow models.RemotePeer
		remotePeerRow.Enabled = true
		repo.On("FindOne", peerId.String()).Once().Return(remotePeerRow, nil)

		called := false
		nextErr := errors.New("next error")
		next := func() error {
			called = true
			return nextErr
		}
		ctx := new(mockedPeer.MockContext)
		defer ctx.AssertExpectations(t)
		ctx.On("PeerId").Once().Return(peerId)
		err := ctrl.Auth(ctx, next)

		assert.Equal(t, err, nextErr)
		assert.True(t, called)

	})

	t.Run("Auth with forbidden", func(t *testing.T) {
		ctrl := setup()
		repo := new(mockedRepo.MockRemotePeerRepository)
		defer repo.AssertExpectations(t)
		ctrl.RemotePeerService.RemotePeerRepo = repo

		peerId := peer.PeerId(uuid.New())
		var remotePeerRow models.RemotePeer
		remotePeerRow.Enabled = true
		repoErr := errors.New("repository error")
		repo.On("FindOne", peerId.String()).Once().Return(remotePeerRow, repoErr)

		called := false
		nextErr := errors.New("next error")
		next := func() error {
			called = true
			return nextErr
		}
		ctx := new(mockedPeer.MockContext)
		defer ctx.AssertExpectations(t)
		ctx.On("PeerId").Once().Return(peerId)
		throwErr := errors.New("throw error")
		ctx.On("ThrowError", http.StatusForbidden, "Forbidden").Once().Return(throwErr)

		err := ctrl.Auth(ctx, next)

		assert.Equal(t, throwErr, err)
		assert.False(t, called)
	})

	t.Run("Auth with disabled", func(t *testing.T) {
		ctrl := setup()
		repo := new(mockedRepo.MockRemotePeerRepository)
		defer repo.AssertExpectations(t)
		ctrl.RemotePeerService.RemotePeerRepo = repo

		peerId := peer.PeerId(uuid.New())
		var remotePeerRow models.RemotePeer
		remotePeerRow.Enabled = false
		repo.On("FindOne", peerId.String()).Once().Return(remotePeerRow, nil)

		called := false
		nextErr := errors.New("next error")
		next := func() error {
			called = true
			return nextErr
		}
		ctx := new(mockedPeer.MockContext)
		defer ctx.AssertExpectations(t)
		ctx.On("PeerId").Once().Return(peerId)
		throwErr := errors.New("throw error")
		ctx.On("ThrowError", http.StatusForbidden, "Forbidden").Once().Return(throwErr)

		err := ctrl.Auth(ctx, next)

		assert.Equal(t, throwErr, err)
		assert.False(t, called)
	})
}
