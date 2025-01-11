package remotefile

import "time"

type RemoteFile struct {
	ID                 string    `json:"id" form:"id"`
	CreatedAt          time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" form:"updatedAt"`
	PeerID             string    `json:"peerId" form:"peerId"`
	ItemID             uint      `json:"itemId" form:"itemId"`
	Name               string    `json:"name" form:"name"`
	FilePath           string    `json:"filePath" form:"filePath"`
	FileType           string    `json:"fileType" form:"fileType"`
	MimeType           string    `json:"mimeType" form:"mimeType"`
	ParentPath         string    `json:"parentPath" form:"parentPath"`
	Size               int64     `json:"size" form:"size"`
	Available          bool      `json:"available" form:"available"`
	TagQuantity        uint      `json:"tagQuantity" form:"tagQuantity"`
	PendingTagQuantity uint      `json:"pendingTagQuantity" form:"pendingTagQuantity"`
}

type RemoteFileSearchCondition struct {
	PeerID     string  `json:"peerId" form:"peerId" binding:"required"`
	ItemID     uint    `json:"itemId" form:"itemId" binding:"required"`
	ParentPath *string `json:"parentPath" form:"parentPath" binding:"omitempty"`
}
