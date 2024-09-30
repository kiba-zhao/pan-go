package models

import "time"

type RemoteNodeItem struct {
	ID                 string    `json:"id" form:"id"`
	CreatedAt          time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" form:"updatedAt"`
	NodeID             string    `json:"nodeId" form:"nodeId"`
	ItemID             uint      `json:"itemId" form:"itemId"`
	Name               string    `json:"name" form:"name"`
	FileType           string    `json:"fileType" form:"fileType"`
	Size               int64     `json:"size" form:"size"`
	Available          bool      `json:"available" form:"available"`
	TagQuantity        uint      `json:"tagQuantity" form:"tagQuantity"`
	PendingTagQuantity uint      `json:"pendingTagQuantity" form:"pendingTagQuantity"`
}

type RemoteNodeItemSearchCondition struct {
	NodeID string `json:"nodeId" form:"nodeId" binding:"required"`
}
