//go:build windows

package vfs

import (
	"os"
)

func generateDefaultMountPath() string {
	for i := 'z'; i >= 'd'; i-- {
		mountPath := string(i) + ":"
		if _, err := os.Stat(mountPath); os.IsNotExist(err) {
			return mountPath
		}
	}
	return ""
}

func init() {
	MountPath = generateDefaultMountPath()
}
