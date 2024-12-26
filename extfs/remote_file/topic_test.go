package remotefile_test

import (
	"bytes"
	"io"
	"pan/app/peer"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	nodefile "pan/extfs/node_file"
	nodeitem "pan/extfs/node_item"
	remotefile "pan/extfs/remote_file"

	mockedNodeFile "pan/mocks/pan/extfs/node_file"
)

func TestRemoteFileTopic(t *testing.T) {

	setup := func() (*peer.App, *remotefile.RemoteFileTopic) {
		ctrl := &remotefile.RemoteFileTopic{}
		app := peer.NewApp()
		ctrl.SetupToPeer(app)

		ctrl.RemoteFileService = &remotefile.RemoteFileService{}
		return app, ctrl
	}

	t.Run("Search", func(t *testing.T) {

		app, ctrl := setup()

		// mock NodeFileService
		nodeFileService := mockedNodeFile.MockNodeFileInternalService{}
		defer nodeFileService.AssertExpectations(t)
		ctrl.RemoteFileService.NodeFileService = &nodeFileService

		var nodeFile nodefile.NodeFile
		nodeFile.ID = "nodeFileId"
		nodeFile.ItemID = 1
		nodeFile.Name = "test.txt"
		nodeFile.Size = 123
		nodeFile.FileType = nodeitem.FileTypeFile
		nodeFile.ParentPath = "parentPath"
		nodeFile.FilePath = "filePath"
		nodeFile.Available = true
		nodeFile.CreatedAt = time.Now()
		nodeFile.UpdatedAt = time.Now()

		var condition remotefile.RemoteFileRecordSearchCondition
		condition.ItemID = 1
		condition.ParentPath = &nodeFile.ParentPath

		nodeFileService.On("TraverseWithCondition", mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			traverseFn := args.Get(0).(func(item nodefile.NodeFile) error)
			traverseFn(nodeFile)
			condition_ := args.Get(1).(nodefile.NodeFileSearchCondition)
			assert.Equal(t, condition.ItemID, int32(condition_.ItemID))
			assert.Equal(t, *condition.ParentPath, *condition_.ParentPath)
		})

		// request and response
		reqBytes, err := proto.Marshal(&condition)
		assert.Nil(t, err)
		req := peer.NewRequest(remotefile.SearchRemoteFiles, bytes.NewReader(reqBytes))
		reqReader := peer.MarshalRequest(req)
		var ctx peer.Context
		peer.InitContext(&ctx)
		err = peer.UnmarshalRequest(reqReader, ctx.Request())
		assert.Nil(t, err)

		err = app.Run(&ctx, nil)
		assert.Nil(t, err)

		body, err := io.ReadAll(ctx)
		assert.Nil(t, err)

		var results remotefile.RemoteFileRecordList
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
		nodeFileService := mockedNodeFile.MockNodeFileInternalService{}
		defer nodeFileService.AssertExpectations(t)
		ctrl.RemoteFileService.NodeFileService = &nodeFileService

		var nodeFile nodefile.NodeFile
		nodeFile.ID = "nodeFileId"
		nodeFile.ItemID = 1
		nodeFile.Name = "test.txt"
		nodeFile.Size = 123
		nodeFile.FileType = nodeitem.FileTypeFile
		nodeFile.ParentPath = "parentPath"
		nodeFile.FilePath = "filePath"
		nodeFile.Available = true
		nodeFile.CreatedAt = time.Now()
		nodeFile.UpdatedAt = time.Now()

		nodeFileService.On("SelectWithCondition", mock.Anything).Once().Return(nodeFile, nil)

		var condition remotefile.RemoteFileRecordSelectCondition
		condition.ItemID = uint32(nodeFile.ItemID)
		condition.ParentPath = nodeFile.ParentPath
		condition.Name = nodeFile.Name

		// request and response
		reqBytes, err := proto.Marshal(&condition)
		assert.Nil(t, err)
		req := peer.NewRequest(remotefile.SelectRemoteFile, bytes.NewReader(reqBytes))
		reqReader := peer.MarshalRequest(req)
		var ctx peer.Context
		peer.InitContext(&ctx)
		err = peer.UnmarshalRequest(reqReader, ctx.Request())
		assert.Nil(t, err)

		err = app.Run(&ctx, nil)
		assert.Nil(t, err)

		body, err := io.ReadAll(ctx)
		assert.Nil(t, err)

		var results remotefile.RemoteFileRecord
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
