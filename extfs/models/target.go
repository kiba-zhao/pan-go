package models

import (
	"time"

	"gorm.io/gorm"
)

type TargetFields struct {
	Name     string `form:"name" binding:"required" json:"name"`
	FilePath string `form:"filepath" binding:"required" json:"filepath"`
	Enabled  bool   `form:"enabled" binding:"required" json:"enabled"`
}

type TargetQueryOptions struct {
	Version *uint8 `form:"version" binding:"omitempty" json:"version"`
}

type Target struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createAt" form:"createAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
	Name      string         `gorm:"size:255" json:"name" form:"name"`
	FilePath  string         `gorm:"size:255"  json:"filepath" form:"filepath"`
	Enabled   bool           `gorm:"index" json:"enabled" form:"enabled"`
	Invalid   *bool          `gorm:"index" json:"invalid" form:"invalid"`
	Version   uint8          `gorm:"index" json:"version" form:"version"`
}

type TargetSearchCondition struct {
	RangeSearchCondition
	SortSearchCondition
	Keyword string `form:"q" binding:"omitempty"`
	Enabled *bool  `form:"enabled" binding:"omitempty"`
}
