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
