package gomobile

import (
	"pan/lib/config"
	"pan/lib/repository"
)

type Settings struct {
	HostName   string
	ConfigPath string
	DBPath     string
	TempDBPath string
}

func initWithSettings(settings Settings) {
	config.SetHostName(settings.HostName)
	config.SetRootPath(settings.ConfigPath)

	repository.SetDBPath(settings.DBPath)
	repository.SetTempDBPath(settings.TempDBPath)
}
