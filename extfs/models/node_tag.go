package models

import (
	"time"

	"gorm.io/gorm"
)

type NodeTagReport struct {
	NodeID          string
	Quantity        uint
	PendingQuantity uint
}

type NodeTag struct {
	ID        uint64         `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
	NodeID    string         `gorm:"size:255;index" json:"nodeId" form:"nodeId"`
	TagID     uint64         `gorm:"index" json:"tagId" form:"tagId"`
	Tag       Tag            `gorm:"foreignKey:TagID" json:"tag" form:"tag"`
	Pending   *bool          `gorm:"index" json:"pending" form:"pending"`
	IsValid   *bool          `gorm:"index" json:"isValid" form:"isValid"`
}
