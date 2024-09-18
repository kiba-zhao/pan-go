package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint64         `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
	Name      string         `gorm:"size:255;index" json:"name" form:"name"`
	NameHash  []byte         `gorm:"-" json:"-" form:"-"`
}
