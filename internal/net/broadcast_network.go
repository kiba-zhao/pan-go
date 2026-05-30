package net

import (
	"bytes"
	"errors"
	"net"
	"pan/internal/log"
	"slices"
	"strings"
	"sync"
)

var ErrBroadcastNetworkNoDeliverConn = errors.New("net.BroadcastNetwork Error: No Deliver Conn")
var ErrBroadcastNetworkUnavailable = errors.New("net.BroadcastNetwork Error: Unavailable")
var ErrBroadcastNetworkNoDeliverAddrs = errors.New("net.BroadcastNetwork Error: No Deliver Addrs")
var ErrBroadcastNetworkDeliverInvalidAddrs = errors.New("net.BroadcastNetwork Error: Deliver Invalid Addrs")

type stdBroadcastNetwork struct {
	logger log.Logger

	ifaces         []BroadcastInterface
	addrs          []string
	deliverMaxSize int
	zoneList       []string
	rw             sync.RWMutex
}

func (network *stdBroadcastNetwork) Deliver(connList []*net.UDPConn, payload []byte, addrs ...string) error {
	network.logger.Debug("net.BroadcastNetwork", "Deliver Addrs begin: "+strings.Join(addrs, ","))
	defer network.logger.Debug("net.BroadcastNetwork", "Deliver Addrs end: "+strings.Join(addrs, ","))

	if len(connList) <= 0 {
		return ErrBroadcastNetworkNoDeliverConn
	}

	network.rw.RLock()
	ifaces := network.ifaces
	maxSize := network.deliverMaxSize
	zoneList := network.zoneList
	deliverAddrs := network.addrs
	network.rw.RUnlock()
	if len(ifaces) <= 0 {
		return nil
	}

	size := len(payload) + 1
	if size > maxSize && maxSize > 0 {
		return bytes.ErrTooLarge
	}

	if len(addrs) > 0 {
		deliverAddrs = addrs
	}

	if len(deliverAddrs) <= 0 {
		return ErrBroadcastNetworkNoDeliverAddrs
	}

	deliverUDPAddrs := make([]*net.UDPAddr, 0)
	for _, addr := range deliverAddrs {
		udpAddrs, err := resolveAddrs(addr, zoneList)
		if err != nil {
			network.logger.Error("net.BroadcastNetwork", "Deliver resolveAddrs Error: "+addr)
			continue
		}
		if len(udpAddrs) <= 0 {
			continue
		}
		deliverUDPAddrs = append(deliverUDPAddrs, udpAddrs...)
	}

	if len(deliverUDPAddrs) <= 0 {
		return ErrBroadcastNetworkDeliverInvalidAddrs
	}

	buffer := packBuffer(payload)
	size = len(buffer)
	for _, iface := range ifaces {
		for _, conn := range connList {
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			localAddrIP := localAddr.IP
			if iface.Addr != localAddrIP.String() {
				continue
			}
			mtu := iface.MTU
			for _, udpAddr := range deliverUDPAddrs {
				if localAddrIP != nil {
					if localAddrIP.To4() != nil && udpAddr.IP.To4() == nil {
						continue
					}
					if localAddrIP.To4() == nil && udpAddr.IP.To4() != nil {
						continue
					}
				}
				for offset := 0; offset < size; offset += mtu {
					var limit int
					if offset+mtu > size {
						limit = size - offset
					} else {
						limit = offset + mtu
					}
					block := buffer[offset:limit]

					_, err := conn.WriteTo(block, udpAddr)
					if err != nil {
						network.logger.Error("net.BroadcastNetwork", "Deliver error: "+err.Error()+" Addr: "+udpAddr.String())
						break
					}
				}
			}
		}
	}

	return nil
}

func (network *stdBroadcastNetwork) Setup(cfg BroadcastConfig) {
	network.logger.Debug("net.BroadcastNetwork", "Setup begin")
	defer network.logger.Debug("net.BroadcastNetwork", "Setup end")

	network.rw.Lock()
	defer network.rw.Unlock()

	addrs := cfg.Addrs()
	deliverMaxSize := cfg.DeliverMaxSize()
	zoneList := cfg.IPv6ZoneList()
	ifaces := cfg.Interfaces()

	if !slices.Equal(network.addrs, addrs) {
		network.addrs = addrs
	}

	if network.deliverMaxSize != deliverMaxSize {
		network.deliverMaxSize = deliverMaxSize
	}

	if !slices.Equal(network.zoneList, zoneList) {
		network.zoneList = zoneList
	}

	network.ifaces = ifaces
}

func resolveAddrs(addr string, zoneArr []string) ([]*net.UDPAddr, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	if udpAddr.IP.To4() != nil || !udpAddr.IP.IsMulticast() || len(udpAddr.Zone) > 0 {
		return []*net.UDPAddr{udpAddr}, nil
	}

	if len(zoneArr) <= 0 {
		return nil, errors.New("net.BroadcastNetwork resolveAddrs Error: Unsupported IPv6 Zone")
	}

	udpAddrs := make([]*net.UDPAddr, 0)
	for _, zone := range zoneArr {
		udpAddr_ := net.UDPAddrFromAddrPort(udpAddr.AddrPort())
		udpAddr_.Zone = zone
		udpAddrs = append(udpAddrs, udpAddr_)
	}
	return udpAddrs, err
}
