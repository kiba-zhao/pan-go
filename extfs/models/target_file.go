package models

import (
	"time"

	"pan/app/models"

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
	HashCode       string         `gorm:"size:28;index" json:"hashCode" form:"hashCode"`
	FilePath       string         ` json:"filepath" form:"filepath"`
	MimeType       string         `gorm:"size:100;index"  json:"mimeType" form:"mimeType"`
	Size           int64          `gorm:"index" json:"size" form:"size"`
	ModTime        time.Time      `json:"modTime" form:"modTime"`
	CheckSum       string         `gorm:"size:88:index" json:"checkSum" form:"checkSum"`
	Available      bool           `gorm:"-:all" json:"available" form:"available"`
	Target         Target         `gorm:"foreignKey:TargetID" json:"-" form:"-"`
}

type TargetFileSearchCondition struct {
	models.RangeSearchCondition
	models.SortSearchCondition
	Keyword   string `form:"q" binding:"omitempty"`
	Available *bool  `form:"available" binding:"omitempty"`
	TargetID  uint   `form:"targetId" binding:"omitempty" json:"targetId"`
}
