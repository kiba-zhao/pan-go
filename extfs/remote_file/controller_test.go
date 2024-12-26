package remotefile_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	appnode "pan/app/app_node"
	"pan/app/web"
	nodeitem "pan/extfs/node_item"
	remotefile "pan/extfs/remote_file"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	mockedSample "pan/mocks/pan/app/sample"
)

func TestRemoteFileController(t *testing.T) {

	setup := func() (web.WebApp, *remotefile.RemoteFileController) {
		ctrl := &remotefile.RemoteFileController{}
		app := web.NewWebApp()
		ctrl.SetupToWeb(app)

		ctrl.RemoteFileService = &remotefile.RemoteFileService{}
		ctrl.RemoteFileService.RemoteFileBroker = &remotefile.RemoteFileBroker{}
		return app, ctrl
	}

	t.Run("GET /remote-files?peerId=&itemId=", func(t *testing.T) {

		app, ctrl := setup()

		// mock SamplePeer
		samplePeer := &mockedSample.MockSamplePeer{}
		defer samplePeer.AssertExpectations(t)
		ctrl.RemoteFileService.RemoteFileBroker.SamplePeer = samplePeer
		//

		// mock response
		peerId := []byte("peerId")
		var record remotefile.RemoteFileRecord
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

		var recordList remotefile.RemoteFileRecordList
		recordList.Items = append(recordList.Items, &record)
		resBody, err := proto.Marshal(&recordList)
		assert.Nil(t, err)

		samplePeer.On("RequestWithProto", context.Background(), peerId, remotefile.SearchRemoteFiles, mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(3).(*remotefile.RemoteFileRecordList)
			err := proto.Unmarshal(resBody, resp)
			assert.Nil(t, err)

			condition := args.Get(4).(*remotefile.RemoteFileRecordSearchCondition)
			assert.Equal(t, record.ItemID, condition.ItemID)
		})

		base64PeerID := appnode.EncodePeerID(peerId)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/remote-files", nil)
		q := req.URL.Query()
		q.Add("peerId", base64PeerID)
		q.Add("itemId", strconv.FormatUint(uint64(record.ItemID), 10))
		req.URL.RawQuery = q.Encode()
		app.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []remotefile.RemoteFile
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
