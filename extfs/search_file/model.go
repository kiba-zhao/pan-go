package searchfile

import (
	"time"
)

type SearchFile struct {
	ID        uint64    `json:"id" form:"id"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`
	Name      string    `json:"name" form:"name"`
	FileType  string    `json:"fileType" form:"fileType"`
	Size      int64     `json:"size" form:"size"`
	Available bool      `json:"available" form:"available"`
	PeerID    string    `json:"peerId" form:"peerId"`
	ItemID    uint      `json:"itemId" form:"itemId"`
	FilePath  string    `json:"filePath" form:"filePath"`
	SearchID  uint64    `json:"searchId" form:"searchId"`
}
