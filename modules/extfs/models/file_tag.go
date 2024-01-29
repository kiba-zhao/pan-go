package models

import "gorm.io/gorm"

type FileTag struct {
	gorm.Model
	FileID uint
	TagID  uint
	Hash   []byte `gorm:"size:64"`
}
