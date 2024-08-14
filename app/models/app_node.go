package models

import (
	"time"

	"gorm.io/gorm"
)

type AppNodeFields struct {
	NodeID  string `gorm:"size:255;json:"nodeId" form:"nodeId"`
	Name    string `gorm:"size:255;json:"name" form:"name"`
	Blocked *bool  `gorm:"index" json:"blocked" form:"blocked"`
}

type AppNode struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	NodeID    string         `gorm:"size:255;uniqueIndex" json:"nodeId" form:"nodeId"`
	Name      string         `gorm:"size:255;index" json:"name" form:"name"`
	Blocked   *bool          `gorm:"index" json:"blocked" form:"blocked"`
	Online    bool           `gorm:"-:all" json:"online" form:"online"`
	CreatedAt time.Time      `json:"createAt" form:"createAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
}

type AppNodeSearchCondition struct {
	RangeSearchCondition
	SortSearchCondition
	Keyword string `form:"q" binding:"omitempty"`
	Blocked *bool  `form:"blocked" binding:"omitempty"`
	Online  *bool  `form:"online" binding:"omitempty"`
}
