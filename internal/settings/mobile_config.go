//go:build android || ios

package settings

import "pan/pkg/net"

type NetInterface = net.BroadcastInterface

type MobileSettingsConfig interface {
	NetInterfaces() []NetInterface
	WifiInterfaces() []NetInterface
	WifiZoneList() []string
}
