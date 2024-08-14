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
	"strconv"
	"testing"

	mockedDispatcher "pan/mocks/pan/extfs/dispatchers"
	mockedRepo "pan/mocks/pan/extfs/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTargetController(t *testing.T) {

	setup := func() (web net.WebApp, ctrl *controllers.TargetController) {
		ctrl = new(controllers.TargetController)
		web = net.NewWebApp()
		ctrl.SetupToWeb(web)

		ctrl.TargetService = &services.TargetService{}
		ctrl.TargetService.TargetFileService = &services.TargetFileService{}
		return web, ctrl
	}

	t.Run("GET /targets", func(t *testing.T) {

		web, ctrl := setup()

		targetRepo := new(mockedRepo.MockTargetRepository)
		defer targetRepo.AssertExpectations(t)
		ctrl.TargetService.TargetRepo = targetRepo

		total := int64(10)
		enabled := false
		targets := []models.Target{
			{Name: "Target A", FilePath: "/path_a", Enabled: &enabled},
			{Name: "Target B", FilePath: "/path_b", Enabled: &enabled},
		}
		targetRepo.On("Search", models.TargetSearchCondition{}).Once().Return(total, targets, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/targets", nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(net.CountHeaderName))
		var results []models.Target
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, targets, results)
	})

	t.Run("GET /targets?q=keyword with query", func(t *testing.T) {

		web, ctrl := setup()

		targetRepo := new(mockedRepo.MockTargetRepository)
		defer targetRepo.AssertExpectations(t)
		ctrl.TargetService.TargetRepo = targetRepo

		enabled := false
		condition := models.TargetSearchCondition{}
		condition.Keyword = "keyword"
		condition.Enabled = &enabled
		condition.RangeStart = 0
		condition.RangeEnd = 12
		condition.SortField = "name"
		condition.SortOrder = "desc"

		total := int64(10)
		targets := []models.Target{
			{Name: "Target A", FilePath: "/path_a", Enabled: &enabled},
			{Name: "Target B", FilePath: "/path_b", Enabled: &enabled},
		}
		targetRepo.On("Search", condition).Once().Return(total, targets, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/targets", nil)
		q := req.URL.Query()
		q.Add("q", condition.Keyword)
		q.Add("enabled", strconv.FormatBool(enabled))
		q.Add("_start", strconv.Itoa(condition.RangeStart))
		q.Add("_end", strconv.Itoa(condition.RangeEnd))
		q.Add("_sort", condition.SortField)
		q.Add("_order", condition.SortOrder)
		req.URL.RawQuery = q.Encode()
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(net.CountHeaderName))
		var results []models.Target
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, targets, results)
	})

	t.Run("GET /targets with invalid query", func(t *testing.T) {

		web, _ := setup()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/targets", nil)
		q := req.URL.Query()
		q.Add("enabled", "enabled")
		q.Add("_start", "range-start")
		q.Add("_end", "range-end")
		q.Add("_sort", "sort-field")
		q.Add("_order", "sort-order")
		req.URL.RawQuery = q.Encode()
		web.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("GET /targets/:id", func(t *testing.T) {

		web, ctrl := setup()

		targetRepo := new(mockedRepo.MockTargetRepository)
		defer targetRepo.AssertExpectations(t)
		ctrl.TargetService.TargetRepo = targetRepo

		id := uint(123)
		enbaled := false
		var version *uint8
		target := models.Target{ID: id, Name: "Target A", FilePath: "/path_a", Enabled: &enbaled}
		targetRepo.On("Select", id, version).Once().Return(target, nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/targets/%d", id)
		req, _ := http.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result models.Target
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, target.ID, result.ID)
		assert.Equal(t, target.Name, result.Name)
		assert.Equal(t, target.FilePath, result.FilePath)
	})

	t.Run("POST /targets", func(t *testing.T) {

		web, ctrl := setup()

		targetRepo := new(mockedRepo.MockTargetRepository)
		defer targetRepo.AssertExpectations(t)
		ctrl.TargetService.TargetRepo = targetRepo

		enabled := true
		available := false
		version := uint8(1)
		fields := models.TargetFields{
			Name:     "Target A",
			FilePath: "/path_a",
			Enabled:  &enabled,
		}
		newTarget := models.Target{ID: 123, Name: fields.Name, FilePath: fields.FilePath, Enabled: fields.Enabled, Available: available, Version: version}
		targetRepo.On("Save", mock.AnythingOfType("models.Target"), false).Once().Return(newTarget, nil)

		targetDispatcher := new(mockedDispatcher.MockTargetDispatcher)
		defer targetDispatcher.AssertExpectations(t)
		ctrl.TargetService.TargetDispatcher = targetDispatcher

		targetDispatcher.On("Scan", newTarget).Once().Return(nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/targets", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		var result models.Target
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, newTarget, result)
	})

	t.Run("PATCH /targets/:id", func(t *testing.T) {

		web, ctrl := setup()

		targetRepo := new(mockedRepo.MockTargetRepository)
		defer targetRepo.AssertExpectations(t)
		ctrl.TargetService.TargetRepo = targetRepo

		id := uint(123)
		var version *uint8
		enabled := false
		available := false
		fields := models.TargetFields{
			Name:     "Target A",
			FilePath: "/path_a",
			Enabled:  &enabled,
		}
		firstVersion := uint8(1)
		target := models.Target{ID: id, Name: "Target B", FilePath: "/path_b", Enabled: fields.Enabled, Available: available, Version: firstVersion}
		sencodVersion := firstVersion + 1
		newTarget := models.Target{ID: target.ID, Name: fields.Name, FilePath: fields.FilePath, Enabled: fields.Enabled, Available: available, Version: sencodVersion}
		targetRepo.On("Select", id, version).Once().Return(target, nil)
		targetRepo.On("Save", mock.AnythingOfType("models.Target"), true).Once().Return(newTarget, nil)

		targetDispatcher := new(mockedDispatcher.MockTargetDispatcher)
		defer targetDispatcher.AssertExpectations(t)
		ctrl.TargetService.TargetDispatcher = targetDispatcher

		targetDispatcher.On("Scan", newTarget).Once().Return(nil)

		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/targets/%d", id)
		req, _ := http.NewRequest("PATCH", url, bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result models.Target
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, newTarget, result)
	})

	t.Run("DELETE /targets/:id", func(t *testing.T) {
		t.Skip("debug")
		web, ctrl := setup()

		targetRepo := new(mockedRepo.MockTargetRepository)
		defer targetRepo.AssertExpectations(t)
		ctrl.TargetService.TargetRepo = targetRepo

		id := uint(123)
		var version *uint8
		target := models.Target{ID: id, Name: "Target A", FilePath: "/path_a"}
		targetRepo.On("Select", id, version).Once().Return(target, nil)
		targetRepo.On("Delete", target).Once().Return(nil)

		targetDispatcher := new(mockedDispatcher.MockTargetDispatcher)
		defer targetDispatcher.AssertExpectations(t)
		ctrl.TargetService.TargetDispatcher = targetDispatcher

		targetDispatcher.On("Clean", target).Once().Return(nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/targets/%d", id)
		req, _ := http.NewRequest("DELETE", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
	})
}
