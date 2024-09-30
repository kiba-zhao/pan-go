package models

import "time"

type RemoteNode struct {
	NodeID             string    `json:"nodeId" form:"nodeId"`
	Name               string    `json:"name" form:"name"`
	Available          bool      `json:"available" form:"available"`
	TagQuantity        uint      `json:"tagQuantity" form:"tagQuantity"`
	PendingTagQuantity uint      `json:"pendingTagQuantity" form:"pendingTagQuantity"`
	CreatedAt          time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" form:"updatedAt"`
}
