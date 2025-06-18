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

func initWithSettings(settings *Settings, logger Logger) {
	var err error
	err = config.InitHostName(settings.HostName)
	if err != nil {
		logger.Error("gomobile", "settings Error:"+err.Error())
	}
	err = config.InitRootPath(settings.ConfigPath)
	if err != nil {
		logger.Error("gomobile", "settings Error:"+err.Error())
	}

	err = repository.InitDBPath(settings.DBPath)
	if err != nil {
		logger.Error("gomobile", "settings Error:"+err.Error())
	}

	err = repository.InitTempDBPath(settings.TempDBPath)
	if err != nil {
		logger.Error("gomobile", "settings Error:"+err.Error())
	}

}
