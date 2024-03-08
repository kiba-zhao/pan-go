package services

import (
	"pan/core"
	"pan/extfs/models"
	"sync"
)

type SettingsService struct {
	Settings *models.Settings
	Config   core.Config
	rw       sync.RWMutex
}

func (s *SettingsService) GetTotalHeaderName() string {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.Settings.TotalHeaderName
}
