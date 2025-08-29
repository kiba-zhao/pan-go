package settings

import (
	"errors"
	"net"
	"slices"
	"strconv"
	"sync"

	"github.com/spf13/viper"
)

var ErrSettingsServiceUnavailable = errors.New("settings.SettingsService Error: Unavailable")

const (
	SettingsDefaultPort          = 9000
	SettingsDefaultBroadcastIP   = "224.0.0.2"
	SettingsDefaultBroadcastPort = SettingsDefaultPort + 1
)

const (
	SettingsNameField             = "name"
	SettingsPeerPortField         = "peerPort"
	SettingsBroadcastAddrsField   = "broadcastAddrs"
	SettingsPublicAddrsField      = "publicAddrs"
	SettingsEnabledField          = "enabled"
	SettingsBroadcastEnabledField = "broadcastEnabled"
)

type SettingsChangedTrigger interface {
	OnSettingsChanged(settings Settings)
}

type SettingsService struct {
	Viper   *viper.Viper
	Trigger SettingsChangedTrigger

	cfg   SettingsConfig
	cfgRW sync.RWMutex
}

func (service *SettingsService) Setup(cfg SettingsConfig) {
	service.cfgRW.Lock()
	defer service.cfgRW.Unlock()
	service.cfg = cfg
}

func (service *SettingsService) Load() (Settings, error) {

	service.cfgRW.RLock()
	cfg := service.cfg
	defer service.cfgRW.RUnlock()

	viper := service.Viper
	return generateSettings(viper, cfg)
}

func (service *SettingsService) Save(fields SettingsFields) (Settings, error) {
	service.cfgRW.RLock()
	cfg := service.cfg
	defer service.cfgRW.RUnlock()

	viper := service.Viper
	settings, err := generateSettings(viper, cfg)
	if err != nil {
		return settings, err
	}

	changed := false

	if len(fields.Name) > 0 && fields.Name != settings.Name {
		viper.Set(SettingsNameField, fields.Name)
		settings.Name = fields.Name
		changed = true
	}

	if fields.PeerPort != nil && *fields.PeerPort != settings.PeerPort {
		peerPort := *fields.PeerPort
		viper.Set(SettingsPeerPortField, peerPort)
		settings.PeerPort = peerPort
		changed = true
	}

	if len(fields.BroadcastAddrs) > 0 && slices.Equal(fields.BroadcastAddrs, settings.BroadcastAddrs) {
		viper.Set(SettingsBroadcastAddrsField, fields.BroadcastAddrs)
		settings.BroadcastAddrs = fields.BroadcastAddrs
		changed = true
	}

	if len(fields.PublicAddrs) > 0 && slices.Equal(fields.PublicAddrs, settings.PublicAddrs) {
		viper.Set(SettingsPublicAddrsField, fields.PublicAddrs)
		settings.PublicAddrs = fields.PublicAddrs
		changed = true
	}

	if fields.Enabled != nil && *fields.Enabled != settings.Enabled {
		enabled := *fields.Enabled
		viper.Set(SettingsEnabledField, enabled)
		settings.Enabled = enabled
		changed = true
	}

	if fields.BroadcastEnabled != nil && *fields.BroadcastEnabled != settings.BroadcastEnabled {
		broadcastEnabled := *fields.BroadcastEnabled
		viper.Set(SettingsBroadcastEnabledField, broadcastEnabled)
		settings.BroadcastEnabled = broadcastEnabled
		changed = true
	}

	if !changed {
		return settings, nil
	}

	err = viper.WriteConfig()
	if service.Trigger != nil {
		service.Trigger.OnSettingsChanged(settings)
	}
	return settings, err
}

func generateSettings(viper *viper.Viper, cfg SettingsConfig) (Settings, error) {

	var settings Settings
	if viper == nil || cfg == nil {
		return settings, ErrSettingsServiceUnavailable
	}

	if viper.IsSet(SettingsNameField) {
		settings.Name = viper.GetString(SettingsNameField)
	} else {
		settings.Name = cfg.HostName()
	}

	if viper.IsSet(SettingsPeerPortField) {
		settings.PeerPort = viper.GetUint16(SettingsPeerPortField)
	} else {
		settings.PeerPort = SettingsDefaultPort
	}

	if viper.IsSet(SettingsBroadcastAddrsField) {
		settings.BroadcastAddrs = viper.GetStringSlice(SettingsBroadcastAddrsField)
	} else {
		broadcastAddr := net.JoinHostPort(SettingsDefaultBroadcastIP, strconv.FormatUint(uint64(SettingsDefaultBroadcastPort), 10))
		settings.BroadcastAddrs = append(settings.BroadcastAddrs, broadcastAddr)
	}

	if viper.IsSet(SettingsPublicAddrsField) {
		settings.PublicAddrs = viper.GetStringSlice(SettingsPublicAddrsField)
	}

	if viper.IsSet(SettingsEnabledField) {
		settings.Enabled = viper.GetBool(SettingsEnabledField)
	} else {
		settings.Enabled = true
	}

	if viper.IsSet(SettingsBroadcastEnabledField) {
		settings.BroadcastEnabled = viper.GetBool(SettingsBroadcastEnabledField)
	} else {
		settings.BroadcastEnabled = true
	}

	return settings, nil
}
