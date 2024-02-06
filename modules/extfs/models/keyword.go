package models

import (
	"gorm.io/gorm"
)

type Keyword struct {
	gorm.Model
	Name   string `gorm:"size:255,index"`
	Origin uint8  `gorm:"index"`
}

type KeywordFindCondition struct {
	Keyword string `form:"keyword" binding:"omitempty"`
	Offset  int    `form:"offset" binding:"omitempty"`
	Limit   int    `form:"limit" binding:"omitempty"`
}
