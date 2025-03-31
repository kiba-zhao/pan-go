//go:build darwin

package vfs

import (
	"path"
)

func init() {
	MountPath = path.Join("/", "Volumes", "extfs")
}
