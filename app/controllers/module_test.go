package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pan/app/controllers"
	"pan/app/models"
	"pan/app/services"
	"pan/core"
	"testing"

	appTestMocked "pan/mocks/pan/app/test"
	mocked "pan/mocks/pan/core"

	"github.com/stretchr/testify/assert"
)

func TestModules(t *testing.T) {

	setup := func() (web *core.WebApp, ctrl *controllers.ModuleController) {
		ctrl = new(controllers.ModuleController)
		web = core.NewWebApp(&core.Settings{})
		ctrl.Init(web)

		ctrl.ModuleService = &services.ModuleService{}
		return web, ctrl
	}

	t.Run("GET /modules", func(t *testing.T) {
		web, ctrl := setup()

		registry := new(mocked.MockRegistry)
		defer registry.AssertExpectations(t)
		moduleA := new(mocked.MockAppModule)
		defer moduleA.AssertExpectations(t)
		moduleA.On("Avatar").Once().Return("Avatar A")
		moduleA.On("Name").Once().Return("Name A")
		moduleA.On("Desc").Once().Return("Desc A")

		moduleB := new(mocked.MockAppModule)
		defer moduleB.AssertExpectations(t)
		moduleB.On("Avatar").Once().Return("Avatar B")
		moduleB.On("Name").Once().Return("Name B")
		moduleB.On("Desc").Once().Return("Desc B")

		modules := []core.AppModule{
			moduleA, moduleB,
		}
		registry.On("GetModules").Once().Return(modules)
		ctrl.ModuleService.Registry = registry

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/modules", nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result models.ModuleSearchResult
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, len(modules), result.Total)
		assert.Equal(t, len(modules), len(result.Items))

	})

	t.Run("GET /modules with special keyword", func(t *testing.T) {
		web, ctrl := setup()

		registry := new(mocked.MockRegistry)
		defer registry.AssertExpectations(t)
		moduleA := new(mocked.MockAppModule)
		defer moduleA.AssertExpectations(t)
		moduleA.On("Avatar").Once().Return("Avatar A")
		moduleA.On("Name").Twice().Return("Name A")
		moduleA.On("Desc").Twice().Return("Desc A")

		moduleB := new(mocked.MockAppModule)
		defer moduleB.AssertExpectations(t)
		moduleB.On("Avatar").Once().Return("Avatar B")
		moduleB.On("Name").Twice().Return("Name B")
		moduleB.On("Desc").Twice().Return("Desc B")

		moduleC := new(mocked.MockAppModule)
		defer moduleC.AssertExpectations(t)
		moduleC.On("Name").Once().Return("Name C")
		moduleC.On("Desc").Once().Return("Desc C")

		modules := []core.AppModule{
			moduleA, moduleB, moduleC,
		}
		registry.On("GetModules").Once().Return(modules)
		ctrl.ModuleService.Registry = registry

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/modules", nil)
		q := req.URL.Query()
		q.Add("keyword", "Name A,Desc B")
		req.URL.RawQuery = q.Encode()
		web.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result models.ModuleSearchResult
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, 2, len(result.Items))
		itemA := result.Items[0]
		assert.Equal(t, "Avatar A", itemA.Avatar)
		assert.Equal(t, "Name A", itemA.Name)
		assert.Equal(t, "Desc A", itemA.Desc)
		assert.Equal(t, true, itemA.Enabled)
		assert.Equal(t, true, itemA.ReadOnly)
		assert.Equal(t, false, itemA.HasWeb)
		itemB := result.Items[1]
		assert.Equal(t, "Avatar B", itemB.Avatar)
		assert.Equal(t, "Name B", itemB.Name)
		assert.Equal(t, "Desc B", itemB.Desc)
		assert.Equal(t, true, itemB.Enabled)
		assert.Equal(t, true, itemB.ReadOnly)
		assert.Equal(t, false, itemB.HasWeb)

	})

	t.Run("GET /modules/:name", func(t *testing.T) {
		web, ctrl := setup()

		name := "name1"
		registry := new(mocked.MockRegistry)
		defer registry.AssertExpectations(t)
		ctrl.ModuleService.Registry = registry

		module := new(mocked.MockAppModule)
		defer module.AssertExpectations(t)
		module.On("Avatar").Once().Return("Avatar A")
		module.On("Name").Once().Return(name)
		module.On("Desc").Once().Return("Desc A")
		registry.On("GetModuleByName", name).Once().Return(module)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/modules/"+name, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result models.Module
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, "Avatar A", result.Avatar)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, "Desc A", result.Desc)

	})

	t.Run("Get /modules/:name with not found", func(t *testing.T) {
		web, ctrl := setup()

		name := "name1"
		registry := new(mocked.MockRegistry)
		defer registry.AssertExpectations(t)
		ctrl.ModuleService.Registry = registry
		registry.On("GetModuleByName", name).Once().Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/modules/"+name, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

	})

	t.Run("PATCH /modules/:name", func(t *testing.T) {
		web, ctrl := setup()

		name := "name1"
		enabled := true
		avatar := "Avatar A"
		desc := "Desc A"
		moduleFields := &models.ModuleFields{Enabled: &enabled}
		reqBody, err := json.Marshal(moduleFields)
		if err != nil {
			t.Fatal(err)
		}

		registry := new(mocked.MockRegistry)
		defer registry.AssertExpectations(t)
		ctrl.ModuleService.Registry = registry

		module := new(appTestMocked.MockAppEnabledModule)
		defer module.AssertExpectations(t)
		module.On("Desc").Once().Return(desc)
		module.On("Avatar").Once().Return(avatar)
		module.On("Name").Once().Return(name)
		module.On("Enabled").Once().Return(enabled)
		module.On("SetEnable", enabled).Once().Return(nil)
		registry.On("GetModuleByName", name).Once().Return(module)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/modules/"+name, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var result models.Module
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, enabled, result.Enabled)
		assert.Equal(t, avatar, result.Avatar)
		assert.Equal(t, desc, result.Desc)

	})

	t.Run("PATCH /modules/:name with Module Not Found", func(t *testing.T) {
		web, ctrl := setup()

		name := "name2"
		enabled := false
		moduleFields := &models.ModuleFields{Enabled: &enabled}
		reqBody, err := json.Marshal(moduleFields)
		if err != nil {
			t.Fatal(err)
		}

		registry := new(mocked.MockRegistry)
		defer registry.AssertExpectations(t)
		ctrl.ModuleService.Registry = registry

		registry.On("GetModuleByName", name).Once().Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/modules/"+name, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

	})

	t.Run("PATCH /modules/:name with Forbidden", func(t *testing.T) {
		web, ctrl := setup()

		name := "name3"
		enabled := false
		moduleFields := &models.ModuleFields{Enabled: &enabled}
		reqBody, err := json.Marshal(moduleFields)
		if err != nil {
			t.Fatal(err)
		}

		registry := new(mocked.MockRegistry)
		defer registry.AssertExpectations(t)
		ctrl.ModuleService.Registry = registry

		module := new(mocked.MockAppModule)
		defer module.AssertExpectations(t)
		registry.On("GetModuleByName", name).Once().Return(module)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/modules/"+name, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("PATCH /modules/:name with Server Error", func(t *testing.T) {
		web, ctrl := setup()

		name := "name4"
		enabled := false
		moduleFields := &models.ModuleFields{Enabled: &enabled}
		reqBody, err := json.Marshal(moduleFields)
		if err != nil {
			t.Fatal(err)
		}

		registry := new(mocked.MockRegistry)
		defer registry.AssertExpectations(t)
		ctrl.ModuleService.Registry = registry

		module := new(appTestMocked.MockAppEnabledModule)
		defer module.AssertExpectations(t)
		module.On("SetEnable", enabled).Once().Return(errors.New("error"))
		registry.On("GetModuleByName", name).Once().Return(module)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/modules/"+name, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("PATCH /modules/:name with Bad Request", func(t *testing.T) {
		web, ctrl := setup()
		name := "name4"

		registry := new(mocked.MockRegistry)
		defer registry.AssertExpectations(t)
		ctrl.ModuleService.Registry = registry

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/modules/"+name, bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

	})
}
