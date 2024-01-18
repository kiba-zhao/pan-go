package models

import "gorm.io/gorm"

type Peer struct {
	gorm.Model
	ID         string `gorm:"size:36,primarykey"`
	Enabled    bool   `gorm:"index"`
	Hash       []byte `gorm:"size:64"`
	RemoteHash []byte `gorm:"size:64"`
	RemoteTime int64
}

type PeerFindCondition struct {
	Enabled bool `form:"disabled" binding:"omitempty"`
	Offset  int  `form:"offset" binding:"omitempty"`
	Limit   int  `form:"limit" binding:"omitempty"`
}
