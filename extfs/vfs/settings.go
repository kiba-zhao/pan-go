package vfs

import (
	"os"
	appConfig "pan/app/config"
	"path"
)

type VFSSettings struct {
	MountPath string
	LocalName string
}

func newDefaultsVFSSettings(appSettings appConfig.AppSettings) *VFSSettings {
	settings := &VFSSettings{}
	settings.LocalName = appSettings.Name
	settings.MountPath = path.Join(os.TempDir(), "extfs")

	return settings
}
