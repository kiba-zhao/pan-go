package models

import (
	"time"

	"gorm.io/gorm"
)

/**
 * Target File Model
 */
type TargetFile struct {
	ID             uint64         `gorm:"primaryKey" json:"id" form:"id"`
	CreatedAt      time.Time      `json:"createAt" form:"createAt"`
	UpdatedAt      time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
	TargetID       uint           `gorm:"index" json:"targetId" form:"targetId"`
	TargetHashCode string         `gorm:"size:28;index" json:"-" form:"-"`
	FilePath       string         ` json:"filepath" form:"filepath"`
	HashCode       string         `gorm:"size:28;index" json:"hashCode" form:"hashCode"`
	MimeType       string         `gorm:"size:100;index"  json:"mimeType" form:"mimeType"`
	Size           int64          `gorm:"index" json:"size" form:"size"`
	ModTime        time.Time      `json:"modTime" form:"modTime"`
	CheckSum       string         `gorm:"size:88:index" json:"checkSum" form:"checkSum"`
}
