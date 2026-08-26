//go:build !(android || ios)

package settings

import (
	"errors"

	"github.com/spf13/viper"
)

var ErrWebHostServiceUnavailable = errors.New("settings.WebHostService Error: Unavailable")

const (
	WebHostDefaultWebPort = uint16(9000)
)

const (
	WebHostWebPortField       = "webPort"
	WebHostLocalHostOnlyField = "localHostOnly"
	WebHostWebEnabledField    = "webEnabled"
)

type WebHostChangedTrigger interface {
	OnWebHostChanged(settings WebHost)
}

type WebHostService struct {
	Viper   *viper.Viper
	Trigger WebHostChangedTrigger
}

func (service *WebHostService) Load() (WebHost, error) {
	viper := service.Viper
	return generateWebHost(viper)
}

func (service *WebHostService) Save(fields WebHostFields) (WebHost, error) {
	viper := service.Viper
	settings, err := generateWebHost(viper)
	if err != nil {
		return settings, err
	}

	changed := false
	if fields.WebPort > 0 && fields.WebPort != settings.WebPort {
		viper.Set(WebHostWebPortField, fields.WebPort)
		settings.WebPort = fields.WebPort
		changed = true
	}

	if fields.LocalHostOnly != nil && *fields.LocalHostOnly != settings.LocalHostOnly {
		localHostOnly := *fields.LocalHostOnly
		viper.Set(WebHostLocalHostOnlyField, localHostOnly)
		settings.LocalHostOnly = localHostOnly
		changed = true
	}

	if fields.WebEnabled != nil && *fields.WebEnabled != settings.WebEnabled {
		webEnabled := *fields.WebEnabled
		viper.Set(WebHostWebEnabledField, webEnabled)
		settings.WebEnabled = webEnabled
		changed = true
	}

	if !changed {
		return settings, nil
	}

	err = viper.WriteConfig()
	if service.Trigger != nil {
		service.Trigger.OnWebHostChanged(settings)
	}
	return settings, err
}

func generateWebHost(viper *viper.Viper) (WebHost, error) {

	var settings WebHost
	if viper == nil {
		return settings, ErrWebHostServiceUnavailable
	}

	if viper.IsSet(WebHostWebPortField) {
		settings.WebPort = viper.GetUint16(WebHostWebPortField)
	} else {
		settings.WebPort = WebHostDefaultWebPort
	}

	if viper.IsSet(WebHostLocalHostOnlyField) {
		settings.LocalHostOnly = viper.GetBool(WebHostLocalHostOnlyField)
	} else {
		settings.LocalHostOnly = true
	}

	if viper.IsSet(WebHostWebEnabledField) {
		settings.WebEnabled = viper.GetBool(WebHostWebEnabledField)
	} else {
		settings.WebEnabled = true
	}
	return settings, nil
}
