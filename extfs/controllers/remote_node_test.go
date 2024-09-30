package controllers_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	appModels "pan/app/models"
	"pan/app/net"
	appNode "pan/app/node"
	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/services"
	"testing"

	mockedAppNode "pan/mocks/pan/app/node"
	mockedAppServices "pan/mocks/pan/app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRemoteNodeController(t *testing.T) {
	setup := func() (net.WebApp, *controllers.RemoteNodeController) {
		app := net.NewWebApp()
		controller := &controllers.RemoteNodeController{}
		controller.SetupToWeb(app)

		controller.RemoteNodeService = &services.RemoteNodeService{}
		return app, controller
	}

	t.Run("GET /remote-nodes", func(t *testing.T) {
		web, ctrl := setup()

		// mock NodeManager
		mgr := &mockedAppNode.MockNodeManager{}
		defer mgr.AssertExpectations(t)
		provider := &mockedAppServices.MockNodeManagerProvider{}
		defer provider.AssertExpectations(t)
		ctrl.RemoteNodeService.Provider = provider
		provider.On("NodeManager").Once().Return(mgr)

		firstNodeId := appNode.NodeID([]byte("1st-node-id"))
		secondNodeId := appNode.NodeID([]byte("2st-node-id"))
		thirdNodeId := appNode.NodeID([]byte("3st-node-id"))
		fourthNodeId := appNode.NodeID([]byte("4st-node-id"))
		mgr.On("TraverseNodeID", mock.Anything).Once().Run(func(args mock.Arguments) {
			traverse := args.Get(0).(func(appNode.NodeID) error)
			traverse(firstNodeId)
			traverse(secondNodeId)
			traverse(thirdNodeId)
			traverse(fourthNodeId)
		}).Return(nil)

		// mock NodeExternalService
		nodeExternalService := &mockedAppServices.MockNodeExternalService{}
		defer nodeExternalService.AssertExpectations(t)
		ctrl.RemoteNodeService.NodeExternalService = nodeExternalService

		firstNode := appModels.Node{}
		firstNode.ID = 1
		firstNode.NodeID = base64.StdEncoding.EncodeToString(firstNodeId)
		firstNode.Name = "first-node"

		secondNode := appModels.Node{}
		secondNode.ID = 2
		secondNode.NodeID = base64.StdEncoding.EncodeToString(secondNodeId)
		secondNode.Name = "second-node"

		nodeExternalService.On("TraverseWithNodeIDs", mock.Anything, mock.Anything).Once().Run(func(args mock.Arguments) {
			traverse := args.Get(0).(func(appModels.Node) error)
			traverse(secondNode)
			traverse(firstNode)
		}).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/remote-nodes", nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []models.RemoteNode
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 4, len(results))

		// assert 1st remote node
		assert.Equal(t, secondNode.Name, results[0].Name)
		assert.Equal(t, secondNode.NodeID, results[0].NodeID)
		assert.Equal(t, secondNode.CreatedAt, results[0].CreatedAt)
		assert.Equal(t, secondNode.UpdatedAt, results[0].UpdatedAt)
		assert.True(t, results[0].Available)

		// assert 2nd remote node
		assert.Equal(t, firstNode.Name, results[1].Name)
		assert.Equal(t, firstNode.NodeID, results[1].NodeID)
		assert.Equal(t, firstNode.CreatedAt, results[1].CreatedAt)
		assert.Equal(t, firstNode.UpdatedAt, results[1].UpdatedAt)
		assert.True(t, results[1].Available)

		// assert 3rd remote node
		thirdNodeIdBase64 := base64.StdEncoding.EncodeToString(thirdNodeId)
		assert.Equal(t, "", results[2].Name)
		assert.Equal(t, thirdNodeIdBase64, results[2].NodeID)
		assert.True(t, results[2].Available)

		// assert 4th remote node
		fourthNodeIdBase64 := base64.StdEncoding.EncodeToString(fourthNodeId)
		assert.Equal(t, "", results[3].Name)
		assert.Equal(t, fourthNodeIdBase64, results[3].NodeID)
		assert.True(t, results[3].Available)
	})
}
