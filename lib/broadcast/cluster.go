package broadcast

import (
	"bytes"
	"errors"
	"iter"
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

	IPv6ZoneList() []string

	SeqForDeliverConn() iter.Seq2[int, *net.UDPConn]
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
	connSeq := clusterRuntime.SeqForDeliverConn()
	if connSeq == nil {
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
		udpAddrs, err := resolveAddrs(addr, clusterRuntime.IPv6ZoneList())
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

	buffer := packBuffer(payload)
	size = len(buffer)
	for mtu, conn := range connSeq {

		localAddr := conn.LocalAddr().(*net.UDPAddr)
		localAddrIP := localAddr.IP

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
					cluster.logger.Error("BroadcastCluster", "Deliver error: "+err.Error()+" Addr: "+udpAddr.String())
				}
			}

		}
	}

	return nil
}

func (cluster *stdBroadcastCluster) Runtime() BroadcastClusterRuntime {
	return cluster.runtime
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
		return nil, errors.New("Unsupported IPv6 Zone")
	}

	udpAddrs := make([]*net.UDPAddr, 0)
	for _, zone := range zoneArr {
		udpAddr_ := net.UDPAddrFromAddrPort(udpAddr.AddrPort())
		udpAddr_.Zone = zone
		udpAddrs = append(udpAddrs, udpAddr_)
	}
	return udpAddrs, err
}

func isGlobalMulticastIP(ip net.IP) bool {
	return !(ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast())
}
