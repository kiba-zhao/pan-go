package appbroadcast

import (
	"time"

	"gorm.io/gorm"
)

type AppBroadcastInfo struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	PeerID    string         `gorm:"size:256;index" json:"peerId" form:"peerId"`
	Hightest  uint64         `json:"hightest" form:"hightest"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
}
