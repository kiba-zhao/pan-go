//go:build darwin

package vfs

import (
	"pan/lib/pkg"
	"path"
)

func initMountPathAsDefault() error {
	filePath := path.Join("/", "Volumes", pkg.Name())
	return InitMountPath(filePath)
}
