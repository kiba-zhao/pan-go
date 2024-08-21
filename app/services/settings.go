package services

import (
	"pan/app/config"
	"pan/app/models"
)

type SettingsProvider interface {
	NodeID() string
	Settings() config.Settings
	SetSettings(config.Settings) error
}

type SettingsService struct {
	Provider SettingsProvider
}

func (s *SettingsService) Load() models.Settings {
	settings := models.Settings{}
	settings.Settings = s.Provider.Settings()
	settings.NodeID = s.Provider.NodeID()
	return settings
}

func (s *SettingsService) Save(fields models.SettingsFields) (models.Settings, error) {
	settings := s.Provider.Settings()
	if fields.RootPath != "" {
		settings.RootPath = fields.RootPath
	}
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

	err := s.Provider.SetSettings(settings)
	if err != nil {
		return models.Settings{}, err
	}

	return s.Load(), nil
}
