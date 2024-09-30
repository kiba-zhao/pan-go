package controllers_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"pan/app/net"
	"pan/extfs/constant"
	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/services"
	"time"

	appNode "pan/app/node"
	MockedAppNode "pan/mocks/pan/app/node"
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

}
