package settings

import (
	"pan/lib/config"
	"sync"
)

type AppSettingsExternalService interface {
	Load() AppSettings
}

type AppSettingsService struct {
	AppConfig config.AppConfig

	settings   AppSettings
	settingsRW sync.RWMutex
}

func (service *AppSettingsService) Load() AppSettings {
	service.settingsRW.RLock()
	defer service.settingsRW.RUnlock()

	settings := service.settings
	settings.ConfigPath = config.RootPath()

	return settings
}

func (service *AppSettingsService) SetConfigSettings(settings config.AppSettings) {
	service.settingsRW.Lock()
	defer service.settingsRW.Unlock()
	service.settings.Settings = *settings
}

func (service *AppSettingsService) SetPeerID(peerId string) {
	service.settingsRW.Lock()
	defer service.settingsRW.Unlock()
	service.settings.PeerID = peerId
}

func (s *AppSettingsService) Save(fields AppSettingsFields) (AppSettings, error) {
	appSettings := s.Load()
	settings := &appSettings.Settings
	if fields.Name != "" {
		settings.Name = fields.Name
	}
	if fields.WebAddr != nil {
		settings.WebAddr = *fields.WebAddr
	}
	if fields.PeerPort != nil {
		settings.PeerPort = *fields.PeerPort
	}
	if fields.BroadcastAddrs != nil {
		settings.BroadcastAddrs = fields.BroadcastAddrs
	}
	if fields.PublicAddrs != nil {
		settings.PublicAddrs = fields.PublicAddrs
	}

	if fields.Enabled != nil {
		settings.Enabled = *fields.Enabled
	}

	err := s.AppConfig.Save(settings)
	if err != nil {
		return AppSettings{}, err
	}

	return s.Load(), nil
}
