package nodesearchfile

import (
	"os"
)

type NodeSearchFileSettings struct {
	DBPath    string
	TaskNum   uint16
	AutoClean bool
}

func newDefaultSettings() *NodeSearchFileSettings {
	var settings NodeSearchFileSettings

	settings.DBPath = os.TempDir()
	settings.TaskNum = 5
	settings.AutoClean = false
	return &settings
}
