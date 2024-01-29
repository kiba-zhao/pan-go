package models

import (
	"time"

	"gorm.io/gorm"
)

type FileInfo struct {
	gorm.Model
	TargetID     uint
	Name         string `gorm:"size:255"`
	Size         int64
	ModifyTime   time.Time
	RelativePath string `gorm:"size:255"`
	Hash         []byte `gorm:"size:64"`
}

type FileInfoTotal struct {
	Size  int64
	Total uint
}
