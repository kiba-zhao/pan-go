//go:build android || ios

package settings

import (
	"errors"

	"github.com/spf13/viper"
)

var ErrMobileNetworkServiceUnavailable = errors.New("settings.MobileNetworkService Error: Unavailable")

const (
	MobileNetworkWifiOnlyField = "mobileWifiOnly"
)

type MobileNetworkChangedTrigger interface {
	OnMobileNetworkChanged(settings MobileNetwork)
}

type MobileNetworkService struct {
	Viper   *viper.Viper
	Trigger MobileNetworkChangedTrigger
}

func (service *MobileNetworkService) Load() (MobileNetwork, error) {
	viper := service.Viper
	return generateMobileNetwork(viper)
}

func (service *MobileNetworkService) Save(fields MobileNetworkFields) (MobileNetwork, error) {
	viper := service.Viper
	settings, err := generateMobileNetwork(viper)
	if err != nil {
		return settings, err
	}

	changed := false
	if fields.WifiOnly != nil && *fields.WifiOnly != *settings.WifiOnly {
		wifiOnly := *fields.WifiOnly
		viper.Set(MobileNetworkWifiOnlyField, wifiOnly)
		settings.WifiOnly = &wifiOnly
		changed = true
	}

	if !changed {
		return settings, nil
	}

	err = viper.WriteConfig()
	if service.Trigger != nil {
		service.Trigger.OnMobileNetworkChanged(settings)
	}

	return settings, err
}

func generateMobileNetwork(viper *viper.Viper) (MobileNetwork, error) {

	var settings MobileNetwork
	if viper == nil {
		return settings, ErrMobileNetworkServiceUnavailable
	}

	if viper.IsSet(MobileNetworkWifiOnlyField) {
		wifiOnly := viper.GetBool(MobileNetworkWifiOnlyField)
		settings.WifiOnly = &wifiOnly
	}

	return settings, nil
}
