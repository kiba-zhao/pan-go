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
	"strconv"
	"testing"
	"time"

	appNode "pan/app/node"
	MockedAppNode "pan/mocks/pan/app/node"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func TestRemoteFileItemController(t *testing.T) {

	setup := func() (net.WebApp, *controllers.RemoteFileItemController) {
		ctrl := &controllers.RemoteFileItemController{}
		app := net.NewWebApp()
		ctrl.SetupToWeb(app)

		ctrl.RemoteFileItemService = &services.RemoteFileItemService{}
		return app, ctrl
	}

	t.Run("GET /remote-file-items?nodeId=&itemId=", func(t *testing.T) {

		app, ctrl := setup()

		nodeModule := &MockedAppNode.MockNodeModule{}
		defer nodeModule.AssertExpectations(t)
		ctrl.RemoteFileItemService.NodeModule = nodeModule

		// mock response
		nodeId := []byte("nodeId")
		var record models.RemoteFileItemRecord
		record.ID = "recordId"
		record.ItemID = 1
		record.Name = "test.txt"
		record.Size = 123
		record.FileType = constant.FileTypeFile
		record.ParentPath = "parentPath"
		record.FilePath = "filePath"
		record.Available = true
		record.CreatedAt = time.Now().Unix()
		record.UpdatedAt = time.Now().Unix()

		var recordList models.RemoteFileItemRecordList
		recordList.Items = append(recordList.Items, &record)
		resBody, err := proto.Marshal(&recordList)
		assert.Nil(t, err)

		var ctx appNode.Context
		appNode.InitContext(&ctx)
		ctx.Respond(bytes.NewReader(resBody))
		nodeModule.On("Do", nodeId, mock.Anything).Once().Return(&ctx.Response, nil).Run(func(args mock.Arguments) {
			request := args.Get(1).(*appNode.Request)
			assert.Equal(t, services.RequestAllRemoteFileItems, request.Name())
		})

		base64NodeId := base64.StdEncoding.EncodeToString(nodeId)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/remote-file-items", nil)
		q := req.URL.Query()
		q.Add("nodeId", base64NodeId)
		q.Add("itemId", strconv.FormatUint(uint64(record.ItemID), 10))
		req.URL.RawQuery = q.Encode()
		app.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []models.RemoteFileItem
		err = json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(results))

		assert.Equal(t, uint(record.ItemID), results[0].ItemID)
		assert.Equal(t, record.Name, results[0].Name)
		assert.Equal(t, record.Size, results[0].Size)
		assert.Equal(t, record.FileType, results[0].FileType)
		assert.Equal(t, record.ParentPath, results[0].ParentPath)
		assert.Equal(t, record.FilePath, results[0].FilePath)
		assert.Equal(t, record.Available, results[0].Available)
		assert.Equal(t, record.CreatedAt, results[0].CreatedAt.Unix())
		assert.Equal(t, record.UpdatedAt, results[0].UpdatedAt.Unix())

	})
}
