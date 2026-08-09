//go:build android || ios

package settings

import "pan/pkg/ptp"

type NetInterface = ptp.BroadcastInterface

type MobileSettingsConfig interface {
	NetInterfaces() []NetInterface
	WifiInterfaces() []NetInterface
	WifiZoneList() []string
}
