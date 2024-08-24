package models

import (
	"pan/app/models"
	"time"

	"gorm.io/gorm"
)

/**
 * Target Model
 */

type TargetFields struct {
	Name     string `form:"name" binding:"required" json:"name"`
	FilePath string `form:"filepath" binding:"required" json:"filepath"`
	Enabled  *bool  `form:"enabled" binding:"required" json:"enabled"`
}

type TargetQueryOptions struct {
	Version *uint8 `form:"version" binding:"omitempty" json:"version"`
}

type Target struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
	Name      string         `gorm:"size:255;index" json:"name" form:"name"`
	FilePath  string         `json:"filepath" form:"filepath"`
	HashCode  string         `gorm:"size:28;index" json:"-" form:"-"`
	Enabled   *bool          `gorm:"index" json:"enabled" form:"enabled"`
	Available bool           `gorm:"-:all" json:"available" form:"available"`
	Version   uint8          `gorm:"index" json:"version" form:"version"`
}

type TargetSearchCondition struct {
	models.RangeSearchCondition
	models.SortSearchCondition
	Keyword   string `form:"q" binding:"omitempty"`
	Enabled   *bool  `form:"enabled" binding:"omitempty"`
	Available *bool  `form:"available" binding:"omitempty"`
}
