package node

import (
	"pan/lib/web"
	"time"

	"gorm.io/gorm"
)

type AppNodeFields struct {
	PeerID       string   `json:"peerId" form:"peerId"`
	Name         string   `json:"name" form:"name"`
	Blocked      *bool    `json:"blocked" form:"blocked"`
	NetworkAddrs []string `json:"networkAddrs" form:"networkAddrs"`
}

type AppNode struct {
	ID               uint           `gorm:"primarykey" json:"id" form:"id"`
	PeerID           string         `gorm:"size:256;index" json:"peerId" form:"peerId"`
	Name             string         `gorm:"size:256;index" json:"name" form:"name"`
	Blocked          bool           `gorm:"index" json:"blocked" form:"blocked"`
	Online           bool           `gorm:"-:all" json:"online" form:"online"`
	NetworkAddrTexts []string       `gorm:"-:all" json:"networkAddrs" form:"networkAddrs"`
	NetworkAddrs     []NetworkAddr  `gorm:"foreignKey:AppNodeID" json:"-" form:"-"`
	CreatedAt        time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
}

type AppNodeSearchCondition struct {
	web.RangeCondition
	web.SortCondition
	Keyword string `form:"q" binding:"omitempty"`
	Blocked *bool  `form:"blocked" binding:"omitempty"`
	Online  *bool  `form:"online" binding:"omitempty"`
}
