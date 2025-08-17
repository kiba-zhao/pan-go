package broadcast

import (
	"bytes"
	"errors"
	"net"
	"pan/lib/log"
)

var ErrBroadcastClusterUnavailable = errors.New("broadcast.BroadcastCluster Error: Unavailable")
var ErrBroadcastClusterDeliverNoAddrs = errors.New("broadcast.BroadcastCluster Error: Deliver No Addrs")
var ErrBroadcastClusterDeliverInvalidAddrs = errors.New("broadcast.BroadcastCluster Error: Deliver Invalid Addrs")

// BroadcastServeModule is the broadcast serve module
type BroadcastServeModule interface {
	// ServeBroadcast serves the broadcast message to the peer.
	ServeBroadcast([]byte, string) error
}

type BroadcastClusterRuntime interface {
	// ServeModules returns the list of broadcast serve modules.
	ServeModules() []BroadcastServeModule
	// Addrs returns the list of broadcast deliver addresses.
	Addrs() []string
	// DeliverLimitSize returns the broadcast deliver limit size.
	DeliverLimitSize() int
	// MTU returns the broadcast MTU.
	MTU() int

	DeliverConn() *net.UDPConn
}

type BroadcastCluster interface {
	// Serve serves the broadcast message to the peer.
	Serve([]byte, string) error

	// Deliver delivers the broadcast message to the peer.
	Deliver([]byte, ...string) error

	Runtime() BroadcastClusterRuntime
}

type stdBroadcastCluster struct {
	logger  log.Logger
	runtime BroadcastClusterRuntime
}

var _ = (BroadcastCluster)((*stdBroadcastCluster)(nil))

func (cluster *stdBroadcastCluster) Serve(payload []byte, addr string) error {
	cluster.logger.Debug("BroadcastCluster", "Serve Addr begin: "+addr)
	defer cluster.logger.Debug("BroadcastCluster", "Serve Addr end: "+addr)

	clusterRuntime := cluster.Runtime()
	if clusterRuntime == nil {
		return ErrBroadcastClusterUnavailable
	}
	serveModules := clusterRuntime.ServeModules()
	if len(serveModules) <= 0 {
		return ErrBroadcastClusterUnavailable
	}

	for _, module := range serveModules {
		err := module.ServeBroadcast(payload, addr)
		if err == nil {
			break
		}
	}

	return nil
}

func (cluster *stdBroadcastCluster) Deliver(payload []byte, addrs ...string) error {
	clusterRuntime := cluster.Runtime()
	if clusterRuntime == nil {
		return ErrBroadcastClusterUnavailable
	}
	conn := clusterRuntime.DeliverConn()
	if conn == nil {
		return ErrBroadcastClusterUnavailable
	}

	size := len(payload) + 1
	limitSize := clusterRuntime.DeliverLimitSize()
	if size > limitSize && limitSize > 0 {
		return bytes.ErrTooLarge
	}

	var deliverAddrs []string
	if len(addrs) <= 0 {
		deliverAddrs = clusterRuntime.Addrs()
	} else {
		deliverAddrs = addrs
	}

	if len(deliverAddrs) <= 0 {
		return ErrBroadcastClusterDeliverNoAddrs
	}

	deliverUDPAddrs := make([]*net.UDPAddr, 0)
	for _, addr := range deliverAddrs {
		udpAddrs, err := resolveAddrs(addr)
		if err != nil {
			cluster.logger.Error("BroadcastCluster", "Deliver resolveAddrs error: "+addr)
			continue
		}
		if len(udpAddrs) <= 0 {
			continue
		}
		deliverUDPAddrs = append(deliverUDPAddrs, udpAddrs...)
	}

	if len(deliverUDPAddrs) <= 0 {
		return ErrBroadcastClusterDeliverInvalidAddrs
	}

	mtu := clusterRuntime.MTU()
	buffer := packBuffer(payload)
	size = len(buffer)
	for offset := 0; offset < size; offset += mtu {
		var limit int
		if offset+mtu > size {
			limit = size - offset
		} else {
			limit = offset + mtu
		}
		block := buffer[offset:limit]

		for _, udpAddr := range deliverUDPAddrs {
			_, err := conn.WriteTo(block, udpAddr)
			if err != nil {
				cluster.logger.Error("BroadcastCluster", "Deliver error: "+err.Error())
			}
		}
	}

	return nil
}

func (cluster *stdBroadcastCluster) Runtime() BroadcastClusterRuntime {
	return cluster.runtime
}

func resolveAddrs(addr string) ([]*net.UDPAddr, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	if udpAddr.IP.To4() != nil || !udpAddr.IP.IsMulticast() || len(udpAddr.Zone) > 0 {
		return []*net.UDPAddr{udpAddr}, nil
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	isGlobalAddr := isGlobalMulticastIP(udpAddr.IP)
	udpAddrs := make([]*net.UDPAddr, 0)
	for _, iface := range ifaces {
		if net.FlagMulticast != (net.FlagMulticast & iface.Flags) {
			continue
		}
		if net.FlagRunning != (net.FlagRunning & iface.Flags) {
			continue
		}
		if net.FlagLoopback == (net.FlagLoopback & iface.Flags) {
			continue
		}
		ifaceAddrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, ifaceAddr := range ifaceAddrs {
			ipNet, ok := ifaceAddr.(*net.IPNet)
			if !ok || ipNet.IP.To4() != nil {
				continue
			}
			if isGlobalAddr == (ipNet.IP.IsGlobalUnicast() && !ipNet.IP.IsPrivate()) {
				udpAddr_ := net.UDPAddrFromAddrPort(udpAddr.AddrPort())
				udpAddr_.Zone = iface.Name
				udpAddrs = append(udpAddrs, udpAddr_)
				break
			}
		}
	}
	return udpAddrs, err
}

func isGlobalMulticastIP(ip net.IP) bool {
	return !(ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast())
}
