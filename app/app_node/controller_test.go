package appnode_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pan/app/web"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	appnode "pan/app/app_node"
	mocked "pan/mocks/pan/app/app_node"
	mockedPeer "pan/mocks/pan/app/peer"
)

func TestAppNodeController(t *testing.T) {

	setup := func() (web.WebApp, *appnode.AppNodeController) {
		ctrl := &appnode.AppNodeController{}
		webApp := web.NewWebApp()
		ctrl.SetupToWeb(webApp)

		ctrl.AppNodeService = &appnode.AppNodeService{}
		return webApp, ctrl
	}

	t.Run("GET /nodes", func(t *testing.T) {

		webApp, ctrl := setup()

		peerNodeRepo := &mocked.MockAppNodeRepository{}
		defer peerNodeRepo.AssertExpectations(t)
		ctrl.AppNodeService.AppNodeRepo = peerNodeRepo
		total := int64(10)
		peerIds := [][]byte{[]byte("peer id 1"), []byte("peer id 2"), []byte("peer id 3")}
		items := []appnode.AppNode{
			{ID: 1, Name: "peer node1", PeerID: appnode.EncodePeerID(peerIds[0]), Blocked: false},
			{ID: 2, Name: "peer node2", PeerID: appnode.EncodePeerID(peerIds[1]), Blocked: false},
			{ID: 2, Name: "peer node2", PeerID: appnode.EncodePeerID(peerIds[2]), Blocked: true},
		}
		peerNodeRepo.On("Search", appnode.AppNodeSearchCondition{}).Once().Return(total, items, nil)

		peerModule := &mockedPeer.MockPeerModule{}
		ctrl.AppNodeService.PeerModule = peerModule
		defer peerModule.AssertExpectations(t)
		peerModule.On("CanReach", peerIds[0]).Once().Return(true)
		peerModule.On("CanReach", peerIds[1]).Once().Return(false)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/nodes", nil)
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(web.CountHeaderName))
		var results []appnode.AppNode
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		item := items[0]
		item.Online = true
		assert.Equal(t, item, results[0])
		assert.Equal(t, items[1], results[1])
		assert.Equal(t, items[2], results[2])
	})

	t.Run("GET /nodes?q=keyword with query online", func(t *testing.T) {

		webApp, ctrl := setup()

		peerNodeRepo := &mocked.MockAppNodeRepository{}
		defer peerNodeRepo.AssertExpectations(t)
		ctrl.AppNodeService.AppNodeRepo = peerNodeRepo
		online := true
		condition := appnode.AppNodeSearchCondition{}
		condition.Keyword = "peer node1"
		condition.Online = &online
		condition.RangeStart = 0
		condition.RangeEnd = 12
		condition.SortField = "name"
		condition.SortOrder = "desc"

		total := int64(10)
		peerIds := [][]byte{[]byte("peer id 1"), []byte("peer id 2"), []byte("peer id 3")}
		items := []appnode.AppNode{
			{ID: 1, Name: "peer node1", PeerID: appnode.EncodePeerID(peerIds[0]), Blocked: false},
			{ID: 2, Name: "peer node2", PeerID: appnode.EncodePeerID(peerIds[1]), Blocked: false},
			{ID: 2, Name: "peer node2", PeerID: appnode.EncodePeerID(peerIds[2]), Blocked: true},
		}
		peerNodeRepo.On("Search", condition).Once().Return(total, items, nil)

		peerModule := &mockedPeer.MockPeerModule{}
		defer peerModule.AssertExpectations(t)
		ctrl.AppNodeService.PeerModule = peerModule
		peerModule.On("CanReach", peerIds[0]).Once().Return(true)
		peerModule.On("CanReach", peerIds[1]).Once().Return(false)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/nodes", nil)
		q := req.URL.Query()
		q.Add("q", condition.Keyword)
		q.Add("online", strconv.FormatBool(online))
		q.Add("_start", strconv.Itoa(condition.RangeStart))
		q.Add("_end", strconv.Itoa(condition.RangeEnd))
		q.Add("_sort", condition.SortField)
		q.Add("_order", condition.SortOrder)
		req.URL.RawQuery = q.Encode()
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(web.CountHeaderName))
		var results []appnode.AppNode
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		item := items[0]
		item.Online = true
		assert.Len(t, results, 1)
		assert.Equal(t, item, results[0])
	})

	t.Run("GET /nodes?q=keyword with blocked", func(t *testing.T) {

		webApp, ctrl := setup()

		peerNodeRepo := &mocked.MockAppNodeRepository{}
		defer peerNodeRepo.AssertExpectations(t)
		ctrl.AppNodeService.AppNodeRepo = peerNodeRepo
		blocked := true
		condition := appnode.AppNodeSearchCondition{}
		condition.Keyword = "peer node1"
		condition.Blocked = &blocked
		condition.RangeStart = 0
		condition.RangeEnd = 12
		condition.SortField = "name"
		condition.SortOrder = "desc"

		total := int64(10)
		peerIds := [][]byte{[]byte("peer id 1"), []byte("peer id 2"), []byte("peer id 3")}
		items := []appnode.AppNode{
			{ID: 1, Name: "peer node1", PeerID: appnode.EncodePeerID(peerIds[0]), Blocked: false},
			{ID: 2, Name: "peer node2", PeerID: appnode.EncodePeerID(peerIds[1]), Blocked: false},
			{ID: 2, Name: "peer node2", PeerID: appnode.EncodePeerID(peerIds[2]), Blocked: true},
		}
		peerNodeRepo.On("Search", condition).Once().Return(total, items, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/nodes", nil)
		q := req.URL.Query()
		q.Add("q", condition.Keyword)
		q.Add("blocked", strconv.FormatBool(blocked))
		q.Add("_start", strconv.Itoa(condition.RangeStart))
		q.Add("_end", strconv.Itoa(condition.RangeEnd))
		q.Add("_sort", condition.SortField)
		q.Add("_order", condition.SortOrder)
		req.URL.RawQuery = q.Encode()
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(web.CountHeaderName))
		var results []appnode.AppNode
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, items, results)
	})

	t.Run("GET /nodes/:id", func(t *testing.T) {

		webApp, ctrl := setup()

		peerNodeRepo := &mocked.MockAppNodeRepository{}
		defer peerNodeRepo.AssertExpectations(t)
		ctrl.AppNodeService.AppNodeRepo = peerNodeRepo
		peerId := []byte("peer id 1")
		item := appnode.AppNode{ID: 1, Name: "peer node1", PeerID: appnode.EncodePeerID(peerId), Blocked: false}
		peerNodeRepo.On("Select", item.ID).Once().Return(item, nil)

		peerModule := &mockedPeer.MockPeerModule{}
		defer peerModule.AssertExpectations(t)
		ctrl.AppNodeService.PeerModule = peerModule
		peerModule.On("CanReach", peerId).Once().Return(true)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req := httptest.NewRequest("GET", url, nil)
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result appnode.AppNode
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		item_ := item
		item_.Online = true
		assert.Equal(t, item_, result)
	})

	t.Run("GET /nodes/:id wth blocked", func(t *testing.T) {

		webApp, ctrl := setup()

		peerNodeRepo := &mocked.MockAppNodeRepository{}
		defer peerNodeRepo.AssertExpectations(t)
		ctrl.AppNodeService.AppNodeRepo = peerNodeRepo
		peerId := []byte("peer id 1")
		item := appnode.AppNode{ID: 1, Name: "peer node1", PeerID: appnode.EncodePeerID(peerId), Blocked: true}
		peerNodeRepo.On("Select", item.ID).Once().Return(item, nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req := httptest.NewRequest("GET", url, nil)
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result appnode.AppNode
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, item, result)
	})

	t.Run("DELETE /nodes/:id", func(t *testing.T) {
		webApp, ctrl := setup()

		peerNodeRepo := &mocked.MockAppNodeRepository{}
		defer peerNodeRepo.AssertExpectations(t)
		ctrl.AppNodeService.AppNodeRepo = peerNodeRepo
		peerId := []byte("peer id 1")
		item := appnode.AppNode{ID: 1, Name: "peer node1", PeerID: appnode.EncodePeerID(peerId), Blocked: false}
		peerNodeRepo.On("Select", item.ID).Once().Return(item, nil)
		peerNodeRepo.On("Delete", item).Once().Return(nil)

		peerModule := &mockedPeer.MockPeerModule{}
		defer peerModule.AssertExpectations(t)
		ctrl.AppNodeService.PeerModule = peerModule
		peerModule.On("Purge", peerId).Once().Return(nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req := httptest.NewRequest("DELETE", url, nil)
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
	})

	t.Run("DELETE /nodes/:id with blocked", func(t *testing.T) {
		webApp, ctrl := setup()

		peerNodeRepo := &mocked.MockAppNodeRepository{}
		defer peerNodeRepo.AssertExpectations(t)
		ctrl.AppNodeService.AppNodeRepo = peerNodeRepo
		peerId := []byte("peer id 1")
		item := appnode.AppNode{ID: 1, Name: "peer node1", PeerID: appnode.EncodePeerID(peerId), Blocked: true}
		peerNodeRepo.On("Select", item.ID).Once().Return(item, nil)
		peerNodeRepo.On("Delete", item).Once().Return(nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req := httptest.NewRequest("DELETE", url, nil)
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
	})

	t.Run("POST /nodes/:id", func(t *testing.T) {

		webApp, ctrl := setup()

		peerNodeRepo := &mocked.MockAppNodeRepository{}
		defer peerNodeRepo.AssertExpectations(t)
		ctrl.AppNodeService.AppNodeRepo = peerNodeRepo
		blocked := true
		fields := appnode.AppNodeFields{
			Name:    "peer node1",
			PeerID:  appnode.EncodePeerID([]byte("peer id 1")),
			Blocked: &blocked,
		}
		item := appnode.AppNode{Name: fields.Name, PeerID: fields.PeerID, Blocked: blocked}
		newItem := item
		newItem.ID = 1
		peerNodeRepo.On("Save", item).Once().Return(newItem, nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/nodes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		var result appnode.AppNode
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, newItem, result)
	})

	t.Run("PATCH /nodes/:id", func(t *testing.T) {

		webApp, ctrl := setup()

		peerNodeRepo := &mocked.MockAppNodeRepository{}
		defer peerNodeRepo.AssertExpectations(t)
		ctrl.AppNodeService.AppNodeRepo = peerNodeRepo
		blocked := true
		peerId := []byte("peer id")
		fields := appnode.AppNodeFields{
			Name:    "peer node1",
			PeerID:  appnode.EncodePeerID([]byte("peer id 1")),
			Blocked: &blocked,
		}
		item := appnode.AppNode{ID: 123, Name: "peer node", PeerID: appnode.EncodePeerID(peerId), Blocked: false}
		newItem := item
		newItem.Name = fields.Name
		newItem.Blocked = blocked
		peerNodeRepo.On("Select", item.ID).Once().Return(item, nil)
		peerNodeRepo.On("Save", newItem).Once().Return(newItem, nil)

		peerModule := &mockedPeer.MockPeerModule{}
		defer peerModule.AssertExpectations(t)
		ctrl.AppNodeService.PeerModule = peerModule
		peerModule.On("Purge", peerId).Once().Return(nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result appnode.AppNode
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, newItem, result)

	})

	t.Run("PATCH /nodes/:id with unblocked", func(t *testing.T) {

		webApp, ctrl := setup()

		peerNodeRepo := &mocked.MockAppNodeRepository{}
		defer peerNodeRepo.AssertExpectations(t)
		ctrl.AppNodeService.AppNodeRepo = peerNodeRepo
		blocked := false
		peerId := []byte("peer id")
		fields := appnode.AppNodeFields{
			Name:    "peer node1",
			PeerID:  appnode.EncodePeerID([]byte("peer id 1")),
			Blocked: &blocked,
		}
		item := appnode.AppNode{ID: 123, Name: "peer node", PeerID: appnode.EncodePeerID(peerId), Blocked: false}
		newItem := item
		newItem.Name = fields.Name
		newItem.Blocked = blocked
		peerNodeRepo.On("Select", item.ID).Once().Return(item, nil)
		peerNodeRepo.On("Save", newItem).Once().Return(newItem, nil)

		peerModule := &mockedPeer.MockPeerModule{}
		defer peerModule.AssertExpectations(t)
		ctrl.AppNodeService.PeerModule = peerModule
		peerModule.On("CanReach", peerId).Once().Return(true)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result appnode.AppNode
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		item_ := newItem
		item_.Online = true
		assert.Equal(t, item_, result)

	})
}
