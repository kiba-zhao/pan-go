//go:build !(android || ios)

package settings

type WebHost struct {
	WebPort       uint16 `json:"webPort" form:"webPort"  binding:"omitempty"`
	WebEnabled    bool   `json:"webEnabled" form:"webEnabled"  binding:"omitempty"`
	LocalHostOnly bool   `json:"localHostOnly" form:"localHostOnly"  binding:"omitempty"`
}

type WebHostFields struct {
	WebPort       uint16 `form:"webPort" json:"webPort"  binding:"omitempty"`
	WebEnabled    *bool  `form:"webEnabled" json:"webEnabled"  binding:"omitempty"`
	LocalHostOnly *bool  `form:"localHostOnly" json:"localHostOnly"  binding:"omitempty"`
}
