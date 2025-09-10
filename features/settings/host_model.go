//go:build !(android || ios)

package settings

type HostSettings struct {
	WebAddr    string `json:"webAddr" form:"webAddr"  binding:"omitempty"`
	WebEnabled bool   `json:"webEnabled" form:"webEnabled"  binding:"omitempty"`
}

type HostSettingsFields struct {
	WebAddr    string `form:"webAddr" json:"webAddr"  binding:"omitempty"`
	WebEnabled *bool  `form:"webEnabled" json:"webEnabled"  binding:"omitempty"`
}
