//go:build !(android || ios) || host

package settings

import (
	"errors"
	"net"
	"strconv"

	"github.com/spf13/viper"
)

var ErrHostSettingsServiceUnavailable = errors.New("settings.HostSettingsService Error: Unavailable")

const (
	HostSettingsDefaultWebIP = "127.0.0.1"
)

const (
	HostSettingsWebAddrField    = "webAddr"
	HostSettingsWebEnabledField = "webEnabled"
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
	if len(fields.WebAddr) > 0 && fields.WebAddr != settings.WebAddr {
		viper.Set(HostSettingsWebAddrField, fields.WebAddr)
		settings.WebAddr = fields.WebAddr
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

	if viper.IsSet(HostSettingsWebAddrField) {
		settings.WebAddr = viper.GetString(HostSettingsWebAddrField)
	} else {
		settings.WebAddr = net.JoinHostPort(HostSettingsDefaultWebIP, strconv.FormatUint(uint64(SettingsDefaultPort), 10))
	}

	if viper.IsSet(HostSettingsWebEnabledField) {
		settings.WebEnabled = viper.GetBool(HostSettingsWebEnabledField)
	} else {
		settings.WebEnabled = true
	}
	return settings, nil
}
