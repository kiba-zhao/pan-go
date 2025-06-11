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

	configMocked "pan/mocks/pan/lib/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		cfgPath := config.RootPath()
		settings := config.Settings{}
		settings.Name = "test name"
		settings.WebAddress = []string{"127.0.0.1:9002"}
		settings.PeerAddress = []string{"127.0.0.1:9001"}
		settings.BroadcastAddress = []string{"127.0.0.1:9000"}
		settings.PublicAddress = []string{"127.0.0.1:9003"}
		settings.GuardEnabled = true
		settings.GuardAccess = true

		ctrl.AppSettingsService.SetConfigSettings(&settings)
		ctrl.AppSettingsService.SetPeerID(peerId)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/settings", nil)
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results appsettings.AppSettings
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, appsettings.AppSettings{Settings: settings, PeerID: peerId, ConfigPath: cfgPath}, results)

	})

	t.Run("PATCH /settings", func(t *testing.T) {
		webApp, ctrl := setup()

		peerId := "test peer id"
		cfgPath := config.RootPath()
		settings := config.Settings{}
		settings.Name = "test name"
		settings.WebAddress = []string{"127.0.0.1:9002"}
		settings.PeerAddress = []string{"127.0.0.1:9001"}
		settings.BroadcastAddress = []string{"127.0.0.1:9000"}
		settings.PublicAddress = []string{"127.0.0.1:9003"}
		settings.GuardEnabled = true
		settings.GuardAccess = true

		ctrl.AppSettingsService.SetConfigSettings(&settings)
		ctrl.AppSettingsService.SetPeerID(peerId)

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

		var settings_ config.AppSettings
		appConfig := &configMocked.MockAppConfig[config.AppSettings]{}
		appConfig.AssertExpectations(t)
		ctrl.AppSettingsService.AppConfig = appConfig
		appConfig.On("Save", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			settings_ = args.Get(0).(config.AppSettings)
			assert.Equal(t, fields.Name, settings_.Name)
			assert.Equal(t, fields.WebAddress, settings_.WebAddress)
			assert.Equal(t, fields.PeerAddress, settings_.PeerAddress)
			assert.Equal(t, fields.BroadcastAddress, settings_.BroadcastAddress)
			assert.Equal(t, fields.PublicAddress, settings_.PublicAddress)
			assert.Equal(t, *fields.GuardEnabled, settings_.GuardEnabled)
			assert.Equal(t, *fields.GuardAccess, settings_.GuardAccess)
			ctrl.AppSettingsService.SetConfigSettings(settings_)
		})

		fieldsData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/settings", bytes.NewReader(fieldsData))
		req.Header.Set("Content-Type", "application/json")
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results appsettings.AppSettings
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, appsettings.AppSettings{Settings: *settings_, PeerID: peerId, ConfigPath: cfgPath}, results)
	})
}
