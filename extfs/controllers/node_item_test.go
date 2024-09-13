package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

		nodeItemRepo := &mockedRepo.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		id := uint(1)
		enabledField := false
		fields := models.NodeItemFields{}
		fields.Enabled = &enabledField
		fields.FilePath = "/path_field"
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
		nodeItemRepo.On("Save", mock.AnythingOfType("models.NodeItem")).Once().Return(newEntity, nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/node-items/%d", id)
		req, _ := http.NewRequest("PATCH", url, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result models.NodeItem
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, newEntity, result)
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

		nodeItemRepo := &mockedRepo.MockNodeItemRepository{}
		defer nodeItemRepo.AssertExpectations(t)
		ctrl.NodeItemService.NodeItemRepo = nodeItemRepo

		enabled := true
		fields := models.NodeItemFields{}
		fields.Enabled = &enabled
		fields.FilePath = "/path_field"
		fields.Name = "Node Item Field Name"
		entity := models.NodeItem{}
		entity.ID = 1
		entity.Name = fields.Name
		entity.FilePath = fields.FilePath
		entity.Enabled = fields.Enabled

		nodeItemRepo.On("Save", mock.AnythingOfType("models.NodeItem")).Once().Return(entity, nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/node-items", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		var result models.NodeItem
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, entity, result)
	})
}
