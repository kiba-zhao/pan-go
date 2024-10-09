package models

import "time"

type RemoteFileItem struct {
	ID                 string    `json:"id" form:"id"`
	CreatedAt          time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" form:"updatedAt"`
	NodeID             string    `json:"nodeId" form:"nodeId"`
	ItemID             uint      `json:"itemId" form:"itemId"`
	Name               string    `json:"name" form:"name"`
	FilePath           string    `json:"filePath" form:"filePath"`
	FileType           string    `json:"fileType" form:"fileType"`
	ParentPath         string    `json:"parentPath" form:"parentPath"`
	Size               int64     `json:"size" form:"size"`
	Available          bool      `json:"available" form:"available"`
	TagQuantity        uint      `json:"tagQuantity" form:"tagQuantity"`
	PendingTagQuantity uint      `json:"pendingTagQuantity" form:"pendingTagQuantity"`
}

type RemoteFileItemSearchCondition struct {
	NodeID     string  `json:"nodeId" form:"nodeId" binding:"required"`
	ItemID     uint    `json:"itemId" form:"itemId" binding:"required"`
	ParentPath *string `json:"parentPath" form:"parentPath" binding:"omitempty"`
}
