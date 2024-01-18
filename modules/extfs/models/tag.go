package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name  string `gorm:"size:255,index"`
	Owner string `gorm:"size:36,index"`
}

type TagFindCondition struct {
	Name   string `form:"name" binding:"omitempty,excludesall=%"`
	Owner  string `form:"owner" binding:"omitempty,uuid"`
	Offset int    `form:"offset" binding:"omitempty"`
	Limit  int    `form:"limit" binding:"omitempty"`
}
