//go:build android || ios

package settings

import (
	"errors"

	"github.com/spf13/viper"
)

var ErrMobileSettingsServiceUnavailable = errors.New("settings.MobileSettingsService Error: Unavailable")

const (
	MobileSettingsWifiOnlyField = "mobileWifiOnly"
)

type MobileSettingsChangedTrigger interface {
	OnMobileSettingsChanged(settings MobileSettings)
}

type MobileSettingsService struct {
	Viper   *viper.Viper
	Trigger MobileSettingsChangedTrigger
}

func (service *MobileSettingsService) Load() (MobileSettings, error) {
	viper := service.Viper
	return generateMobileSettings(viper)
}

func (service *MobileSettingsService) Save(fields MobileSettingsFields) (MobileSettings, error) {
	viper := service.Viper
	settings, err := generateMobileSettings(viper)
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

	err = viper.WriteConfig()
	if service.Trigger != nil {
		service.Trigger.OnMobileSettingsChanged(settings)
	}

	return settings, err
}

func generateMobileSettings(viper *viper.Viper) (MobileSettings, error) {

	var settings MobileSettings
	if viper == nil {
		return settings, ErrMobileSettingsServiceUnavailable
	}

	if viper.IsSet(MobileSettingsWifiOnlyField) {
		wifiOnly := viper.GetBool(MobileSettingsWifiOnlyField)
		settings.WifiOnly = &wifiOnly
	}

	return settings, nil
}
