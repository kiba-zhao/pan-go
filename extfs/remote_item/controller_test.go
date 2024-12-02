package remoteitem_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"pan/app/web"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	nodeitem "pan/extfs/node_item"
	remoteitem "pan/extfs/remote_item"
	mockedSample "pan/mocks/pan/app/sample"
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

	t.Run("GET /remote-items?peerId=", func(t *testing.T) {

		app, ctrl := setup()

		// mock SamplePeer
		samplePeer := &mockedSample.MockSamplePeer{}
		defer samplePeer.AssertExpectations(t)
		ctrl.RemoteItemService.RemoteItemBroker.SamplePeer = samplePeer
		//

		// mock response
		peerId := []byte("peerId")
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

		samplePeer.On("RequestWithProto", context.Background(), peerId, remoteitem.SelectAllRemoteItems, mock.Anything, nil).Once().Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(3).(*remoteitem.RemoteItemRecordList)
			err := proto.Unmarshal(resBody, resp)
			assert.Nil(t, err)
		})

		base64PeerID := base64.StdEncoding.EncodeToString(peerId)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/remote-items", nil)
		q := req.URL.Query()
		q.Add("peerId", base64PeerID)
		req.URL.RawQuery = q.Encode()
		app.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []remoteitem.RemoteItem
		err = json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(results))

		assert.Equal(t, uint(record.ID), results[0].ItemID)
		assert.Equal(t, base64PeerID, results[0].PeerID)
		assert.Equal(t, record.Name, results[0].Name)
		assert.Equal(t, record.FileType, results[0].FileType)
		assert.Equal(t, record.Size, results[0].Size)
		assert.Equal(t, record.Available, results[0].Available)
		assert.Equal(t, record.CreatedAt, results[0].CreatedAt.Unix())
		assert.Equal(t, record.UpdatedAt, results[0].UpdatedAt.Unix())

	})

}
