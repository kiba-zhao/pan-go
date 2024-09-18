package models

import (
	"time"

	"gorm.io/gorm"
)

type NodeItemFields struct {
	Name     string `form:"name" binding:"required" json:"name"`
	FilePath string `form:"filepath" binding:"required" json:"filepath"`
	Enabled  *bool  `form:"enabled" binding:"required" json:"enabled"`
}

type NodeItem struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
	Name      string         `gorm:"size:255;index" json:"name" form:"name"`
	FilePath  string         `json:"filepath" form:"filepath"`
	FileType  string         `gorm:"size:1;index;" json:"filetype" form:"filetype"`
	Enabled   *bool          `gorm:"index" json:"enabled" form:"enabled"`
	Available bool           `gorm:"-:all" json:"available" form:"available"`
}
