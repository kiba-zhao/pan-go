package remoteitem_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"pan/lib/peer"
	"pan/lib/web"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	nodeitem "pan/features/extfs/node_item"
	remoteitem "pan/features/extfs/remote_item"
	libApp "pan/lib/app"
	mockedFeature "pan/mocks/pan/lib/feature"
)

func TestRemoteItemController(t *testing.T) {

	setup := func() (web.WebApp, *remoteitem.RemoteItemController) {
		ctrl := &remoteitem.RemoteItemController{}
		app := web.NewWebApp()
		ctrl.SetupToWeb(app)

		ctrl.RemoteItemService = &remoteitem.RemoteItemService{}
		ctrl.RemoteItemService.RemoteItemBroker = &remoteitem.RemoteItemBroker{}
		return app, ctrl
	}

	t.Run("GET /remotes/:peerId/remote-items/:id", func(t *testing.T) {
		app, ctrl := setup()

		// mock for Broker
		brokerHelper := &mockedFeature.MockBrokerHelper{}
		defer brokerHelper.AssertExpectations(t)
		ctrl.RemoteItemService.RemoteItemBroker.InitBroker(brokerHelper)
		//

		// mock response
		peerIdBytes := []byte("peerId")
		var record remoteitem.RemoteItemRecord
		record.ID = 1
		record.Name = "test.txt"
		record.Size = 123
		record.FileType = nodeitem.FileTypeFile
		record.Available = true
		record.CreatedAt = time.Now().Unix()
		record.UpdatedAt = time.Now().Unix()

		resBody, err := proto.Marshal(&record)
		assert.Nil(t, err)

		appCtx := libApp.NewAppContext()
		appCtx.Respond(bytes.NewReader(resBody))

		brokerHelper.On("Do", context.Background(), peerIdBytes, mock.Anything).Once().Return(&appCtx.Response, nil).Run(func(args mock.Arguments) {
			peerReq := args.Get(2).(peer.PeerRequest)
			assert.Equal(t, remoteitem.SelectRemoteItem, peerReq.Name())

			bodyBytes, err := io.ReadAll(peerReq.Reader)
			assert.Nil(t, err)

			var condition remoteitem.RemoteItemRecordSelectCondition
			err = proto.Unmarshal(bodyBytes, &condition)
			assert.Nil(t, err)
			assert.Equal(t, record.ID, *condition.ID)
		})
		//

		peerId := peer.EncodePeerID(peerIdBytes)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/remotes/%s/remote-items/%d", peerId, record.ID)
		req := httptest.NewRequest("GET", url, nil)
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp remoteitem.RemoteItem
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Nil(t, err)

		assert.Equal(t, peerId, resp.PeerID)
		assert.Equal(t, uint(record.ID), resp.ItemID)
		assert.Equal(t, record.Name, resp.Name)
		assert.Equal(t, record.FileType, resp.FileType)
		assert.Equal(t, record.Size, resp.Size)
		assert.Equal(t, record.Available, resp.Available)
		assert.Equal(t, time.Unix(record.CreatedAt, 0), resp.CreatedAt)
		assert.Equal(t, time.Unix(record.UpdatedAt, 0), resp.UpdatedAt)
	})

	t.Run("GET /remotes/:peerId/remote-items", func(t *testing.T) {

		app, ctrl := setup()

		// mock BrokerHelper
		brokerHelper := &mockedFeature.MockBrokerHelper{}
		defer brokerHelper.AssertExpectations(t)
		ctrl.RemoteItemService.RemoteItemBroker.InitBroker(brokerHelper)
		//

		// mock response
		peerIdBytes := []byte("peerId")
		var record remoteitem.RemoteItemRecord
		record.ID = 1
		record.Name = "test.txt"
		record.Size = 123
		record.FileType = nodeitem.FileTypeFile
		record.Available = true
		record.CreatedAt = time.Now().Unix()
		record.UpdatedAt = time.Now().Unix()

		var recordList remoteitem.RemoteItemRecordList
		recordList.Items = append(recordList.Items, &record)
		resBody, err := proto.Marshal(&recordList)
		assert.Nil(t, err)

		appCtx := libApp.NewAppContext()
		appCtx.Respond(bytes.NewReader(resBody))

		brokerHelper.On("Do", context.Background(), peerIdBytes, mock.Anything).Once().Return(&appCtx.Response, nil).Run(func(args mock.Arguments) {
			peerReq := args.Get(2).(peer.PeerRequest)

			assert.Equal(t, remoteitem.SelectAllRemoteItems, peerReq.Name())
		})

		peerId := peer.EncodePeerID(peerIdBytes)
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/remotes/%s/remote-items", peerId)
		req := httptest.NewRequest("GET", url, nil)
		app.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []remoteitem.RemoteItem
		err = json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(results))

		assert.Equal(t, uint(record.ID), results[0].ItemID)
		assert.Equal(t, peerId, results[0].PeerID)
		assert.Equal(t, record.Name, results[0].Name)
		assert.Equal(t, record.FileType, results[0].FileType)
		assert.Equal(t, record.Size, results[0].Size)
		assert.Equal(t, record.Available, results[0].Available)
		assert.Equal(t, record.CreatedAt, results[0].CreatedAt.Unix())
		assert.Equal(t, record.UpdatedAt, results[0].UpdatedAt.Unix())

	})

}
