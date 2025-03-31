package settings_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pan/lib/config"
	"pan/lib/web"
	"testing"

	appsettings "pan/features/app/settings"

	mocked "pan/mocks/pan/features/app/settings"

	"github.com/stretchr/testify/assert"
)

func TestAppSettingsController(t *testing.T) {

	setup := func() (web.WebApp, *appsettings.AppSettingsController) {
		ctrl := &appsettings.AppSettingsController{}
		webApp := web.NewWebApp()
		ctrl.SetupToWeb(webApp)

		ctrl.AppSettingsService = &appsettings.AppSettingsService{}
		return webApp, ctrl
	}

	t.Run("GET /settings", func(t *testing.T) {
		webApp, ctrl := setup()

		peerId := "test peer id"
		rootPath := "test root path"
		settings := config.Settings{}
		settings.Name = "test name"
		settings.WebAddress = []string{"127.0.0.1:9002"}
		settings.PeerAddress = []string{"127.0.0.1:9001"}
		settings.BroadcastAddress = []string{"127.0.0.1:9000"}
		settings.PublicAddress = []string{"127.0.0.1:9003"}
		settings.GuardEnabled = true
		settings.GuardAccess = true

		provider := &mocked.MockAppSettingsProvider{}
		provider.AssertExpectations(t)
		ctrl.AppSettingsService.Provider = provider
		provider.On("Settings").Once().Return(settings)
		provider.On("PeerID").Once().Return(peerId)
		provider.On("RootPath").Once().Return(rootPath)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/settings", nil)
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results appsettings.AppSettings
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, appsettings.AppSettings{Settings: settings, PeerID: peerId, RootPath: rootPath}, results)

	})

	t.Run("PATCH /settings", func(t *testing.T) {
		webApp, ctrl := setup()

		peerId := "test peer id"
		rootPath := "test root path"
		settings := config.Settings{}
		settings.Name = "test name"
		settings.WebAddress = []string{"127.0.0.1:9002"}
		settings.PeerAddress = []string{"127.0.0.1:9001"}
		settings.BroadcastAddress = []string{"127.0.0.1:9000"}
		settings.PublicAddress = []string{"127.0.0.1:9003"}
		settings.GuardEnabled = true
		settings.GuardAccess = true

		fields := appsettings.AppSettingsFields{}
		fields.Name = "field name"
		fields.WebAddress = []string{"0.0.0.0:9002"}
		fields.PeerAddress = []string{"0.0.0.0:9001"}
		fields.BroadcastAddress = []string{"0.0.0.0:9000"}
		fields.PublicAddress = []string{"0.0.0.0:9003"}
		fields.GuardEnabled = new(bool)
		*fields.GuardEnabled = false
		fields.GuardAccess = new(bool)
		*fields.GuardAccess = false

		settings_ := settings
		settings_.Name = fields.Name
		settings_.WebAddress = fields.WebAddress
		settings_.PeerAddress = fields.PeerAddress
		settings_.BroadcastAddress = fields.BroadcastAddress
		settings_.PublicAddress = fields.PublicAddress
		settings_.GuardEnabled = *fields.GuardEnabled
		settings_.GuardAccess = *fields.GuardAccess

		provider := &mocked.MockAppSettingsProvider{}
		provider.AssertExpectations(t)
		ctrl.AppSettingsService.Provider = provider
		provider.On("Settings").Once().Return(settings)
		provider.On("SetSettings", settings_).Once().Return(nil)
		provider.On("Settings").Once().Return(settings_)
		provider.On("PeerID").Once().Return(peerId)
		provider.On("RootPath").Once().Return(rootPath)

		fieldsData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/settings", bytes.NewReader(fieldsData))
		req.Header.Set("Content-Type", "application/json")
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results appsettings.AppSettings
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, appsettings.AppSettings{Settings: settings_, PeerID: peerId, RootPath: rootPath}, results)
	})
}
