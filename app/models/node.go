package models

import (
	"time"

	"gorm.io/gorm"
)

type NodeFields struct {
	NodeID  string `json:"nodeId" form:"nodeId"`
	Name    string `json:"name" form:"name"`
	Blocked *bool  `json:"blocked" form:"blocked"`
}

type Node struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	NodeID    string         `gorm:"size:255;uniqueIndex" json:"nodeId" form:"nodeId"`
	Name      string         `gorm:"size:255;index" json:"name" form:"name"`
	Blocked   bool           `gorm:"index" json:"blocked" form:"blocked"`
	Online    bool           `gorm:"-:all" json:"online" form:"online"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
}

type NodeSearchCondition struct {
	RangeSearchCondition
	SortSearchCondition
	Keyword string `form:"q" binding:"omitempty"`
	Blocked *bool  `form:"blocked" binding:"omitempty"`
	Online  *bool  `form:"online" binding:"omitempty"`
}
