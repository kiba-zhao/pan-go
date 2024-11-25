package remotenode_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"pan/app/peer"
	"pan/app/web"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	appnode "pan/app/app_node"
	remotenode "pan/extfs/remote_node"

	mockedAppNode "pan/mocks/pan/app/app_node"
	mockedPeer "pan/mocks/pan/app/peer"
)

func TestRemoteNodeController(t *testing.T) {

	setup := func() (web.WebApp, *remotenode.RemoteNodeController) {
		app := web.NewWebApp()
		ctrl := &remotenode.RemoteNodeController{}
		ctrl.SetupToWeb(app)

		ctrl.RemoteNodeService = &remotenode.RemoteNodeService{}
		return app, ctrl
	}

	t.Run("GET /remote-nodes", func(t *testing.T) {
		webApp, ctrl := setup()

		// mock PeerManager
		mgr := &mockedPeer.MockPeerManager{}
		defer mgr.AssertExpectations(t)
		ctrl.RemoteNodeService.PeerManager = mgr

		firstPeerID := peer.PeerID([]byte("1st-peer-id"))
		secondPeerID := peer.PeerID([]byte("2st-peer-id"))
		thirdPeerID := peer.PeerID([]byte("3st-peer-id"))
		fourthPeerID := peer.PeerID([]byte("4st-peer-id"))
		mgr.On("TraversePeerID", mock.Anything).Once().Run(func(args mock.Arguments) {
			traverse := args.Get(0).(func(peer.PeerID) error)
			traverse(firstPeerID)
			traverse(secondPeerID)
			traverse(thirdPeerID)
			traverse(fourthPeerID)
		}).Return(nil)

		// mock NodeExternalService
		appNodeExternalService := &mockedAppNode.MockAppNodeExternalService{}
		defer appNodeExternalService.AssertExpectations(t)
		ctrl.RemoteNodeService.AppNodeExternalService = appNodeExternalService

		firstNode := appnode.AppNode{}
		firstNode.ID = 1
		firstNode.PeerID = base64.StdEncoding.EncodeToString(firstPeerID)
		firstNode.Name = "first-node"

		secondNode := appnode.AppNode{}
		secondNode.ID = 2
		secondNode.PeerID = base64.StdEncoding.EncodeToString(secondPeerID)
		secondNode.Name = "second-node"

		appNodeExternalService.On("TraverseWithPeerIDs", mock.Anything, mock.Anything).Once().Run(func(args mock.Arguments) {
			traverse := args.Get(0).(func(appnode.AppNode) error)
			traverse(secondNode)
			traverse(firstNode)
		}).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/remote-nodes", nil)
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []remotenode.RemoteNode
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(results))

		// assert 1st remote node
		assert.Equal(t, secondNode.Name, results[0].Name)
		assert.Equal(t, secondNode.PeerID, results[0].PeerID)
		assert.Equal(t, secondNode.CreatedAt, results[0].CreatedAt)
		assert.Equal(t, secondNode.UpdatedAt, results[0].UpdatedAt)
		assert.True(t, results[0].Available)

		// assert 2nd remote node
		assert.Equal(t, firstNode.Name, results[1].Name)
		assert.Equal(t, firstNode.PeerID, results[1].PeerID)
		assert.Equal(t, firstNode.CreatedAt, results[1].CreatedAt)
		assert.Equal(t, firstNode.UpdatedAt, results[1].UpdatedAt)
		assert.True(t, results[1].Available)
	})
}
