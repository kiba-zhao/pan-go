// Define vfs settings
package vfs

import (
	appConfig "pan/app/config"
)

var MountPath string

type VFSSettings struct {
	MountPath string
	Enabled   bool
	LocalName string
}

func newDefaultsVFSSettings(appSettings appConfig.AppSettings) *VFSSettings {
	settings := &VFSSettings{}
	settings.MountPath = MountPath
	settings.Enabled = true
	settings.LocalName = appSettings.Name

	return settings
}
