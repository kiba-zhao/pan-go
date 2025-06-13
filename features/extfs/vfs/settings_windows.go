//go:build windows

package vfs

import (
	"os"
)

func initMountPathAsDefault() error {
	var err error
	var filePath string
	for i := 'z'; i >= 'd'; i-- {
		filePath = string(i) + ":"
		if _, err = os.Stat(filePath); os.IsNotExist(err) {
			break
		}
	}

	if !os.IsNotExist(err) {
		return err
	}
	return InitMountPath(filePath)
}
