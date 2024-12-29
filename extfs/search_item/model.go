package searchitem

import (
	"pan/app/web"
	"time"

	"gorm.io/gorm"
)

type SearchItemFields struct {
	Query string `json:"query" form:"query"`
}

type SearchItem struct {
	ID        uint64         `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
	Query     string         `gorm:"size:255;uniqueIndex" json:"query" form:"query"`
}

type SearchItemCondition struct {
	web.RangeSearchCondition
	Query string `form:"q" binding:"omitempty"`
}
