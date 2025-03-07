// Define vfs settings
package vfs

import (
	"os"
	appConfig "pan/app/config"
	"path"
	"runtime"
)

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

func generateDefaultMountPath() string {
	var mountPath string
	switch runtime.GOOS {
	case "windows":
		for i := 'd'; i <= 'z'; i++ {
			mountPath = string(i)
			if _, err := os.Stat(mountPath); os.IsNotExist(err) {
				break
			}
		}

	case "linux":
		// mountPath = path.Join("/", "media", "extfs")
		home, err := os.UserHomeDir()
		if err != nil {
			home = path.Join("/", "media")
		}
		mountPath = path.Join(home, "extfs")
	case "darwin":
		mountPath = path.Join("/", "Volumes", "extfs")
	default:
		return ""
	}
	return mountPath
}

var MountPath string

func init() {
	MountPath = generateDefaultMountPath()
}
