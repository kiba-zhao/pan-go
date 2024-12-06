package vfs

import (
	"os"
	appConfig "pan/app/config"
	"path"
)

type VFSSettings struct {
	MountPath string
	Enabled   bool
	LocalName string
}

func newDefaultsVFSSettings(appSettings appConfig.AppSettings) *VFSSettings {
	settings := &VFSSettings{}
	settings.MountPath = path.Join(os.TempDir(), "extfs")
	settings.Enabled = true
	settings.LocalName = appSettings.Name

	return settings
}
