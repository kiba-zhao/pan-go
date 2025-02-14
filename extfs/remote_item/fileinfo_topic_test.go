package remoteitem_test

import (
	"bytes"
	"io"
	"pan/app/peer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	nodeitem "pan/extfs/node_item"
	remoteitem "pan/extfs/remote_item"

	mockedNodeItem "pan/mocks/pan/extfs/node_item"
)

func TestRemoteFileInfoTopic(t *testing.T) {

	setup := func() (*peer.App, *remoteitem.RemoteFileInfoTopic) {
		ctrl := &remoteitem.RemoteFileInfoTopic{}
		app := peer.NewApp()
		ctrl.SetupToPeer(app)

		ctrl.RemoteFileInfoService = &remoteitem.RemoteFileInfoService{}
		return app, ctrl
	}

	t.Run("Search", func(t *testing.T) {

		app, ctrl := setup()

		// mock NodeFileService
		nodeFileInfoService := mockedNodeItem.MockNodeFileInfoInternalService{}
		defer nodeFileInfoService.AssertExpectations(t)
		ctrl.RemoteFileInfoService.NodeFileInfoService = &nodeFileInfoService

		var nodeFile nodeitem.NodeFileInfo

		nodeFile.ItemID = 1
		nodeFile.Name = "test.txt"
		nodeFile.Size = 123
		nodeFile.FileType = nodeitem.FileTypeFile
		nodeFile.ParentPath = "parentPath"
		nodeFile.FilePath = "filePath"
		nodeFile.Available = true
		nodeFile.CreatedAt = time.Now()
		nodeFile.UpdatedAt = time.Now()

		var condition remoteitem.RemoteFileInfoRecordSearchCondition
		condition.ItemID = 1
		condition.ParentPath = nodeFile.ParentPath

		nodeFileInfoService.On("TraverseWithCondition", mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			traverseFn := args.Get(0).(func(item nodeitem.NodeFileInfo) error)
			traverseFn(nodeFile)
			condition_ := args.Get(1).(nodeitem.NodeFileInfoSearchCondition)
			assert.Equal(t, condition.ItemID, uint32(condition_.ItemID))
			assert.Equal(t, condition.ParentPath, condition_.ParentPath)
		})

		// request and response
		reqBytes, err := proto.Marshal(&condition)
		assert.Nil(t, err)
		req := peer.NewRequest(remoteitem.SearchRemoteFileInfos, bytes.NewReader(reqBytes))
		reqReader := peer.MarshalRequest(req)
		var ctx peer.Context
		peer.InitContext(&ctx)
		err = peer.UnmarshalRequest(reqReader, ctx.Request())
		assert.Nil(t, err)

		err = app.Run(&ctx, nil)
		assert.Nil(t, err)

		body, err := io.ReadAll(ctx)
		assert.Nil(t, err)

		var results remoteitem.RemoteFileInfoRecordList
		err = proto.Unmarshal(body, &results)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(results.Items))
		assert.Equal(t, nodeFile.ItemID, uint(results.Items[0].ItemID))
		assert.Equal(t, nodeFile.Name, results.Items[0].Name)
		assert.Equal(t, nodeFile.Size, results.Items[0].Size)
		assert.Equal(t, nodeFile.FileType, results.Items[0].FileType)
		assert.Equal(t, nodeFile.ParentPath, results.Items[0].ParentPath)
		assert.Equal(t, nodeFile.FilePath, results.Items[0].FilePath)
		assert.Equal(t, nodeFile.Available, results.Items[0].Available)
		assert.Equal(t, nodeFile.CreatedAt.Unix(), results.Items[0].CreatedAt)
		assert.Equal(t, nodeFile.UpdatedAt.Unix(), results.Items[0].UpdatedAt)
	})

	t.Run("Select", func(t *testing.T) {
		app, ctrl := setup()

		// mock NodeFileService
		nodeFileInfoService := mockedNodeItem.MockNodeFileInfoInternalService{}
		defer nodeFileInfoService.AssertExpectations(t)
		ctrl.RemoteFileInfoService.NodeFileInfoService = &nodeFileInfoService

		var nodeFile nodeitem.NodeFileInfo

		nodeFile.ItemID = 1
		nodeFile.Name = "test.txt"
		nodeFile.Size = 123
		nodeFile.FileType = nodeitem.FileTypeFile
		nodeFile.ParentPath = "parentPath"
		nodeFile.FilePath = "filePath"
		nodeFile.Available = true
		nodeFile.CreatedAt = time.Now()
		nodeFile.UpdatedAt = time.Now()

		nodeFileInfoService.On("Select", nodeFile.ItemID, nodeFile.FilePath).Once().Return(nodeFile, nil)

		var condition remoteitem.RemoteFileInfoRecordSelectCondition
		condition.ItemID = uint32(nodeFile.ItemID)
		condition.FilePath = nodeFile.FilePath

		// request and response
		reqBytes, err := proto.Marshal(&condition)
		assert.Nil(t, err)
		req := peer.NewRequest(remoteitem.SelectRemoteFileInfo, bytes.NewReader(reqBytes))
		reqReader := peer.MarshalRequest(req)
		var ctx peer.Context
		peer.InitContext(&ctx)
		err = peer.UnmarshalRequest(reqReader, ctx.Request())
		assert.Nil(t, err)

		err = app.Run(&ctx, nil)
		assert.Nil(t, err)

		body, err := io.ReadAll(ctx)
		assert.Nil(t, err)

		var results remoteitem.RemoteFileInfoRecord
		err = proto.Unmarshal(body, &results)
		assert.Nil(t, err)

		assert.Equal(t, nodeFile.ItemID, uint(results.ItemID))
		assert.Equal(t, nodeFile.Name, results.Name)
		assert.Equal(t, nodeFile.Size, results.Size)
		assert.Equal(t, nodeFile.FileType, results.FileType)
		assert.Equal(t, nodeFile.ParentPath, results.ParentPath)
		assert.Equal(t, nodeFile.FilePath, results.FilePath)
		assert.Equal(t, nodeFile.Available, results.Available)
		assert.Equal(t, nodeFile.CreatedAt.Unix(), results.CreatedAt)
		assert.Equal(t, nodeFile.UpdatedAt.Unix(), results.UpdatedAt)
	})
}
