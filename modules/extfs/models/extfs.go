package models

import "gorm.io/gorm"

type ExtFS struct {
	gorm.Model
	ID   uint64 `gorm:"primaryKey"`
	Hash []byte `gorm:"size:64,index"`
}
