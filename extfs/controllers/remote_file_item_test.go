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
	"strconv"
	"testing"
	"time"

	appNode "pan/app/node"
	MockedAppNode "pan/mocks/pan/app/node"
	MockedServices "pan/mocks/pan/extfs/services"

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

	setupForNode := func() (*appNode.App, *controllers.RemoteFileItemController) {
		ctrl := &controllers.RemoteFileItemController{}
		app := appNode.NewApp()
		ctrl.SetupToNode(app)

		ctrl.RemoteFileItemService = &services.RemoteFileItemService{}
		return app, ctrl
	}

	t.Run("GET /remote-file-items?nodeId=&itemId=", func(t *testing.T) {

		app, ctrl := setup()

		// mock NodeScopeModule
		nodeScopeModule := &MockedAppNode.MockNodeScopeModule{}
		defer nodeScopeModule.AssertExpectations(t)
		ctrl.RemoteFileItemService.NodeScopeModule = nodeScopeModule
		scope := []byte("scope")
		nodeScopeModule.On("NodeScope").Once().Return(scope)

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
		requestName := appNode.GenerateRouteName(scope, services.RequestAllRemoteFileItems)
		nodeModule.On("Do", nodeId, mock.Anything).Once().Return(&ctx.Response, nil).Run(func(args mock.Arguments) {
			request := args.Get(1).(*appNode.Request)
			assert.Equal(t, requestName, request.Name())
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

	t.Run("SearchForNode", func(t *testing.T) {

		app, ctrl := setupForNode()

		// mock FileItemService
		fileItemService := MockedServices.MockFileItemInternalService{}
		defer fileItemService.AssertExpectations(t)
		ctrl.RemoteFileItemService.FileItemService = &fileItemService

		var fileItem models.FileItem
		fileItem.ID = "fileItemId"
		fileItem.ItemID = 1
		fileItem.Name = "test.txt"
		fileItem.Size = 123
		fileItem.FileType = constant.FileTypeFile
		fileItem.ParentPath = "parentPath"
		fileItem.FilePath = "filePath"
		fileItem.Available = true
		fileItem.CreatedAt = time.Now()
		fileItem.UpdatedAt = time.Now()

		var condition models.RemoteFileItemRecordSearchCondition
		condition.ItemID = 1
		condition.ParentPath = &fileItem.ParentPath

		fileItemService.On("TraverseWithCondition", mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			traverseFn := args.Get(0).(func(item models.FileItem) error)
			traverseFn(fileItem)
			condition_ := args.Get(1).(models.FileItemSearchCondition)
			assert.Equal(t, condition.ItemID, int32(condition_.ItemID))
			assert.Equal(t, *condition.ParentPath, *condition_.ParentPath)
		})

		// request and response
		reqBytes, err := proto.Marshal(&condition)
		assert.Nil(t, err)
		req := appNode.NewRequest(services.RequestAllRemoteFileItems, bytes.NewReader(reqBytes))
		reqReader := appNode.MarshalRequest(req)
		var ctx appNode.Context
		appNode.InitContext(&ctx)
		err = appNode.UnmarshalRequest(reqReader, ctx.Request())
		assert.Nil(t, err)

		err = app.Run(&ctx, nil)
		assert.Nil(t, err)

		body, err := io.ReadAll(ctx.Body())
		assert.Nil(t, err)

		var results models.RemoteFileItemRecordList
		err = proto.Unmarshal(body, &results)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(results.Items))
		assert.Equal(t, fileItem.ItemID, uint(results.Items[0].ItemID))
		assert.Equal(t, fileItem.Name, results.Items[0].Name)
		assert.Equal(t, fileItem.Size, results.Items[0].Size)
		assert.Equal(t, fileItem.FileType, results.Items[0].FileType)
		assert.Equal(t, fileItem.ParentPath, results.Items[0].ParentPath)
		assert.Equal(t, fileItem.FilePath, results.Items[0].FilePath)
		assert.Equal(t, fileItem.Available, results.Items[0].Available)
		assert.Equal(t, fileItem.CreatedAt.Unix(), results.Items[0].CreatedAt)
		assert.Equal(t, fileItem.UpdatedAt.Unix(), results.Items[0].UpdatedAt)
	})
}
