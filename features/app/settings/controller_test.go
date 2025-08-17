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
		settings.WebAddr = "127.0.0.1:9002"
		settings.PeerPort = 9001
		settings.BroadcastAddrs = []string{"127.0.0.1:9000"}
		settings.PublicAddrs = []string{"127.0.0.1:9003"}
		settings.Enabled = true

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
		settings.WebAddr = "127.0.0.1:9002"
		settings.PeerPort = 9001
		settings.BroadcastAddrs = []string{"127.0.0.1:9000"}
		settings.PublicAddrs = []string{"127.0.0.1:9003"}
		settings.Enabled = true

		ctrl.AppSettingsService.SetConfigSettings(&settings)
		ctrl.AppSettingsService.SetPeerID(peerId)

		fields := appsettings.AppSettingsFields{}
		fields.Name = "field name"
		webAddr := "0.0.0.0:9002"
		fields.WebAddr = &webAddr
		peerPort := uint16(9007)
		fields.PeerPort = &peerPort
		fields.BroadcastAddrs = []string{"0.0.0.0:9000"}
		fields.PublicAddrs = []string{"0.0.0.0:9003"}
		enabled := false
		fields.Enabled = &enabled

		var settings_ config.AppSettings
		appConfig := &configMocked.MockAppConfig[config.AppSettings]{}
		appConfig.AssertExpectations(t)
		ctrl.AppSettingsService.AppConfig = appConfig
		appConfig.On("Save", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			settings_ = args.Get(0).(config.AppSettings)
			assert.Equal(t, fields.Name, settings_.Name)
			assert.Equal(t, *fields.WebAddr, settings_.WebAddr)
			assert.Equal(t, *fields.PeerPort, settings_.PeerPort)
			assert.Equal(t, fields.BroadcastAddrs, settings_.BroadcastAddrs)
			assert.Equal(t, fields.PublicAddrs, settings_.PublicAddrs)
			assert.Equal(t, *fields.Enabled, settings_.Enabled)
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
