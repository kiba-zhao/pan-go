package remoteitem_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pan/lib/peer"
	"pan/lib/web"
	mockedFeature "pan/mocks/pan/lib/feature"
	"path"
	"testing"
	"time"

	nodeitem "pan/features/extfs/node_item"
	remoteitem "pan/features/extfs/remote_item"

	mockedNodeItem "pan/mocks/pan/features/extfs/node_item"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func TestRemoteFileInfoController(t *testing.T) {

	setup := func() (web.WebApp, *remoteitem.RemoteFileInfoController) {
		app := web.NewWebApp()
		ctrl := &remoteitem.RemoteFileInfoController{}
		ctrl.SetupToWeb(app)

		ctrl.RemoteFileInfoService = &remoteitem.RemoteFileInfoService{}
		ctrl.RemoteFileInfoService.RemoteFileInfoBroker = &remoteitem.RemoteFileInfoBroker{}
		return app, ctrl
	}

	t.Run("GET /remotes/:peerId/remote-items/:id/_files/*filepath", func(t *testing.T) {
		app, ctrl := setup()

		// mock BrokerHelper
		brokerHelper := &mockedFeature.MockBrokerHelper{}
		defer brokerHelper.AssertExpectations(t)
		ctrl.RemoteFileInfoService.RemoteFileInfoBroker.BrokerHelper = brokerHelper
		//

		peerIdBytes := []byte("peerId")
		var record remoteitem.RemoteFileInfoRecord
		record.ItemID = 1
		record.Name = "test.txt"
		record.Size = 123
		record.FileType = nodeitem.FileTypeFile
		record.ParentPath = "parentPath"
		record.Available = true
		record.CreatedAt = time.Now().Unix()
		record.UpdatedAt = time.Now().Unix()
		record.FilePath = path.Join(record.ParentPath, record.Name)

		resBody, err := proto.Marshal(&record)
		assert.Nil(t, err)

		brokerHelper.On("RequestWithProto", context.Background(), peerIdBytes, remoteitem.SelectRemoteFileInfo, mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(3).(*remoteitem.RemoteFileInfoRecord)
			err := proto.Unmarshal(resBody, resp)
			assert.Nil(t, err)

			condition := args.Get(4).(*remoteitem.RemoteFileInfoRecordSelectCondition)
			assert.Equal(t, record.ItemID, condition.ItemID)
			assert.Equal(t, record.FilePath, condition.FilePath)
		})

		// mock NodeFileInfoService
		nodeFileInfoService := &mockedNodeItem.MockNodeFileInfoInternalService{}
		defer nodeFileInfoService.AssertExpectations(t)
		ctrl.RemoteFileInfoService.NodeFileInfoService = nodeFileInfoService
		nodeFileInfoService.On("IsNotExist", mock.Anything).Once().Return(false)
		//

		peerId := peer.EncodePeerID(peerIdBytes)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/remotes/%s/remote-items/%d/_files/%s", peerId, record.ItemID, record.FilePath)
		req := httptest.NewRequest("GET", url, nil)
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp remoteitem.RemoteFileInfo
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Nil(t, err)

		assert.Equal(t, uint(record.ItemID), resp.ItemID)
		assert.Equal(t, record.Name, resp.Name)
		assert.Equal(t, record.Size, resp.Size)
		assert.Equal(t, record.FileType, resp.FileType)
		assert.Equal(t, record.ParentPath, resp.ParentPath)
		assert.Equal(t, record.Available, resp.Available)
		assert.Equal(t, record.CreatedAt, resp.CreatedAt.Unix())
		assert.Equal(t, record.UpdatedAt, resp.UpdatedAt.Unix())
		assert.Equal(t, record.FilePath, resp.FilePath)

	})

	t.Run("GET /remotes/:peerId/remote-items/:id/_files?parentPath=", func(t *testing.T) {
		app, ctrl := setup()

		// mock BrokerHelper
		brokerHelper := &mockedFeature.MockBrokerHelper{}
		defer brokerHelper.AssertExpectations(t)
		ctrl.RemoteFileInfoService.RemoteFileInfoBroker.BrokerHelper = brokerHelper
		//

		// mock response
		peerIdBytes := []byte("peerId")
		var record remoteitem.RemoteFileInfoRecord
		record.ID = "recordId"
		record.ItemID = 1
		record.Name = "test.txt"
		record.Size = 123
		record.FileType = nodeitem.FileTypeFile
		record.ParentPath = "parentPath"
		record.FilePath = "filePath"
		record.Available = true
		record.CreatedAt = time.Now().Unix()
		record.UpdatedAt = time.Now().Unix()

		var recordList remoteitem.RemoteFileInfoRecordList
		recordList.Items = append(recordList.Items, &record)
		resBody, err := proto.Marshal(&recordList)
		assert.Nil(t, err)

		brokerHelper.On("RequestWithProto", context.Background(), peerIdBytes, remoteitem.SearchRemoteFileInfos, mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(3).(*remoteitem.RemoteFileInfoRecordList)
			err := proto.Unmarshal(resBody, resp)
			assert.Nil(t, err)

			condition := args.Get(4).(*remoteitem.RemoteFileInfoRecordSearchCondition)
			assert.Equal(t, record.ItemID, condition.ItemID)
			assert.Equal(t, record.ParentPath, condition.ParentPath)
		})

		// mock NodeFileInfoService
		nodeFileInfoService := &mockedNodeItem.MockNodeFileInfoInternalService{}
		defer nodeFileInfoService.AssertExpectations(t)
		ctrl.RemoteFileInfoService.NodeFileInfoService = nodeFileInfoService
		nodeFileInfoService.On("IsNotExist", mock.Anything).Once().Return(false)
		//

		peerId := peer.EncodePeerID(peerIdBytes)
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/remotes/%s/remote-items/%d/_files", peerId, record.ItemID)
		req := httptest.NewRequest("GET", url, nil)
		q := req.URL.Query()
		q.Add("parentPath", record.ParentPath)
		req.URL.RawQuery = q.Encode()
		app.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []remoteitem.RemoteFileInfo
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
