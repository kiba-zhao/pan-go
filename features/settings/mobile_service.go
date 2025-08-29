//go:build android || ios

package settings

import (
	"errors"
	"sync/atomic"

	"github.com/spf13/viper"
)

var ErrMobileSettingsServiceUnavailable = errors.New("settings.MobileSettingsService Error: Unavailable")

const (
	MobileSettingsWifiOnlyField = "mobileWifiOnly"
)

type MobileSettingsService struct {
	Viper *viper.Viper

	version atomic.Uint32
}

func (service *MobileSettingsService) Load() (MobileSettings, error) {
	viper := service.Viper
	version := service.version.Load()
	return generateMobileSettings(viper, version)
}

func (service *MobileSettingsService) Save(fields MobileSettingsFields) (MobileSettings, error) {
	viper := service.Viper
	version := service.version.Load()
	settings, err := generateMobileSettings(viper, version)
	if err != nil {
		return settings, err
	}

	changed := false
	if fields.WifiOnly != nil && *fields.WifiOnly != *settings.WifiOnly {
		wifiOnly := *fields.WifiOnly
		viper.Set(MobileSettingsWifiOnlyField, wifiOnly)
		settings.WifiOnly = &wifiOnly
		changed = true
	}

	if !changed {
		return settings, nil
	}

	settings.Version = service.version.Add(1)
	err = viper.WriteConfig()
	return settings, err
}

func generateMobileSettings(viper *viper.Viper, version uint32) (MobileSettings, error) {

	var settings MobileSettings
	if viper == nil {
		return settings, ErrMobileSettingsServiceUnavailable
	}

	if viper.IsSet(MobileSettingsWifiOnlyField) {
		wifiOnly := viper.GetBool(MobileSettingsWifiOnlyField)
		settings.WifiOnly = &wifiOnly
	}

	settings.Version = version
	return settings, nil
}
