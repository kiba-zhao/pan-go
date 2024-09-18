package models

type Item struct {
	ID                 string  `json:"id" form:"id"`
	ItemType           string  `json:"itemType" form:"itemType"`
	Name               string  `json:"name" form:"name"`
	TagQuantity        uint    `json:"tagQuantity" form:"tagQuantity"`
	PendingTagQuantity uint    `json:"pendingTagQuantity" form:"pendingTagQuantity"`
	Available          bool    `json:"available" form:"available"`
	LinkID             *string `json:"linkId" form:"linkId"`
	ParentID           *string `json:"parentId" form:"parentId"`
}

type ItemSearchCondition struct {
	ItemTypes []string `form:"itemType" binding:"omitempty"`
	ParentID  *string  `form:"parentId" binding:"omitempty"`
}
