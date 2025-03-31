//go:build linux

package vfs

import (
	"os"
	"path"
)

func generateDefaultMountPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path.Join("/", "media")
	}
	return path.Join(home, "extfs")
}

func init() {
	MountPath = generateDefaultMountPath()
}
