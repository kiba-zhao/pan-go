package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"pan/app/net"
	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/services"
	"testing"

	mockedRepo "pan/mocks/pan/extfs/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodeItemController(t *testing.T) {

	setup := func() (web net.WebApp, ctrl *controllers.NodeItemController) {
		ctrl = new(controllers.NodeItemController)
		web = net.NewWebApp()
		ctrl.SetupToWeb(web)

		ctrl.NodeItemService = &services.NodeItemService{}
		return web, ctrl
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

		nodeItemRepo := &mockedRepo.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		id := uint(1)
		enabled := false
		entity := models.NodeItem{}
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
		var result models.NodeItem
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, entity, result)
	})

	t.Run("PATCH /node-items/:id", func(t *testing.T) {

		web, ctrl := setup()
		filepath, err := setupTemp("extfs-node-items")
		assert.Nil(t, err)
		defer teardownTemp(filepath)

		stat, err := os.Stat(filepath)
		assert.Nil(t, err)
		assert.True(t, stat.IsDir())

		nodeItemRepo := &mockedRepo.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		id := uint(1)
		enabledField := false
		fields := models.NodeItemFields{}
		fields.Enabled = &enabledField
		fields.FilePath = filepath
		fields.Name = "Node Item Field Name"
		enabled := !enabledField
		entity := models.NodeItem{}
		entity.ID = id
		entity.Enabled = &enabled
		entity.Name = "Node Item Name"
		entity.FilePath = "/path_a"
		nodeItemRepo.On("Select", id).Once().Return(entity, nil)
		newEntity := entity
		newEntity.Enabled = &enabled
		newEntity.FilePath = fields.FilePath
		newEntity.Name = fields.Name
		newEntity.FileType = services.FileTypeFolder
		nodeItemRepo.On("Save", mock.AnythingOfType("models.NodeItem")).Once().Return(newEntity, nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/node-items/%d", id)
		req, _ := http.NewRequest("PATCH", url, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result models.NodeItem
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

		nodeItemRepo := &mockedRepo.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		id := uint(1)
		enabled := false
		entity := models.NodeItem{}
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
		filepath, err := setupTemp("extfs-node-items")
		assert.Nil(t, err)
		defer teardownTemp(filepath)

		stat, err := os.Stat(filepath)
		assert.Nil(t, err)
		assert.True(t, stat.IsDir())

		nodeItemRepo := &mockedRepo.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		enabled := true
		fields := models.NodeItemFields{}
		fields.Enabled = &enabled
		fields.FilePath = filepath
		fields.Name = "Node Item Field Name"
		entity := models.NodeItem{}
		entity.ID = 1
		entity.Name = fields.Name
		entity.FilePath = fields.FilePath
		entity.FileType = services.FileTypeFolder
		entity.Enabled = fields.Enabled

		nodeItemRepo.On("Save", mock.AnythingOfType("models.NodeItem")).Once().Return(entity, nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/node-items", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		var result models.NodeItem
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

		filepath, err := setupTemp("extfs-node-items")
		assert.Nil(t, err)
		defer teardownTemp(filepath)

		stat, err := os.Stat(filepath)
		assert.Nil(t, err)
		assert.True(t, stat.IsDir())

		nodeItemRepo := &mockedRepo.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		enabled := true
		entity := models.NodeItem{}
		entity.ID = 1
		entity.Name = "Node Item Name"
		entity.FilePath = filepath
		entity.FileType = services.FileTypeFolder
		entity.Enabled = &enabled
		nodeItemRepo.On("TraverseAll", mock.AnythingOfType("func(models.NodeItem) error")).Once().Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(0).(func(models.NodeItem) error)
			fn(entity)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/node-items", nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result []models.NodeItem
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
