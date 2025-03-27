// Define broadcast server for broadcast
package broadcast

import (
	"bytes"
	"context"
	"errors"
	"net"
	"slices"
	"strconv"
	"sync"
	"time"

	appNet "pan/app/net"
	"pan/logger"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var ErrBroadcastServerUnavailable = errors.New("broadcast.BroadcastServer Error: Unavailable")
var ErrBroadcastInvalidConn = errors.New("broadcast.BroadcastServer Error: Invalid Connection")

type PacketConn interface {
	JoinGroup(ifi *net.Interface, group net.Addr) error
	LeaveGroup(ifi *net.Interface, group net.Addr) error
}

type broadcastServer struct {
	module  *broadcastModule
	locker  sync.RWMutex
	connArr []*net.UDPConn
	mtu     int
}

// Shutdown shuts down the broadcast server.
func (bs *broadcastServer) Shutdown() error {

	bs.locker.Lock()
	defer bs.locker.Unlock()
	return bs.close()
}

// ListenAndServe starts the broadcast server to listen on the specified addresses for incoming multicast messages.
// It initializes the necessary UDP connections for each address and serves the multicast groups.
// Returns an error if the broadcast server is unavailable or if there are issues with the provided addresses.
func (bs *broadcastServer) ListenAndServe(addrs []string) error {
	bs.locker.Lock()
	defer bs.locker.Unlock()

	if bs.module == nil {
		return ErrBroadcastServerUnavailable
	}

	if len(bs.connArr) > 0 {
		return bs.close()
	}

	mAddrs := make(map[int][]net.IP)
	for _, addr := range addrs {
		var addrPort int
		host, port, err := net.SplitHostPort(addr)
		if err == nil {
			addrPort, err = strconv.Atoi(port)
		}
		if err != nil {
			logger.Default().Log(context.Background(), logger.LevelWarning, "app.broadcastServer Addr Error: "+err.Error())
			continue
		}
		if _, ok := mAddrs[addrPort]; !ok {
			mAddrs[addrPort] = make([]net.IP, 0)
		}
		ip := net.ParseIP(host)
		mAddrs[addrPort] = append(mAddrs[addrPort], ip)
	}

	for port, groups := range mAddrs {
		conn, err := newMulticastConn(port, groups)
		if err != nil {
			logger.Default().Log(context.Background(), logger.LevelWarning, "app.broadcastServer Conn Error: "+err.Error())
			continue
		}
		go bs.serve(conn)
	}

	return nil
}

// serve serves the UDP connection.
// It reads the UDP packets and parses them using the parsePacketBuffer function.
// If the packet is a complete packet, it will be served using the Serve method of the broadcast module.
// If the packet is incomplete, it will be stored in the packet buffer and wait for the complete packet.
// If the packet is invalid, it will be ignored.
// The function will return an error if the connection is invalid or if there is an error in the connection.
func (bs *broadcastServer) serve(conn *net.UDPConn) error {
	packetBuffers := make([]*PacketBuffer, 0)
	var packetBuffersRW sync.RWMutex
	var err error

	for {
		block := make([]byte, bs.mtu)
		byteLen, addr, err := conn.ReadFromUDP(block)
		if errors.Is(err, net.ErrClosed) {
			break
		}
		if err != nil {
			continue
		}

		buffer, size := parsePacketBuffer(block[:byteLen])
		packetBuffersRW.RLock()
		idx, ok := slices.BinarySearchFunc(packetBuffers, addr.String(), comparePacketBuffer)
		var bufferItem *PacketBuffer
		if ok {
			bufferItem = packetBuffers[idx]
		}
		packetBuffersRW.RUnlock()

		if size > 0 && ok {
			bufferItem.cancel()
			bufferItem.wg.Wait()
			ok = false
		} else if size == 0 {
			if !ok {
				continue
			}
			len_ := len(bufferItem.content) + len(block)
			if len_ > bufferItem.size {
				continue
			}
			buffer = bytes.Join([][]byte{bufferItem.content, block}, nil)
			size = bufferItem.size
		}

		if size > len(buffer) {
			if !ok {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				bufferItem = &PacketBuffer{
					addr:   addr.String(),
					cancel: cancel,
				}
				bufferItem.wg.Add(1)
				packetBuffersRW.Lock()
				idx_, ok_ := slices.BinarySearchFunc(packetBuffers, addr.String(), comparePacketBuffer)
				if !ok_ {
					packetBuffers = slices.Insert(packetBuffers, idx_, bufferItem)
				}
				packetBuffersRW.Unlock()
				go func(item *PacketBuffer, ctx context.Context) {
					defer item.wg.Done()
					<-ctx.Done()
					packetBuffersRW.Lock()
					defer packetBuffersRW.Unlock()
					idx_, ok_ := slices.BinarySearchFunc(packetBuffers, item.addr, comparePacketBuffer)
					if ok_ {
						packetBuffers = slices.Delete(packetBuffers, idx_, idx_+1)
					}
				}(bufferItem, ctx)
			}

			bufferItem.size = size
			bufferItem.content = buffer
			continue
		}

		if ok {
			bufferItem.cancel()
			bufferItem.wg.Wait()
		}
		go bs.module.Serve(buffer, addr.String())
	}
	return err
}

// close closes all UDP connections in the broadcast server's connection array.
// If no connections exist, it returns nil. Otherwise, it iterates through the
// connections, closing each one and then resets the connection array to nil.

func (bs *broadcastServer) close() error {
	if len(bs.connArr) <= 0 {
		return nil
	}

	for _, conn := range bs.connArr {
		conn.Close()
	}

	bs.connArr = nil
	return nil
}

// newMulticastConn creates a new UDP connection for a given port and list of
// multicast groups.
//
// It iterates through the list of network interfaces and for each interface that
// supports multicast and is running, it joins the multicast group for each IP
// address in the list.
//
// The function returns a *net.UDPConn and an error. If an error occurs while
// joining a multicast group, it will be returned. If no multicast groups are
// specified, the function will return nil as the error.
func newMulticastConn(port int, groups []net.IP) (*net.UDPConn, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	addrStat, err := appNet.StatAddr()
	if err != nil {
		return nil, err
	}

	var mAddr *net.UDPAddr
	if addrStat.IPv6Enabled {
		mAddr = &net.UDPAddr{IP: net.IPv6zero, Port: port}
	} else {
		mAddr = &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: port}
	}

	conn, err := net.ListenUDP("udp", mAddr)
	if err != nil {
		return nil, err
	}

	var ipv4Conn *ipv4.PacketConn
	var ipv6Conn *ipv6.PacketConn
	for _, iface := range ifaces {

		if net.FlagMulticast != (net.FlagMulticast & iface.Flags) {
			continue
		}
		if net.FlagRunning != (net.FlagRunning & iface.Flags) {
			continue
		}

		for _, group := range groups {
			var packetConn PacketConn
			var mIP net.IP
			if mIP = group.To4(); mIP != nil {
				if ipv4Conn == nil {
					ipv4Conn = ipv4.NewPacketConn(conn)
				}
				packetConn = ipv4Conn
			} else {
				if ipv6Conn == nil {
					ipv6Conn = ipv6.NewPacketConn(conn)
				}
				packetConn = ipv6Conn
				mIP = mAddr.IP
			}

			packetConn.JoinGroup(&iface, &net.UDPAddr{IP: mIP})
		}

	}

	return conn, nil
}
