package controllers_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http/httptest"
	"pan/app/net"
	"pan/extfs/constant"
	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/services"
	"time"

	appNode "pan/app/node"
	MockedAppNode "pan/mocks/pan/app/node"
	MockedServices "pan/mocks/pan/extfs/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func TestRemoteNodeItemController(t *testing.T) {

	setup := func() (net.WebApp, *controllers.RemoteNodeItemController) {
		ctrl := &controllers.RemoteNodeItemController{}
		app := net.NewWebApp()
		ctrl.SetupToWeb(app)

		ctrl.RemoteNodeItemService = &services.RemoteNodeItemService{}
		return app, ctrl
	}

	setupForNode := func() (*appNode.App, *controllers.RemoteNodeItemController) {
		ctrl := &controllers.RemoteNodeItemController{}
		app := appNode.NewApp()
		ctrl.SetupToNode(app)

		ctrl.RemoteNodeItemService = &services.RemoteNodeItemService{}
		return app, ctrl
	}

	t.Run("GET /remote-node-items?nodeId=", func(t *testing.T) {

		app, ctrl := setup()

		nodeModule := &MockedAppNode.MockNodeModule{}
		defer nodeModule.AssertExpectations(t)
		ctrl.RemoteNodeItemService.NodeModule = nodeModule

		// mock response
		nodeId := []byte("nodeId")
		var record models.RemoteNodeItemRecord
		record.ID = 1
		record.Name = "test.txt"
		record.Size = 123
		record.FileType = constant.FileTypeFile
		record.Available = true
		record.CreatedAt = time.Now().Unix()
		record.UpdatedAt = time.Now().Unix()

		var recordList models.RemoteNodeItemRecordList
		recordList.Items = append(recordList.Items, &record)

		resBody, err := proto.Marshal(&recordList)
		assert.Nil(t, err)

		var ctx appNode.Context
		appNode.InitContext(&ctx)
		ctx.Respond(bytes.NewReader(resBody))
		nodeModule.On("Do", nodeId, mock.Anything).Once().Return(&ctx.Response, nil).Run(func(args mock.Arguments) {
			request := args.Get(1).(*appNode.Request)
			assert.Equal(t, services.RequestAllRemoteItems, request.Name())
		})

		base64NodeId := base64.StdEncoding.EncodeToString(nodeId)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/remote-node-items", nil)
		q := req.URL.Query()
		q.Add("nodeId", base64NodeId)
		req.URL.RawQuery = q.Encode()
		app.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []models.RemoteNodeItem
		err = json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(results))

		assert.Equal(t, uint(record.ID), results[0].ItemID)
		assert.Equal(t, base64NodeId, results[0].NodeID)
		assert.Equal(t, record.Name, results[0].Name)
		assert.Equal(t, record.FileType, results[0].FileType)
		assert.Equal(t, record.Size, results[0].Size)
		assert.Equal(t, record.Available, results[0].Available)
		assert.Equal(t, record.CreatedAt, results[0].CreatedAt.Unix())
		assert.Equal(t, record.UpdatedAt, results[0].UpdatedAt.Unix())

	})

	t.Run("SearchForNode", func(t *testing.T) {

		app, ctrl := setupForNode()

		// mock NodeItemService
		nodeItemService := MockedServices.MockNodeItemInternalService{}
		defer nodeItemService.AssertExpectations(t)
		ctrl.RemoteNodeItemService.NodeItemService = &nodeItemService

		var nodeItem models.NodeItem
		nodeItem.ID = 1
		nodeItem.Name = "test.txt"
		nodeItem.FileType = constant.FileTypeFile
		nodeItem.Size = 123
		nodeItem.Available = true
		nodeItem.CreatedAt = time.Now()
		nodeItem.UpdatedAt = time.Now()

		nodeItemService.On("TraverseAll", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			traverseFn := args.Get(0).(func(models.NodeItem) error)
			traverseFn(nodeItem)
		})

		// request and response
		req := appNode.NewRequest(services.RequestAllRemoteItems, nil)
		reqReader := appNode.MarshalRequest(req)
		var ctx appNode.Context
		appNode.InitContext(&ctx)
		err := appNode.UnmarshalRequest(reqReader, ctx.Request())
		assert.Nil(t, err)

		err = app.Run(&ctx, nil)
		assert.Nil(t, err)

		body, err := io.ReadAll(ctx.Body())
		assert.Nil(t, err)

		var results models.RemoteNodeItemRecordList
		err = proto.Unmarshal(body, &results)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(results.Items))
		assert.Equal(t, int32(nodeItem.ID), results.Items[0].ID)
		assert.Equal(t, nodeItem.Name, results.Items[0].Name)
		assert.Equal(t, nodeItem.FileType, results.Items[0].FileType)
		assert.Equal(t, nodeItem.Size, results.Items[0].Size)
		assert.Equal(t, nodeItem.Available, results.Items[0].Available)
		assert.Equal(t, nodeItem.CreatedAt.Unix(), results.Items[0].CreatedAt)
		assert.Equal(t, nodeItem.UpdatedAt.Unix(), results.Items[0].UpdatedAt)
	})
}
