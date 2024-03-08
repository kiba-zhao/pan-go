package extfs

import (
	"os"
	"pan/extfs/models"
	"path"
)

// TODO: defines a struct named Settings with one fields: Enabled of type bool

type Settings struct {
	models.Settings
	Enabled bool
}

func (s *Settings) init() {
	s.Enabled = true
	s.Settings.TotalHeaderName = "X-Total-Count"
	s.Settings.DBFilePath = path.Join(os.TempDir(), "extfs.db")
}
