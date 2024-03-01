package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pan/app/controllers"
	"pan/app/models"
	"pan/app/services"
	"pan/core"
	"testing"

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
		moduleA := new(mocked.MockAppModule)
		moduleA.On("Avatar").Once().Return("Avatar A")
		moduleA.On("Name").Once().Return("Name A")
		moduleA.On("Desc").Once().Return("Desc A")

		moduleB := new(mocked.MockAppModule)
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
		moduleA := new(mocked.MockAppModule)
		moduleA.On("Avatar").Twice().Return("Avatar A")
		moduleA.On("Name").Twice().Return("Name A")
		moduleA.On("Desc").Twice().Return("Desc A")

		moduleB := new(mocked.MockAppModule)
		moduleB.On("Avatar").Twice().Return("Avatar B")
		moduleB.On("Name").Twice().Return("Name B")
		moduleB.On("Desc").Twice().Return("Desc B")

		moduleC := new(mocked.MockAppModule)
		moduleC.On("Avatar").Once().Return("Avatar C")
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
}
