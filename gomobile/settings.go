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
	logger.Debug("gomobile", "settings.HostName:"+settings.HostName)
	logger.Debug("gomobile", "settings.ConfigPath:"+settings.ConfigPath)
	logger.Debug("gomobile", "settings.DBPath:"+settings.DBPath)
	logger.Debug("gomobile", "settings.TempDBPath:"+settings.TempDBPath)

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
