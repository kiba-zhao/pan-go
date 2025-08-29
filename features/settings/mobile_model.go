//go:build android || ios

package settings

type MobileSettings struct {
	WifiOnly *bool  `json:"wifiOnly" form:"wifiOnly"  binding:"omitempty"`
	Version  uint32 `json:"-" form:"-"`
}

type MobileSettingsFields struct {
	WifiOnly *bool `form:"wifiOnly" json:"wifiOnly"  binding:"omitempty"`
}
