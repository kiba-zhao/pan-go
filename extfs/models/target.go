package models

import "gorm.io/gorm"

type Target struct {
	gorm.Model
	Name     string `gorm:"size:255"`
	FilePath string `gorm:"size:255"`
	Enabled  bool   `gorm:"index"`
	Version  uint8  `gorm:"index"`
}

type TargetSearchCondition struct {
	RangeSearchCondition
	SortSearchCondition
	Keyword string `form:"keyword" binding:"omitempty"`
	Enabled *bool  `form:"enabled" binding:"omitempty"`
}
