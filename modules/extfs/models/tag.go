package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name   string `gorm:"size:255,index"`
	Origin uint8  `gorm:"index"`
}

type TagFindCondition struct {
	Name   string `form:"name" binding:"omitempty,excludesall=%"`
	Offset int    `form:"offset" binding:"omitempty"`
	Limit  int    `form:"limit" binding:"omitempty"`
}
