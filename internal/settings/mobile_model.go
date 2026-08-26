//go:build android || ios

package settings

type MobileNetwork struct {
	WifiOnly *bool `json:"wifiOnly" form:"wifiOnly"  binding:"omitempty"`
}

type MobileNetworkFields struct {
	WifiOnly *bool `form:"wifiOnly" json:"wifiOnly"  binding:"omitempty"`
}
