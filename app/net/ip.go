package net

import "net"

type AddrStat struct {
	IPv4Enabled       bool
	IPv6Enabled       bool
	IPv6GlobalEnabled bool
	IPv4GlobalEnabled bool
}

func StatAddr() (*AddrStat, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var addrStat AddrStat
	for _, addr := range addrs {
		ipNet := addr.(*net.IPNet)
		if ipNet.IP.IsLoopback() {
			continue
		}
		if ipNet.IP.IsPrivate() {
			continue
		}
		if ipNet.IP.IsLinkLocalUnicast() {
			continue
		}
		if ipNet.IP.IsMulticast() || ipNet.IP.IsLinkLocalMulticast() || ipNet.IP.IsInterfaceLocalMulticast() {
			continue
		}

		if ipNet.IP.To4() != nil {
			addrStat.IPv4GlobalEnabled = ipNet.IP.IsGlobalUnicast()
			addrStat.IPv4Enabled = true
			continue
		}

		addrStat.IPv6Enabled = true
		addrStat.IPv6GlobalEnabled = ipNet.IP.IsGlobalUnicast()

	}

	return &addrStat, nil
}
