//go:build linux

package vfs

import (
	"os"
	"pan/lib/pkg"
	"path/filepath"
)

func initMountPathAsDefault() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(home, pkg.Name())
	return InitMountPath(filePath)
}
