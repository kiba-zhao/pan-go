//go:build !(android || ios) || host

package settings

import (
	"errors"
	"net"
	"strconv"
	"sync/atomic"

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

type HostSettingsService struct {
	Viper *viper.Viper

	version atomic.Uint32
}

func (service *HostSettingsService) Load() (HostSettings, error) {
	viper := service.Viper
	version := service.version.Load()
	return generateHostSettings(viper, version)
}

func (service *HostSettingsService) Save(fields HostSettingsFields) (HostSettings, error) {
	viper := service.Viper
	version := service.version.Load()
	settings, err := generateHostSettings(viper, version)
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

	settings.Version = service.version.Add(1)
	err = viper.WriteConfig()
	return settings, err
}

func generateHostSettings(viper *viper.Viper, version uint32) (HostSettings, error) {

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
	settings.Version = version
	return settings, nil
}
