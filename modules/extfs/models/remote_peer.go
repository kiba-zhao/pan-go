package models

import "gorm.io/gorm"

type RemotePeer struct {
	gorm.Model
	ID      string `gorm:"size:36,primarykey"`
	Enabled bool   `gorm:"index"`
}
