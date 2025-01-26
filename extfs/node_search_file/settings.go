package nodesearchfile

import (
	"os"
	"time"
)

type NodeSearchFileSettings struct {
	DBPath    string
	TaskNum   uint16
	Lifecycle uint64
}

func newDefaultSettings() *NodeSearchFileSettings {
	var settings NodeSearchFileSettings

	settings.DBPath = os.TempDir()
	settings.TaskNum = 5
	// settings.Lifecycle = uint64(time.Minute * 20)
	settings.Lifecycle = uint64(time.Second * 20)
	return &settings
}
