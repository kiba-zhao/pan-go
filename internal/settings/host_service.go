//go:build !(android || ios)

package settings

import (
	"errors"

	"github.com/spf13/viper"
)

var ErrHostSettingsServiceUnavailable = errors.New("settings.HostSettingsService Error: Unavailable")

const (
	HostSettingsDefaultWebPort = uint16(9000)
)

const (
	HostSettingsWebPortField       = "webPort"
	HostSettingsLocalHostOnlyField = "localHostOnly"
	HostSettingsWebEnabledField    = "webEnabled"
)

type HostSettingsChangedTrigger interface {
	OnHostSettingsChanged(settings HostSettings)
}

type HostSettingsService struct {
	Viper   *viper.Viper
	Trigger HostSettingsChangedTrigger
}

func (service *HostSettingsService) Load() (HostSettings, error) {
	viper := service.Viper
	return generateHostSettings(viper)
}

func (service *HostSettingsService) Save(fields HostSettingsFields) (HostSettings, error) {
	viper := service.Viper
	settings, err := generateHostSettings(viper)
	if err != nil {
		return settings, err
	}

	changed := false
	if fields.WebPort > 0 && fields.WebPort != settings.WebPort {
		viper.Set(HostSettingsWebPortField, fields.WebPort)
		settings.WebPort = fields.WebPort
		changed = true
	}

	if fields.LocalHostOnly != nil && *fields.LocalHostOnly != settings.LocalHostOnly {
		localHostOnly := *fields.LocalHostOnly
		viper.Set(HostSettingsLocalHostOnlyField, localHostOnly)
		settings.LocalHostOnly = localHostOnly
		changed = true
	}

	if fields.WebEnabled != nil && *fields.WebEnabled != settings.WebEnabled {
		webEnabled := *fields.WebEnabled
		viper.Set(HostSettingsWebEnabledField, webEnabled)
		settings.WebEnabled = webEnabled
		changed = true
	}

	if !changed {
		return settings, nil
	}

	err = viper.WriteConfig()
	if service.Trigger != nil {
		service.Trigger.OnHostSettingsChanged(settings)
	}
	return settings, err
}

func generateHostSettings(viper *viper.Viper) (HostSettings, error) {

	var settings HostSettings
	if viper == nil {
		return settings, ErrHostSettingsServiceUnavailable
	}

	if viper.IsSet(HostSettingsWebPortField) {
		settings.WebPort = viper.GetUint16(HostSettingsWebPortField)
	} else {
		settings.WebPort = HostSettingsDefaultWebPort
	}

	if viper.IsSet(HostSettingsLocalHostOnlyField) {
		settings.LocalHostOnly = viper.GetBool(HostSettingsLocalHostOnlyField)
	} else {
		settings.LocalHostOnly = true
	}

	if viper.IsSet(HostSettingsWebEnabledField) {
		settings.WebEnabled = viper.GetBool(HostSettingsWebEnabledField)
	} else {
		settings.WebEnabled = true
	}
	return settings, nil
}
