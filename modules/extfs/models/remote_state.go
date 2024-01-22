package models

import (
	"time"

	"gorm.io/gorm"
)

type RemoteState struct {
	ID         string `gorm:"size:36,primarykey"`
	Hash       []byte `gorm:"size:64"`
	RemoteHash []byte `gorm:"size:64"`
	RemoteTime int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
