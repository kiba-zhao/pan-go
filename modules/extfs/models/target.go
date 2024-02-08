package models

import (
	"time"

	"gorm.io/gorm"
)

type Target struct {
	gorm.Model
	FilePath   string `gorm:"size:255"`
	Name       string `gorm:"size:255"`
	Size       int64
	ModifyTime time.Time
	Enabled    bool `gorm:"index"`
	Total      uint
}

type TargetSearchCondition struct {
	SearchCondition
	Keyword string `form:"keyword" binding:"omitempty"`
}

type TargetSearchResult = SearchResult[Target]
