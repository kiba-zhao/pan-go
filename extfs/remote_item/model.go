package remoteitem

import "time"

type RemoteItem struct {
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`
	PeerID    string    `json:"peerId" form:"peerId"`
	ItemID    uint      `json:"itemId" form:"itemId"`
	Name      string    `json:"name" form:"name"`
	FileType  string    `json:"fileType" form:"fileType"`
	Size      int64     `json:"size" form:"size"`
	Available bool      `json:"available" form:"available"`
}

type RemoteItemSearchCondition struct {
	PeerID string `json:"peerId" form:"peerId" binding:"required"`
}
