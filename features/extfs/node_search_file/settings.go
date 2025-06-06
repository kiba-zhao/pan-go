// Define node search file settings
package nodesearchfile

import (
	"time"
)

type NodeSearchFileSettings struct {
	ParallelThreshold uint16
	Lifecycle         uint64
}

func newDefaultSettings() *NodeSearchFileSettings {
	var settings NodeSearchFileSettings

	settings.ParallelThreshold = 3
	settings.Lifecycle = uint64(time.Minute * 20)
	// settings.Lifecycle = uint64(time.Second * 20)
	return &settings
}
