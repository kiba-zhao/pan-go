package settings

type DeviceInfo struct {
	Name string `json:"name" form:"name"  binding:"omitempty"`
	Memo string `json:"memo" form:"memo"  binding:"omitempty"`
}

type DeviceInfoFields struct {
	Name string  `form:"name" json:"name"  binding:"omitempty"`
	Memo *string `form:"memo" json:"memo"  binding:"omitempty"`
}

type DeviceNetwork struct {
	Enabled          bool     `json:"enabled" form:"enabled"  binding:"omitempty"`
	Port             uint16   `json:"port" form:"port"  binding:"omitempty"`
	PublicAddrs      []string `json:"publicAddrs" form:"publicAddrs"  binding:"omitempty"`
	BroadcastAddrs   []string `json:"broadcastAddrs" form:"broadcastAddrs"  binding:"omitempty"`
	BroadcastEnabled bool     `json:"broadcastEnabled" form:"broadcastEnabled"  binding:"omitempty"`
}

type DeviceNetworkFields struct {
	Enabled          *bool    `form:"enabled" json:"enabled"  binding:"omitempty"`
	Port             *uint16  `form:"port" json:"port"  binding:"omitempty"`
	PublicAddrs      []string `form:"publicAddrs" json:"publicAddrs"  binding:"omitempty"`
	BroadcastAddrs   []string `form:"broadcastAddrs" json:"broadcastAddrs"  binding:"omitempty"`
	BroadcastEnabled *bool    `form:"broadcastEnabled" json:"broadcastEnabled"  binding:"omitempty"`
}
