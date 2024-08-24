package controllers_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pan/app/controllers"
	"pan/app/models"
	"pan/app/net"
	"pan/app/services"
	"strconv"
	"testing"

	mockedNode "pan/mocks/pan/app/node"
	mockedRepo "pan/mocks/pan/app/repositories"
	mockedServices "pan/mocks/pan/app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNode(t *testing.T) {

	setup := func() (net.WebApp, *controllers.NodeController) {
		ctrl := &controllers.NodeController{}
		web := net.NewWebApp()
		ctrl.SetupToWeb(web)

		ctrl.NodeService = &services.NodeService{}
		return web, ctrl
	}

	t.Run("GET /nodes", func(t *testing.T) {

		web, ctrl := setup()

		nodeRepo := &mockedRepo.MockNodeRepository{}
		defer nodeRepo.AssertExpectations(t)
		ctrl.NodeService.NodeRepo = nodeRepo
		total := int64(10)
		nodeIds := [][]byte{[]byte("node id 1"), []byte("node id 2"), []byte("node id 3")}
		items := []models.Node{
			{ID: 1, Name: "node1", NodeID: base64.StdEncoding.EncodeToString(nodeIds[0]), Blocked: false},
			{ID: 2, Name: "node2", NodeID: base64.StdEncoding.EncodeToString(nodeIds[1]), Blocked: false},
			{ID: 2, Name: "node2", NodeID: base64.StdEncoding.EncodeToString(nodeIds[2]), Blocked: true},
		}
		nodeRepo.On("Search", models.NodeSearchCondition{}).Once().Return(total, items, nil)

		mgr := &mockedNode.MockNodeManager{}
		defer mgr.AssertExpectations(t)
		mgr.On("Count", nodeIds[0]).Once().Return(3)
		mgr.On("Count", nodeIds[1]).Once().Return(0)
		provider := &mockedServices.MockNodeManagerProvider{}
		defer provider.AssertExpectations(t)
		ctrl.NodeService.Provider = provider
		provider.On("NodeManager").Once().Return(mgr)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/nodes", nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(net.CountHeaderName))
		var results []models.Node
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		item := items[0]
		item.Online = true
		assert.Equal(t, item, results[0])
		assert.Equal(t, items[1], results[1])
		assert.Equal(t, items[2], results[2])
	})

	t.Run("GET /nodes?q=keyword with query online", func(t *testing.T) {

		web, ctrl := setup()

		nodeRepo := &mockedRepo.MockNodeRepository{}
		defer nodeRepo.AssertExpectations(t)
		ctrl.NodeService.NodeRepo = nodeRepo
		online := true
		condition := models.NodeSearchCondition{}
		condition.Keyword = "node1"
		condition.Online = &online
		condition.RangeStart = 0
		condition.RangeEnd = 12
		condition.SortField = "name"
		condition.SortOrder = "desc"

		total := int64(10)
		nodeIds := [][]byte{[]byte("node id 1"), []byte("node id 2"), []byte("node id 3")}
		items := []models.Node{
			{ID: 1, Name: "node1", NodeID: base64.StdEncoding.EncodeToString(nodeIds[0]), Blocked: false},
			{ID: 2, Name: "node2", NodeID: base64.StdEncoding.EncodeToString(nodeIds[1]), Blocked: false},
			{ID: 2, Name: "node2", NodeID: base64.StdEncoding.EncodeToString(nodeIds[2]), Blocked: true},
		}
		nodeRepo.On("Search", condition).Once().Return(total, items, nil)

		mgr := &mockedNode.MockNodeManager{}
		defer mgr.AssertExpectations(t)
		mgr.On("Count", nodeIds[0]).Once().Return(3)
		mgr.On("Count", nodeIds[1]).Once().Return(0)
		provider := &mockedServices.MockNodeManagerProvider{}
		defer provider.AssertExpectations(t)
		ctrl.NodeService.Provider = provider
		provider.On("NodeManager").Once().Return(mgr)

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
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(net.CountHeaderName))
		var results []models.Node
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		item := items[0]
		item.Online = true
		assert.Len(t, results, 1)
		assert.Equal(t, item, results[0])
	})

	t.Run("GET /nodes?q=keyword with blocked", func(t *testing.T) {

		web, ctrl := setup()

		nodeRepo := &mockedRepo.MockNodeRepository{}
		defer nodeRepo.AssertExpectations(t)
		ctrl.NodeService.NodeRepo = nodeRepo
		blocked := true
		condition := models.NodeSearchCondition{}
		condition.Keyword = "node1"
		condition.Blocked = &blocked
		condition.RangeStart = 0
		condition.RangeEnd = 12
		condition.SortField = "name"
		condition.SortOrder = "desc"

		total := int64(10)
		nodeIds := [][]byte{[]byte("node id 1"), []byte("node id 2"), []byte("node id 3")}
		items := []models.Node{
			{ID: 1, Name: "node1", NodeID: base64.StdEncoding.EncodeToString(nodeIds[0]), Blocked: false},
			{ID: 2, Name: "node2", NodeID: base64.StdEncoding.EncodeToString(nodeIds[1]), Blocked: false},
			{ID: 2, Name: "node2", NodeID: base64.StdEncoding.EncodeToString(nodeIds[2]), Blocked: true},
		}
		nodeRepo.On("Search", condition).Once().Return(total, items, nil)

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
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(net.CountHeaderName))
		var results []models.Node
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, items, results)
	})

	t.Run("GET /nodes/:id", func(t *testing.T) {

		web, ctrl := setup()

		nodeRepo := &mockedRepo.MockNodeRepository{}
		defer nodeRepo.AssertExpectations(t)
		ctrl.NodeService.NodeRepo = nodeRepo
		nodeId := []byte("node id 1")
		item := models.Node{ID: 1, Name: "node1", NodeID: base64.StdEncoding.EncodeToString(nodeId), Blocked: false}
		nodeRepo.On("Select", item.ID).Once().Return(item, nil)

		mgr := &mockedNode.MockNodeManager{}
		defer mgr.AssertExpectations(t)
		mgr.On("Count", nodeId).Once().Return(1)
		provider := &mockedServices.MockNodeManagerProvider{}
		defer provider.AssertExpectations(t)
		ctrl.NodeService.Provider = provider
		provider.On("NodeManager").Once().Return(mgr)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req := httptest.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result models.Node
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		item_ := item
		item_.Online = true
		assert.Equal(t, item_, result)
	})

	t.Run("GET /nodes/:id wth blocked", func(t *testing.T) {

		web, ctrl := setup()

		nodeRepo := &mockedRepo.MockNodeRepository{}
		defer nodeRepo.AssertExpectations(t)
		ctrl.NodeService.NodeRepo = nodeRepo
		nodeId := []byte("node id 1")
		item := models.Node{ID: 1, Name: "node1", NodeID: base64.StdEncoding.EncodeToString(nodeId), Blocked: true}
		nodeRepo.On("Select", item.ID).Once().Return(item, nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req := httptest.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result models.Node
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, item, result)
	})

	t.Run("DELETE /nodes/:id", func(t *testing.T) {
		web, ctrl := setup()

		nodeRepo := &mockedRepo.MockNodeRepository{}
		defer nodeRepo.AssertExpectations(t)
		ctrl.NodeService.NodeRepo = nodeRepo
		nodeId := []byte("node id 1")
		item := models.Node{ID: 1, Name: "node1", NodeID: base64.StdEncoding.EncodeToString(nodeId), Blocked: false}
		nodeRepo.On("Select", item.ID).Once().Return(item, nil)
		nodeRepo.On("Delete", item).Once().Return(nil)

		mgr := &mockedNode.MockNodeManager{}
		defer mgr.AssertExpectations(t)
		mgr.On("TraverseNode", nodeId, mock.Anything).Once().Return(nil)

		provider := &mockedServices.MockNodeManagerProvider{}
		defer provider.AssertExpectations(t)
		ctrl.NodeService.Provider = provider
		provider.On("NodeManager").Once().Return(mgr)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req := httptest.NewRequest("DELETE", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
	})

	t.Run("DELETE /nodes/:id with blocked", func(t *testing.T) {
		web, ctrl := setup()

		nodeRepo := &mockedRepo.MockNodeRepository{}
		defer nodeRepo.AssertExpectations(t)
		ctrl.NodeService.NodeRepo = nodeRepo
		nodeId := []byte("node id 1")
		item := models.Node{ID: 1, Name: "node1", NodeID: base64.StdEncoding.EncodeToString(nodeId), Blocked: true}
		nodeRepo.On("Select", item.ID).Once().Return(item, nil)
		nodeRepo.On("Delete", item).Once().Return(nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req := httptest.NewRequest("DELETE", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
	})

	t.Run("POST /nodes/:id", func(t *testing.T) {

		web, ctrl := setup()

		nodeRepo := &mockedRepo.MockNodeRepository{}
		defer nodeRepo.AssertExpectations(t)
		ctrl.NodeService.NodeRepo = nodeRepo
		blocked := true
		fields := models.NodeFields{
			Name:    "node1",
			NodeID:  base64.StdEncoding.EncodeToString([]byte("node id 1")),
			Blocked: &blocked,
		}
		item := models.Node{Name: fields.Name, NodeID: fields.NodeID, Blocked: blocked}
		newItem := item
		newItem.ID = 1
		nodeRepo.On("Save", item).Once().Return(newItem, nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/nodes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		var result models.Node
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, newItem, result)
	})

	t.Run("PATCH /nodes/:id", func(t *testing.T) {

		web, ctrl := setup()

		nodeRepo := &mockedRepo.MockNodeRepository{}
		defer nodeRepo.AssertExpectations(t)
		ctrl.NodeService.NodeRepo = nodeRepo
		blocked := true
		nodeId := []byte("node id")
		fields := models.NodeFields{
			Name:    "node1",
			NodeID:  base64.StdEncoding.EncodeToString([]byte("node id 1")),
			Blocked: &blocked,
		}
		item := models.Node{ID: 123, Name: "node", NodeID: base64.StdEncoding.EncodeToString(nodeId), Blocked: false}
		newItem := item
		newItem.Name = fields.Name
		newItem.Blocked = blocked
		nodeRepo.On("Select", item.ID).Once().Return(item, nil)
		nodeRepo.On("Save", newItem).Once().Return(newItem, nil)

		mgr := &mockedNode.MockNodeManager{}
		defer mgr.AssertExpectations(t)
		mgr.On("TraverseNode", nodeId, mock.Anything).Once().Return(nil)

		provider := &mockedServices.MockNodeManagerProvider{}
		defer provider.AssertExpectations(t)
		ctrl.NodeService.Provider = provider
		provider.On("NodeManager").Once().Return(mgr)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result models.Node
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, newItem, result)

	})

	t.Run("PATCH /nodes/:id with unblocked", func(t *testing.T) {

		web, ctrl := setup()

		nodeRepo := &mockedRepo.MockNodeRepository{}
		defer nodeRepo.AssertExpectations(t)
		ctrl.NodeService.NodeRepo = nodeRepo
		blocked := false
		nodeId := []byte("node id")
		fields := models.NodeFields{
			Name:    "node1",
			NodeID:  base64.StdEncoding.EncodeToString([]byte("node id 1")),
			Blocked: &blocked,
		}
		item := models.Node{ID: 123, Name: "node", NodeID: base64.StdEncoding.EncodeToString(nodeId), Blocked: false}
		newItem := item
		newItem.Name = fields.Name
		newItem.Blocked = blocked
		nodeRepo.On("Select", item.ID).Once().Return(item, nil)
		nodeRepo.On("Save", newItem).Once().Return(newItem, nil)

		mgr := &mockedNode.MockNodeManager{}
		defer mgr.AssertExpectations(t)
		mgr.On("Count", nodeId).Once().Return(4)
		provider := &mockedServices.MockNodeManagerProvider{}
		defer provider.AssertExpectations(t)
		ctrl.NodeService.Provider = provider
		provider.On("NodeManager").Once().Return(mgr)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/nodes/%d", item.ID)
		req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result models.Node
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		item_ := newItem
		item_.Online = true
		assert.Equal(t, item_, result)

	})
}
