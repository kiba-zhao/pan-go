// Define node search file model
package nodesearchfile

import (
	"pan/lib/web"
	"time"

	"gorm.io/gorm"
)

type Tokens []string

type NodeSearchFile struct {
	ID        uint64    `gorm:"primarykey" json:"id" form:"id"`
	ItemID    uint      `json:"itemId" form:"itemId"`
	Name      string    `json:"name" form:"name"`
	FilePath  string    `json:"filePath" form:"filePath"`
	FileType  string    `json:"fileType" form:"fileType"`
	MimeType  string    `json:"mimeType" form:"mimeType"`
	Tokens    Tokens    `gorm:"type:text" json:"tokens" form:"tokens"`
	Score     uint      `json:"score" form:"score"`
	Size      int64     `json:"size" form:"size"`
	Available bool      `gorm:"-:all" json:"available" form:"available"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`
}

type NodeSearchFileSearchCondition struct {
	web.RangeCondition
	Query string `form:"query" binding:"omitempty"`
	Hash  string `form:"hash" binding:"omitempty"`
}

type NodeSearchTaskFields struct {
	Query string `json:"query" form:"query"`
	Hash  string `json:"hash" form:"hash"`
}

type NodeSearchTask struct {
	ID        uint64    `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt
	Query     string `gorm:"size:256,Index" json:"query" form:"query"`
	Hash      string `gorm:"size:256,Index" json:"hash" form:"hash"`
	Status    uint8  `gorm:"Index" json:"status" form:"status"`
}

const (
	NodeSearchTaskStatusPending uint8 = iota
	NodeSearchTaskStatusSuccess
	NodeSearchTaskStatusWarning
	NodeSearchTaskStatusCancel
)
