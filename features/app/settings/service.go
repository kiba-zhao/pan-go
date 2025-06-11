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
	if fields.WebAddress != nil {
		settings.WebAddress = fields.WebAddress
	}
	if fields.PeerAddress != nil {
		settings.PeerAddress = fields.PeerAddress
	}
	if fields.BroadcastAddress != nil {
		settings.BroadcastAddress = fields.BroadcastAddress
	}
	if fields.PublicAddress != nil {
		settings.PublicAddress = fields.PublicAddress
	}

	if fields.GuardEnabled != nil {
		settings.GuardEnabled = *fields.GuardEnabled
	}
	if fields.GuardAccess != nil {
		settings.GuardAccess = *fields.GuardAccess
	}

	err := s.AppConfig.Save(settings)
	if err != nil {
		return AppSettings{}, err
	}

	return s.Load(), nil
}
