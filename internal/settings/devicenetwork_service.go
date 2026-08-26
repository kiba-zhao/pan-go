package settings

import (
	"errors"
	"net"
	"slices"
	"strconv"

	"github.com/spf13/viper"
)

var ErrDeviceNetworkServiceUnavailable = errors.New("settings.DeviceNetworkService Error: Unavailable")

const (
	DeviceNetworkDefaultPort          = 9000
	DeviceNetworkDefaultBroadcastIP   = "224.0.0.2"
	DeviceNetworkDefaultBroadcastPort = DeviceNetworkDefaultPort + 1
)

const (
	DeviceNetworkEnabledField          = "enabled"
	DeviceNetworkPortField             = "port"
	DeviceNetworkPublicAddrsField      = "publicAddrs"
	DeviceNetworkBroadcastAddrsField   = "broadcastAddrs"
	DeviceNetworkBroadcastEnabledField = "broadcastEnabled"
)

type DeviceNetworkChangedTrigger interface {
	OnDeviceNetworkChanged(deviceNetwork DeviceNetwork)
}

type DeviceNetworkService struct {
	Viper   *viper.Viper
	Trigger DeviceNetworkChangedTrigger
}

func (service *DeviceNetworkService) Load() (DeviceNetwork, error) {
	return generateDeviceNetwork(service.Viper)
}

func (service *DeviceNetworkService) Save(fields DeviceNetworkFields) (DeviceNetwork, error) {

	viper := service.Viper
	settings, err := generateDeviceNetwork(viper)
	if err != nil {
		return settings, err
	}

	changed := false

	if fields.Enabled != nil && *fields.Enabled != settings.Enabled {
		enabled := *fields.Enabled
		viper.Set(DeviceNetworkEnabledField, enabled)
		settings.Enabled = enabled
		changed = true
	}

	if fields.Port != nil && *fields.Port != settings.Port {
		port := *fields.Port
		viper.Set(DeviceNetworkPortField, port)
		settings.Port = port
		changed = true
	}

	if fields.PublicAddrs != nil && !slices.Equal(fields.PublicAddrs, settings.PublicAddrs) {
		viper.Set(DeviceNetworkPublicAddrsField, fields.PublicAddrs)
		settings.PublicAddrs = fields.PublicAddrs
		changed = true
	}

	if fields.BroadcastAddrs != nil && !slices.Equal(fields.BroadcastAddrs, settings.BroadcastAddrs) {
		viper.Set(DeviceNetworkBroadcastAddrsField, fields.BroadcastAddrs)
		settings.BroadcastAddrs = fields.BroadcastAddrs
		changed = true
	}

	if fields.BroadcastEnabled != nil && *fields.BroadcastEnabled != settings.BroadcastEnabled {
		broadcastEnabled := *fields.BroadcastEnabled
		viper.Set(DeviceNetworkBroadcastEnabledField, broadcastEnabled)
		settings.BroadcastEnabled = broadcastEnabled
		changed = true
	}

	if !changed {
		return settings, nil
	}

	err = viper.WriteConfig()
	if service.Trigger != nil {
		service.Trigger.OnDeviceNetworkChanged(settings)
	}
	return settings, err
}

func generateDeviceNetwork(viper *viper.Viper) (DeviceNetwork, error) {

	var settings DeviceNetwork
	if viper == nil {
		return settings, ErrDeviceNetworkServiceUnavailable
	}

	if viper.IsSet(DeviceNetworkEnabledField) {
		settings.Enabled = viper.GetBool(DeviceNetworkEnabledField)
	} else {
		settings.Enabled = true
	}

	if viper.IsSet(DeviceNetworkPortField) {
		settings.Port = viper.GetUint16(DeviceNetworkPortField)
	} else {
		settings.Port = DeviceNetworkDefaultPort
	}

	if viper.IsSet(DeviceNetworkPublicAddrsField) {
		settings.PublicAddrs = viper.GetStringSlice(DeviceNetworkPublicAddrsField)
	}

	if viper.IsSet(DeviceNetworkBroadcastAddrsField) {
		settings.BroadcastAddrs = viper.GetStringSlice(DeviceNetworkBroadcastAddrsField)
	} else {
		broadcastAddr := net.JoinHostPort(DeviceNetworkDefaultBroadcastIP, strconv.FormatUint(uint64(DeviceNetworkDefaultBroadcastPort), 10))
		settings.BroadcastAddrs = append(settings.BroadcastAddrs, broadcastAddr)
	}

	if viper.IsSet(DeviceNetworkBroadcastEnabledField) {
		settings.BroadcastEnabled = viper.GetBool(DeviceNetworkBroadcastEnabledField)
	} else {
		settings.BroadcastEnabled = true
	}

	return settings, nil
}
