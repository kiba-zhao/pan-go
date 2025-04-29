package node

type NetworkAddr struct {
	ID        uint64  `gorm:"primarykey" json:"id" form:"id"`
	AppNodeID uint    `gorm:"uniqueIndex:idx_network_addr_uniq" json:"appNodeId" form:"appNodeId"`
	Address   string  `gorm:"size:256;uniqueIndex:idx_network_addr_uniq" json:"address" form:"address"`
	CreatedAt string  `json:"createdAt" form:"createdAt"`
	UpdatedAt string  `json:"updatedAt" form:"updatedAt"`
	AppNode   AppNode `gorm:"foreignKey:AppNodeID"`
}

type NetworkAddrFields struct {
	ID      uint64 `gorm:"primarykey" json:"id" form:"id"`
	Address string `json:"address" form:"address"`
}
