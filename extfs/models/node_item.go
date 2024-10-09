package models

import (
	"time"

	"gorm.io/gorm"
)

type NodeItemFields struct {
	Name     string `form:"name" binding:"required" json:"name"`
	FilePath string `form:"filePath" binding:"required" json:"filePath"`
	Enabled  *bool  `form:"enabled" binding:"required" json:"enabled"`
}

type NodeItem struct {
	ID                 uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt          time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
	Name               string         `gorm:"size:255;index" json:"name" form:"name"`
	FilePath           string         `json:"filePath" form:"filePath"`
	FileType           string         `gorm:"size:1;index;" json:"fileType" form:"fileType"`
	Enabled            *bool          `gorm:"index" json:"enabled" form:"enabled"`
	Available          bool           `gorm:"-:all" json:"available" form:"available"`
	Size               int64          `gorm:"-:all" json:"size" form:"size"`
	TagQuantity        uint           `gorm:"-:all" json:"tagQuantity" form:"tagQuantity"`
	PendingTagQuantity uint           `gorm:"-:all" json:"pendingTagQuantity" form:"pendingTagQuantity"`
}
