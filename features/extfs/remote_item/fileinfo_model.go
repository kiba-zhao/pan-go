package remoteitem

type RemoteFileInfo struct {
	RemoteFileStat
	PeerID     string `json:"peerId" form:"peerId"`
	ItemID     uint   `json:"itemId" form:"itemId"`
	Name       string `json:"name" form:"name"`
	FilePath   string `json:"filePath" form:"filePath"`
	ParentPath string `json:"parentPath" form:"parentPath"`
	Available  bool   `json:"available" form:"available"`
}

type RemoteFileInfoSearchCondition struct {
	PeerID     string `json:"peerId" form:"peerId" binding:"required"`
	ItemID     uint   `json:"itemId" form:"itemId" binding:"required"`
	ParentPath string `json:"parentPath" form:"parentPath" binding:"omitempty"`
}
