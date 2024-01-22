package controllers_test

import (
	"errors"
	"io"
	"net/http"
	"pan/modules/extfs/models"
	"pan/modules/extfs/peer/controllers"
	"pan/modules/extfs/services"
	"pan/peer"
	"testing"
	"time"

	mockedRepo "pan/mocks/pan/modules/extfs/repositories"
	mockedPeer "pan/mocks/pan/peer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func TestRemoteFilesStateInfo(t *testing.T) {

	// setup function return a pointer to controllers.RemoteFilesStateController
	setup := func() *controllers.RemoteFilesStateController {
		ctrl := new(controllers.RemoteFilesStateController)
		ctrl.FilesStateService = new(services.FilesStateService)
		ctrl.FilesStateService.PeerService = new(services.PeerService)
		return ctrl
	}

	t.Run("Get with respond", func(t *testing.T) {
		ctrl := setup()

		repo := new(mockedRepo.MockFilesStateRepository)
		defer repo.AssertExpectations(t)
		ctrl.FilesStateService.FilesStateRepo = repo
		var stateRow models.FilesState
		stateRow.Hash = []byte("hash")
		stateRow.CreatedAt = time.Now()
		repo.On("GetLastOne").Once().Return(stateRow, nil)

		generator := new(mockedPeer.MockPeerIdGenerator)
		defer generator.AssertExpectations(t)
		ctrl.FilesStateService.PeerService.PeerIdGenerator = generator
		peerId := peer.PeerId(uuid.New())
		generator.On("LocalPeerId").Once().Return(peerId)

		called := false
		defer assert.False(t, called)
		next := func() error {
			called = true
			return nil
		}
		ctx := new(mockedPeer.MockContext)
		defer ctx.AssertExpectations(t)
		resErr := errors.New("res error")
		var resReader io.Reader
		ctx.On("Respond", mock.Anything).Once().Return(resErr).Run(func(args mock.Arguments) {
			resReader = args.Get(0).(io.Reader)
		})

		err := ctrl.Get(ctx, next)
		assert.Equal(t, err, resErr)

		bodyBytes, err := io.ReadAll(resReader)
		assert.Nil(t, err)
		var info models.RemoteStateInfo
		err = proto.Unmarshal(bodyBytes, &info)
		assert.Nil(t, err)
		assert.Equal(t, stateRow.Hash, info.Hash)
		assert.Equal(t, stateRow.CreatedAt.Unix(), info.Time)
	})

	t.Run("Get with error", func(t *testing.T) {
		ctrl := setup()
		repo := new(mockedRepo.MockFilesStateRepository)
		defer repo.AssertExpectations(t)
		ctrl.FilesStateService.FilesStateRepo = repo
		repoErr := errors.New("repo error")
		repo.On("GetLastOne").Once().Return(models.FilesState{}, repoErr)

		called := false
		defer assert.False(t, called)
		next := func() error {
			called = true
			return nil
		}
		ctx := new(mockedPeer.MockContext)
		defer ctx.AssertExpectations(t)
		throwErr := errors.New("throw error")
		ctx.On("ThrowError", http.StatusInternalServerError, repoErr.Error()).Once().Return(throwErr)

		err := ctrl.Get(ctx, next)
		assert.Equal(t, err, throwErr)
	})
}
