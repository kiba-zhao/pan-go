package vfs

import (
	"os"
	appConfig "pan/app/config"
	"path"
)

type VFSSettings struct {
	MountPath     string
	Enabled       bool
	LocalDirName  string
	RemoteDirName string
}

func newDefaultsVFSSettings(appSettings appConfig.AppSettings) *VFSSettings {
	settings := &VFSSettings{}
	settings.MountPath = path.Join(os.TempDir(), "extfs")
	settings.Enabled = true
	settings.LocalDirName = "Local"
	settings.RemoteDirName = "Remote"

	return settings
}
