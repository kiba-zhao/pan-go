package settings

type Settings struct {
	Name             string   `json:"name" form:"name"  binding:"omitempty"`
	PeerPort         uint16   `json:"peerPort" form:"peerPort"  binding:"omitempty"`
	BroadcastAddrs   []string `json:"broadcastAddrs" form:"broadcastAddrs"  binding:"omitempty"`
	PublicAddrs      []string `json:"publicAddrs" form:"publicAddrs"  binding:"omitempty"`
	Enabled          bool     `json:"enabled" form:"enabled"  binding:"omitempty"`
	BroadcastEnabled bool     `json:"broadcastEnabled" form:"broadcastEnabled"  binding:"omitempty"`
}

type SettingsFields struct {
	Name             string   `form:"name" json:"name"  binding:"omitempty"`
	PeerPort         *uint16  `form:"peerPort" json:"peerPort"  binding:"omitempty"`
	BroadcastAddrs   []string `form:"broadcastAddrs" json:"broadcastAddrs"  binding:"omitempty"`
	PublicAddrs      []string `form:"publicAddrs" json:"publicAddrs"  binding:"omitempty"`
	Enabled          *bool    `form:"enabled" json:"enabled"  binding:"omitempty"`
	BroadcastEnabled *bool    `form:"broadcastEnabled" json:"broadcastEnabled"  binding:"omitempty"`
}
