// Define vfs settings
package vfs

import (
	"errors"
	appConfig "pan/lib/config"
	"sync"
)

type VFSSettings struct {
	MountPath string
	HostName  string
}

func newDefaultsVFSSettings(appSettings appConfig.AppSettings) *VFSSettings {
	settings := &VFSSettings{}
	settings.MountPath = MountPath()
	settings.HostName = appSettings.Name

	return settings
}

var mountPath string
var mountPathRW sync.RWMutex
var mountPathAlready bool
var errMountPathConflict = errors.New("vfs.mountPath Error: Conflict")

func MountPath() string {
	mountPathRW.RLock()
	defer mountPathRW.RUnlock()

	return mountPath
}

func InitMountPath(path string) error {
	mountPathRW.Lock()
	defer mountPathRW.Unlock()
	if mountPathAlready {
		return errMountPathConflict
	}
	mountPath = path
	mountPathAlready = true
	return nil
}
