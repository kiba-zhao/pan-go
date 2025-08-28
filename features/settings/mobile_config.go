//go:build android || ios

package settings

type MobileSettingsConfig interface {
	WifiInterface() []NetInterface
}
