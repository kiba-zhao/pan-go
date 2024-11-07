package vfs

import (
	"os"
	appConfig "pan/app/config"
	"path"
)

type VFSSettings struct {
	MountPath string
	LocalName string
	Enabled   bool
}

func newDefaultsVFSSettings(appSettings appConfig.AppSettings) *VFSSettings {
	settings := &VFSSettings{}
	settings.LocalName = appSettings.Name
	settings.MountPath = path.Join(os.TempDir(), "extfs")
	settings.Enabled = true

	return settings
}
