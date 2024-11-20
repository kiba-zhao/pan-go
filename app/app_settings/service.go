package appsettings

import (
	"pan/app/config"
)

type AppSettingsProvider interface {
	RootPath() string
	PeerID() string
	Settings() config.Settings
	SetSettings(config.Settings) error
}

type AppSettingsExternalService interface {
	Load() AppSettings
}

type AppSettingsService struct {
	Provider AppSettingsProvider
}

func (s *AppSettingsService) Load() AppSettings {
	settings := AppSettings{}
	settings.Settings = s.Provider.Settings()
	settings.PeerID = s.Provider.PeerID()
	settings.RootPath = s.Provider.RootPath()
	return settings
}

func (s *AppSettingsService) Save(fields AppSettingsFields) (AppSettings, error) {
	settings := s.Provider.Settings()
	if fields.Name != "" {
		settings.Name = fields.Name
	}
	if fields.WebAddress != nil {
		settings.WebAddress = fields.WebAddress
	}
	if fields.NodeAddress != nil {
		settings.NodeAddress = fields.NodeAddress
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

	err := s.Provider.SetSettings(settings)
	if err != nil {
		return AppSettings{}, err
	}

	return s.Load(), nil
}
