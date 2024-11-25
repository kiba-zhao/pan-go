package nodeitem_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"pan/app/web"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	nodeitem "pan/extfs/node_item"
	mocked "pan/mocks/pan/extfs/node_item"
)

func TestNodeItemController(t *testing.T) {

	setup := func() (web.WebApp, *nodeitem.NodeItemController) {
		ctrl := new(nodeitem.NodeItemController)
		webApp := web.NewWebApp()
		ctrl.SetupToWeb(webApp)

		ctrl.NodeItemService = &nodeitem.NodeItemService{}
		return webApp, ctrl
	}

	setupTemp := func(name string) (string, error) {
		dir, err := os.MkdirTemp(os.TempDir(), name)
		return dir, err
	}

	teardownTemp := func(dir string) error {
		return os.RemoveAll(dir)
	}

	t.Run("GET /node-items/:id", func(t *testing.T) {

		web, ctrl := setup()

		nodeItemRepo := &mocked.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		id := uint(1)
		enabled := false
		entity := nodeitem.NodeItem{}
		entity.ID = id
		entity.Enabled = &enabled
		entity.Name = "Node Item Name"
		entity.FilePath = "/path_a"
		nodeItemRepo.On("Select", id).Once().Return(entity, nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/node-items/%d", id)
		req, _ := http.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result nodeitem.NodeItem
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, entity, result)
	})

	t.Run("PATCH /node-items/:id", func(t *testing.T) {

		web, ctrl := setup()
		filePath, err := setupTemp("extfs-node-items")
		assert.Nil(t, err)
		defer teardownTemp(filePath)

		stat, err := os.Stat(filePath)
		assert.Nil(t, err)
		assert.True(t, stat.IsDir())

		nodeItemRepo := &mocked.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		id := uint(1)
		enabledField := false
		fields := nodeitem.NodeItemFields{}
		fields.Enabled = &enabledField
		fields.FilePath = filePath
		fields.Name = "Node Item Field Name"
		enabled := !enabledField
		entity := nodeitem.NodeItem{}
		entity.ID = id
		entity.Enabled = &enabled
		entity.Name = "Node Item Name"
		entity.FilePath = "/path_a"
		nodeItemRepo.On("Select", id).Once().Return(entity, nil)
		newEntity := entity
		newEntity.Enabled = &enabled
		newEntity.FilePath = fields.FilePath
		newEntity.Name = fields.Name
		newEntity.FileType = nodeitem.FileTypeFolder
		nodeItemRepo.On("Save", mock.AnythingOfType("nodeitem.NodeItem")).Once().Return(newEntity, nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/node-items/%d", id)
		req, _ := http.NewRequest("PATCH", url, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result nodeitem.NodeItem
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, newEntity.ID, result.ID)
		assert.Equal(t, newEntity.Enabled, result.Enabled)
		assert.Equal(t, newEntity.Name, result.Name)
		assert.Equal(t, newEntity.FilePath, result.FilePath)
		assert.Equal(t, newEntity.FileType, result.FileType)
		assert.Equal(t, stat.Size(), result.Size)
		assert.Equal(t, newEntity.TagQuantity, result.TagQuantity)
		assert.Equal(t, newEntity.PendingTagQuantity, result.PendingTagQuantity)
	})

	t.Run("DELETE /node-items/:id", func(t *testing.T) {

		web, ctrl := setup()

		nodeItemRepo := &mocked.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		id := uint(1)
		enabled := false
		entity := nodeitem.NodeItem{}
		entity.ID = id
		entity.Name = "Node Item Name"
		entity.FilePath = "/path_a"
		entity.Enabled = &enabled
		nodeItemRepo.On("Select", id).Once().Return(entity, nil)
		nodeItemRepo.On("Delete", entity).Once().Return(nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/node-items/%d", id)
		req, _ := http.NewRequest("DELETE", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
	})

	t.Run("POST /node-items", func(t *testing.T) {

		web, ctrl := setup()
		filePath, err := setupTemp("extfs-node-items")
		assert.Nil(t, err)
		defer teardownTemp(filePath)

		stat, err := os.Stat(filePath)
		assert.Nil(t, err)
		assert.True(t, stat.IsDir())

		nodeItemRepo := &mocked.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		enabled := true
		fields := nodeitem.NodeItemFields{}
		fields.Enabled = &enabled
		fields.FilePath = filePath
		fields.Name = "Node Item Field Name"
		entity := nodeitem.NodeItem{}
		entity.ID = 1
		entity.Name = fields.Name
		entity.FilePath = fields.FilePath
		entity.FileType = nodeitem.FileTypeFolder
		entity.Enabled = fields.Enabled

		nodeItemRepo.On("Save", mock.AnythingOfType("nodeitem.NodeItem")).Once().Return(entity, nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/node-items", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		var result nodeitem.NodeItem
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, entity.ID, result.ID)
		assert.Equal(t, entity.Name, result.Name)
		assert.Equal(t, entity.FilePath, result.FilePath)
		assert.Equal(t, entity.Enabled, result.Enabled)
		assert.Equal(t, entity.FileType, result.FileType)
		assert.True(t, result.Available)
		assert.Equal(t, stat.Size(), result.Size)
		assert.Equal(t, entity.UpdatedAt, result.UpdatedAt)
		assert.Equal(t, entity.CreatedAt, result.CreatedAt)
		assert.Equal(t, entity.TagQuantity, result.TagQuantity)
		assert.Equal(t, entity.PendingTagQuantity, result.PendingTagQuantity)
	})

	t.Run("GET /node-items", func(t *testing.T) {

		web, ctrl := setup()

		filePath, err := setupTemp("extfs-node-items")
		assert.Nil(t, err)
		defer teardownTemp(filePath)

		stat, err := os.Stat(filePath)
		assert.Nil(t, err)
		assert.True(t, stat.IsDir())

		nodeItemRepo := &mocked.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		enabled := true
		entity := nodeitem.NodeItem{}
		entity.ID = 1
		entity.Name = "Node Item Name"
		entity.FilePath = filePath
		entity.FileType = nodeitem.FileTypeFolder
		entity.Enabled = &enabled
		nodeItemRepo.On("TraverseAll", mock.AnythingOfType("func(nodeitem.NodeItem) error")).Once().Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(nodeitem.NodeItem) error)
			fn(entity)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/node-items", nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result []nodeitem.NodeItem
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, entity.ID, result[0].ID)
		assert.Equal(t, entity.Name, result[0].Name)
		assert.Equal(t, entity.FilePath, result[0].FilePath)
		assert.Equal(t, entity.Enabled, result[0].Enabled)
		assert.Equal(t, entity.FileType, result[0].FileType)
		assert.True(t, result[0].Available)
		assert.Equal(t, stat.Size(), result[0].Size)
	})
}
