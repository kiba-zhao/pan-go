package controllers_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"pan/app/net"
	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/services"
	"strconv"
	"strings"
	"testing"

	appModels "pan/app/models"
	appNode "pan/app/node"
	mockedAppNode "pan/mocks/pan/app/node"
	mockedAppServices "pan/mocks/pan/app/services"
	mockedServices "pan/mocks/pan/extfs/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestItemController(t *testing.T) {

	setup := func() (web net.WebApp, ctrl *controllers.ItemController) {
		ctrl = new(controllers.ItemController)
		web = net.NewWebApp()
		ctrl.SetupToWeb(web)

		ctrl.ItemService = &services.ItemService{}
		return web, ctrl
	}

	t.Run("GET /items?fileType=N&fileType=RN", func(t *testing.T) {

		web, ctrl := setup()

		// mock SettingsExternalService
		settingsExternalService := &mockedAppServices.MockSettingsExternalService{}
		defer settingsExternalService.AssertExpectations(t)
		ctrl.ItemService.SettingsExternalService = settingsExternalService

		settings := appModels.Settings{}
		settings.NodeID = "test-node-id"
		settings.Name = "test-name"
		settingsExternalService.On("Load").Once().Return(settings)

		// mock NodeManager
		mgr := &mockedAppNode.MockNodeManager{}
		defer mgr.AssertExpectations(t)
		provider := &mockedAppServices.MockNodeManagerProvider{}
		defer provider.AssertExpectations(t)
		ctrl.ItemService.Provider = provider
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
		ctrl.ItemService.NodeExternalService = nodeExternalService

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

		// request and response
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/items", nil)
		q := req.URL.Query()
		q.Add("itemType", services.ItemTypeNode)
		q.Add("itemType", services.ItemTypeRemoteNode)
		req.URL.RawQuery = q.Encode()
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []models.Item
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 5, len(results))
		// assert local node
		assert.Equal(t, strings.Join([]string{services.ItemTypeNode, settings.NodeID}, services.ItemIDSep), results[0].ID)
		assert.Equal(t, settings.Name, results[0].Name)

		// assert 1st remote node
		assert.Equal(t, strings.Join([]string{services.ItemTypeRemoteNode, secondNode.NodeID}, services.ItemIDSep), results[1].ID)
		assert.Equal(t, secondNode.Name, results[1].Name)
		assert.Equal(t, secondNode.NodeID, *results[1].LinkID)

		// assert 2nd remote node
		assert.Equal(t, strings.Join([]string{services.ItemTypeRemoteNode, firstNode.NodeID}, services.ItemIDSep), results[2].ID)
		assert.Equal(t, firstNode.Name, results[2].Name)
		assert.Equal(t, firstNode.NodeID, *results[2].LinkID)

		// assert 3rd remote node
		thirdNodeIdBase64 := base64.StdEncoding.EncodeToString(thirdNodeId)
		assert.Equal(t, strings.Join([]string{services.ItemTypeRemoteNode, thirdNodeIdBase64}, services.ItemIDSep), results[3].ID)
		assert.Equal(t, "", results[3].Name)
		assert.Equal(t, thirdNodeIdBase64, *results[3].LinkID)

		// assert 4th remote node
		fourthNodeIdBase64 := base64.StdEncoding.EncodeToString(fourthNodeId)
		assert.Equal(t, strings.Join([]string{services.ItemTypeRemoteNode, fourthNodeIdBase64}, services.ItemIDSep), results[4].ID)
		assert.Equal(t, "", results[4].Name)
		assert.Equal(t, fourthNodeIdBase64, *results[4].LinkID)

	})

	t.Run("GET /items?parentId=? for node item", func(t *testing.T) {
		web, ctrl := setup()

		nodeId := "test node id"
		parentId := strings.Join([]string{services.ItemTypeNode, nodeId}, services.ItemIDSep)

		// mock NodeItemInternalService
		nodeItemInternalService := &mockedServices.MockNodeItemInternalService{}
		defer nodeItemInternalService.AssertExpectations(t)
		ctrl.ItemService.NodeItemInternalService = nodeItemInternalService

		var nodeItem models.NodeItem
		nodeItem.ID = 123
		nodeItem.Name = "test node item"
		nodeItem.Available = true
		nodeItem.FileType = services.FileTypeFolder

		nodeItemInternalService.On("TraverseAll", mock.Anything).Once().Run(func(args mock.Arguments) {
			traverse := args.Get(0).(func(models.NodeItem) error)
			traverse(nodeItem)
		}).Return(nil)

		// request and response
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/items", nil)
		q := req.URL.Query()
		q.Add("parentId", parentId)
		req.URL.RawQuery = q.Encode()
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []models.Item
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(results))
		linkId := strconv.FormatUint(uint64(nodeItem.ID), 10)
		id := strings.Join([]string{services.ItemTypeNodeItem, linkId}, services.ItemIDSep)
		assert.Equal(t, id, results[0].ID)
		assert.Equal(t, nodeItem.Name, results[0].Name)
		assert.Equal(t, services.ItemTypeFolder, results[0].ItemType)
		assert.Equal(t, nodeItem.Available, results[0].Available)
		assert.Equal(t, linkId, *results[0].LinkID)
		assert.Equal(t, parentId, *results[0].ParentID)
	})

	t.Run("GET /items?parentId=? for remote node item", func(t *testing.T) {
		web, ctrl := setup()

		nodeId := []byte("test node id")
		base64NodeID := base64.StdEncoding.EncodeToString(nodeId)
		parentId := strings.Join([]string{services.ItemTypeRemoteNode, base64NodeID}, services.ItemIDSep)

		// mock RemoteNodeItemInternalService
		remoteNodeItemInternalService := &mockedServices.MockRemoteNodeItemInternalService{}
		defer remoteNodeItemInternalService.AssertExpectations(t)
		ctrl.ItemService.RemoteNodeItemInternalService = remoteNodeItemInternalService

		var remoteNodeItem models.RemoteNodeItem
		remoteNodeItem.ID = 123
		remoteNodeItem.Name = "test remote node item"
		remoteNodeItem.Available = true
		remoteNodeItem.FileType = services.FileTypeFolder

		remoteNodeItemInternalService.On("TraverseAllWithNodeID", mock.Anything, nodeId).Once().Run(func(args mock.Arguments) {
			traverse := args.Get(0).(func(*models.RemoteNodeItem) error)
			traverse(&remoteNodeItem)
		}).Return(nil)

		// request and response
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/items", nil)
		q := req.URL.Query()
		q.Add("parentId", parentId)
		req.URL.RawQuery = q.Encode()
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []models.Item
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(results))
		linkId := strconv.FormatUint(uint64(remoteNodeItem.ID), 10)
		id := strings.Join([]string{services.ItemTypeRemoteNodeItem, linkId}, services.ItemIDSep)
		assert.Equal(t, id, results[0].ID)
		assert.Equal(t, remoteNodeItem.Name, results[0].Name)
		assert.Equal(t, services.ItemTypeFolder, results[0].ItemType)
		assert.Equal(t, remoteNodeItem.Available, results[0].Available)
		assert.Equal(t, linkId, *results[0].LinkID)
		assert.Equal(t, parentId, *results[0].ParentID)
	})
}
