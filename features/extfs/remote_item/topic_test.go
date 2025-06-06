package remoteitem_test

import (
	"bytes"
	"io"
	libApp "pan/lib/app"
	"pan/lib/peer"
	"testing"
	"time"

	nodeitem "pan/features/extfs/node_item"
	remoteitem "pan/features/extfs/remote_item"
	mockedNodeItem "pan/mocks/pan/features/extfs/node_item"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func TestReRemoteItemTopic(t *testing.T) {

	setup := func() (peer.PeerApp, *remoteitem.RemoteItemTopic) {
		topic := &remoteitem.RemoteItemTopic{}
		app := libApp.NewApp()
		topic.SetupToPeer(app)

		topic.RemoteItemService = &remoteitem.RemoteItemService{}
		return app, topic
	}

	t.Run("SelectAll", func(t *testing.T) {

		app, ctrl := setup()

		// mock NodeItemService
		nodeItemService := mockedNodeItem.MockNodeItemInternalService{}
		defer nodeItemService.AssertExpectations(t)
		ctrl.RemoteItemService.NodeItemService = &nodeItemService

		var nodeItem nodeitem.NodeItem
		nodeItem.ID = 1
		nodeItem.Name = "test.txt"
		nodeItem.FileType = nodeitem.FileTypeFile
		nodeItem.Size = 123
		nodeItem.Available = true
		nodeItem.CreatedAt = time.Now()
		nodeItem.UpdatedAt = time.Now()

		nodeItemService.On("TraverseAll", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			traverseFn := args.Get(0).(func(nodeitem.NodeItem) error)
			traverseFn(nodeItem)
		})

		// request and response
		req := libApp.NewRequest(remoteitem.SelectAllRemoteItems, nil)
		reqReader := libApp.MarshalRequest(req)
		ctx := libApp.NewAppContext()
		libApp.InitContext(ctx)
		err := libApp.UnmarshalRequest(reqReader, ctx.Request())
		assert.Nil(t, err)

		err = app.Run(ctx, nil)
		assert.Nil(t, err)

		body, err := io.ReadAll(ctx)
		assert.Nil(t, err)

		var results remoteitem.RemoteItemRecordList
		err = proto.Unmarshal(body, &results)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(results.Items))
		assert.Equal(t, uint32(nodeItem.ID), results.Items[0].ID)
		assert.Equal(t, nodeItem.Name, results.Items[0].Name)
		assert.Equal(t, nodeItem.FileType, results.Items[0].FileType)
		assert.Equal(t, nodeItem.Size, results.Items[0].Size)
		assert.Equal(t, nodeItem.Available, results.Items[0].Available)
		assert.Equal(t, nodeItem.CreatedAt.Unix(), results.Items[0].CreatedAt)
		assert.Equal(t, nodeItem.UpdatedAt.Unix(), results.Items[0].UpdatedAt)
	})

	t.Run("Select", func(t *testing.T) {
		app, ctrl := setup()

		// mock NodeItemService
		nodeItemService := mockedNodeItem.MockNodeItemInternalService{}
		defer nodeItemService.AssertExpectations(t)
		ctrl.RemoteItemService.NodeItemService = &nodeItemService

		var nodeItem nodeitem.NodeItem
		nodeItem.ID = 1
		nodeItem.Name = "test.txt"
		nodeItem.FileType = nodeitem.FileTypeFile
		nodeItem.Size = 123
		nodeItem.Available = true
		nodeItem.CreatedAt = time.Now()
		nodeItem.UpdatedAt = time.Now()

		nodeItemService.On("SelectByName", nodeItem.Name).Once().Return(nodeItem, nil)

		// request and response
		var condition remoteitem.RemoteItemRecordSelectCondition
		condition.Name = &nodeItem.Name
		reqBody, err := proto.Marshal(&condition)
		assert.Nil(t, err)
		req := libApp.NewRequest(remoteitem.SelectRemoteItem, bytes.NewReader(reqBody))
		reqReader := libApp.MarshalRequest(req)
		ctx := libApp.NewAppContext()
		libApp.InitContext(ctx)
		err = libApp.UnmarshalRequest(reqReader, ctx.Request())
		assert.Nil(t, err)

		err = app.Run(ctx, nil)
		assert.Nil(t, err)

		body, err := io.ReadAll(ctx)
		assert.Nil(t, err)

		var result remoteitem.RemoteItemRecord
		err = proto.Unmarshal(body, &result)
		assert.Nil(t, err)

		assert.Equal(t, uint32(nodeItem.ID), result.ID)
		assert.Equal(t, nodeItem.Name, result.Name)
		assert.Equal(t, nodeItem.FileType, result.FileType)
		assert.Equal(t, nodeItem.Size, result.Size)
		assert.Equal(t, nodeItem.Available, result.Available)
		assert.Equal(t, nodeItem.CreatedAt.Unix(), result.CreatedAt)
		assert.Equal(t, nodeItem.UpdatedAt.Unix(), result.UpdatedAt)

	})
}
