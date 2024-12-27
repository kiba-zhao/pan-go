package searchfile

import (
	"time"

	"gorm.io/gorm"
)

type SearchFile struct {
	ID        uint64         `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" form:"deletedAt"`
	ReferId   string         `json:"referId" form:"referId"`
	ReferType string         `gorm:"size:2;" json:"referType" form:"referType"`
	Name      string         `gorm:"-" json:"name" form:"name"`
	FileType  string         `gorm:"-" json:"fileType" form:"fileType"`
	Size      int64          `gorm:"-" json:"size" form:"size"`
	Available bool           `gorm:"-" json:"available" form:"available"`
}
