//go:build android || ios

package settings

type MobileSettingsConfig interface {
	WifiInterfaces() []NetInterface
}
