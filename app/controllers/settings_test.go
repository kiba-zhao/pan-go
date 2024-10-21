package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pan/app/config"
	"pan/app/controllers"
	"pan/app/models"
	"pan/app/net"
	"pan/app/services"
	"testing"

	mocked "pan/mocks/pan/app/services"

	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {

	setup := func() (net.WebApp, *controllers.SettingsController) {
		ctrl := &controllers.SettingsController{}
		web := net.NewWebApp()
		ctrl.SetupToWeb(web)

		ctrl.SettingsService = &services.SettingsService{}
		return web, ctrl
	}

	t.Run("GET /settings", func(t *testing.T) {
		web, ctrl := setup()

		nodeId := "test node id"
		rootPath := "test root path"
		settings := config.Settings{}
		settings.Name = "test name"
		settings.WebAddress = []string{"127.0.0.1:9002"}
		settings.NodeAddress = []string{"127.0.0.1:9001"}
		settings.BroadcastAddress = []string{"127.0.0.1:9000"}
		settings.PublicAddress = []string{"127.0.0.1:9003"}
		settings.GuardEnabled = true
		settings.GuardAccess = true

		provider := &mocked.MockSettingsProvider{}
		provider.AssertExpectations(t)
		ctrl.SettingsService.Provider = provider
		provider.On("Settings").Once().Return(settings)
		provider.On("NodeID").Once().Return(nodeId)
		provider.On("RootPath").Once().Return(rootPath)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/settings", nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results models.Settings
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, models.Settings{Settings: settings, NodeID: nodeId, RootPath: rootPath}, results)

	})

	t.Run("PATCH /settings", func(t *testing.T) {
		web, ctrl := setup()

		nodeId := "test node id"
		rootPath := "test root path"
		settings := config.Settings{}
		settings.Name = "test name"
		settings.WebAddress = []string{"127.0.0.1:9002"}
		settings.NodeAddress = []string{"127.0.0.1:9001"}
		settings.BroadcastAddress = []string{"127.0.0.1:9000"}
		settings.PublicAddress = []string{"127.0.0.1:9003"}
		settings.GuardEnabled = true
		settings.GuardAccess = true

		fields := models.SettingsFields{}
		fields.Name = "field name"
		fields.WebAddress = []string{"0.0.0.0:9002"}
		fields.NodeAddress = []string{"0.0.0.0:9001"}
		fields.BroadcastAddress = []string{"0.0.0.0:9000"}
		fields.PublicAddress = []string{"0.0.0.0:9003"}
		fields.GuardEnabled = new(bool)
		*fields.GuardEnabled = false
		fields.GuardAccess = new(bool)
		*fields.GuardAccess = false

		settings_ := settings
		settings_.Name = fields.Name
		settings_.WebAddress = fields.WebAddress
		settings_.NodeAddress = fields.NodeAddress
		settings_.BroadcastAddress = fields.BroadcastAddress
		settings_.PublicAddress = fields.PublicAddress
		settings_.GuardEnabled = *fields.GuardEnabled
		settings_.GuardAccess = *fields.GuardAccess

		provider := &mocked.MockSettingsProvider{}
		provider.AssertExpectations(t)
		ctrl.SettingsService.Provider = provider
		provider.On("Settings").Once().Return(settings)
		provider.On("SetSettings", settings_).Once().Return(nil)
		provider.On("Settings").Once().Return(settings_)
		provider.On("NodeID").Once().Return(nodeId)
		provider.On("RootPath").Once().Return(rootPath)

		fieldsData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/settings", bytes.NewReader(fieldsData))
		req.Header.Set("Content-Type", "application/json")
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results models.Settings
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, models.Settings{Settings: settings_, NodeID: nodeId, RootPath: rootPath}, results)
	})
}
