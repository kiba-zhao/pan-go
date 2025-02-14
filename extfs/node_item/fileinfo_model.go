package nodeitem

import "time"

type NodeFileInfo struct {
	ItemID     uint      `json:"itemId" form:"itemId"`
	Name       string    `json:"name" form:"name"`
	FilePath   string    `json:"filePath" form:"filePath"`
	FileType   string    `json:"fileType" form:"fileType"`
	ParentPath string    `json:"parentPath" form:"parentPath"`
	Size       int64     `json:"size" form:"size"`
	Available  bool      `json:"available" form:"available"`
	CreatedAt  time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" form:"updatedAt"`
}

type NodeFileInfoSearchCondition struct {
	ItemID     uint
	ParentPath string
}

type NodeFileInfoSelectCondition struct {
	ItemID     uint
	ParentPath string
	Name       string
}
