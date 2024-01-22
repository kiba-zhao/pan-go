package models

import (
	"time"

	"gorm.io/gorm"
)

type State struct {
	ID        uint64 `gorm:"primarykey"`
	Hash      []byte `gorm:"size:64"`
	Origin    uint8
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
