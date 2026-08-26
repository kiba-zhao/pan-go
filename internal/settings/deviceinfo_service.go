package settings

import (
	"errors"
	"sync"

	"github.com/spf13/viper"
)

var ErrDeviceInfoServiceUnavailable = errors.New("settings.DeviceInfoService Error: Unavailable")

const (
	DeviceInfoNameField = "name"
	DeviceInfoMemoField = "memo"
)

type DeviceInfoService struct {
	Viper *viper.Viper

	cfg   SettingsConfig
	cfgRW sync.RWMutex
}

func (service *DeviceInfoService) setup(cfg SettingsConfig) {
	service.cfgRW.Lock()
	defer service.cfgRW.Unlock()
	service.cfg = cfg
}

func (service *DeviceInfoService) Load() (DeviceInfo, error) {
	service.cfgRW.RLock()
	cfg := service.cfg
	defer service.cfgRW.RUnlock()

	viper := service.Viper
	return generateDeviceInfo(viper, cfg)
}

func (service *DeviceInfoService) Save(fields DeviceInfoFields) (DeviceInfo, error) {
	service.cfgRW.RLock()
	cfg := service.cfg
	defer service.cfgRW.RUnlock()

	viper := service.Viper
	settings, err := generateDeviceInfo(viper, cfg)
	if err != nil {
		return settings, err
	}

	changed := false

	if len(fields.Name) > 0 && fields.Name != settings.Name {
		viper.Set(DeviceInfoNameField, fields.Name)
		settings.Name = fields.Name
		changed = true
	}

	if fields.Memo != nil && *fields.Memo != settings.Memo {
		viper.Set(DeviceInfoMemoField, fields.Memo)
		settings.Memo = *fields.Memo
		changed = true
	}

	if !changed {
		return settings, nil
	}

	err = viper.WriteConfig()
	return settings, err

}

func generateDeviceInfo(viper *viper.Viper, cfg SettingsConfig) (DeviceInfo, error) {

	var settings DeviceInfo
	if viper == nil || cfg == nil {
		return settings, ErrDeviceInfoServiceUnavailable
	}

	if viper.IsSet(DeviceInfoNameField) {
		settings.Name = viper.GetString(DeviceInfoNameField)
	} else {
		settings.Name = cfg.HostName()
	}

	if viper.IsSet(DeviceInfoMemoField) {
		settings.Memo = viper.GetString(DeviceInfoMemoField)
	} else {
		settings.Memo = ""
	}

	return settings, nil
}
