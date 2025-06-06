// Define broadcast agent for broadcast
//
// The broadcast agent is used to serve and deliver broadcast messages.
package broadcast

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"net"
	"pan/lib/log"
	"pan/lib/peer"
	"pan/lib/quic"
	"slices"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

var ErrBroadcastAgentDeliverExit = errors.New("broadcast.BroadcastAgent Error: Deliver Exit")
var ErrBroadcastAgentPayloadInvalid = errors.New("broadcast.BroadcastAgent Error: Payload Invalid")
var ErrBroadcastAgentDeliverOnlineUnavailable = errors.New("broadcast.BroadcastAgent Error: Deliver Online Unavailable")

const (
	BroadcastTypeOnline = uint8(iota + 1)
)

const (
	BroadcastMulticastTypeGlobal = uint8(iota + 1)
	BroadcastMulticastTypeIPV6Global
	BroadcastMulticastTypeLocal
	BroadcastMulticastTypeIPV6Local
)

type stdBroadcastAgent struct {
	logger log.Logger

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool

	publicAddrs   []string
	publicAddrsRW sync.RWMutex

	peerSettings   *peer.PeerSettings
	peerSettingsRW sync.RWMutex

	store   BroadcastStore
	storeRW sync.RWMutex

	quicCluster   quic.QuicCluster
	quicClusterRW sync.RWMutex

	cluster *stdBroadcastCluster
}

// PublicAddrs returns the list of public addresses that the broadcast agent can use.
//
// The addresses are used to broadcast messages publicly.
//
// The addresses are retrieved from the registry, and the list is deduplicated
// before being returned.
func (agent *stdBroadcastAgent) PublicAddrs() []string {
	agent.publicAddrsRW.RLock()
	defer agent.publicAddrsRW.RUnlock()
	return agent.publicAddrs
}

func (agent *stdBroadcastAgent) SetPublicAddrs(addrs []string) {
	agent.publicAddrsRW.Lock()
	defer agent.publicAddrsRW.Unlock()
	if slices.Equal(agent.publicAddrs, addrs) {
		return
	}
	agent.publicAddrs = addrs
	agent.Reload()
}

func (agent *stdBroadcastAgent) Reload() {
	agent.logger.Debug("BroadcastAgent", "Reload")

	agent.reloadLock.Lock()
	defer agent.reloadLock.Unlock()
	if agent.reload {
		return
	}

	agent.reload = true
	agent.reloadChan <- struct{}{}
}

// ServeBroadcast serves the broadcast message to the peer.
// Returns an error if the broadcast agent is unavailable or if there is an issue serving the broadcast.
func (agent *stdBroadcastAgent) ServeBroadcast(payload []byte, addr string) error {
	if len(payload) <= 0 {
		return nil
	}

	var err error
	switch flag := payload[0]; flag {
	case BroadcastTypeOnline:
		err = agent.AcceptOnline(slices.Clone(payload[1:]), addr)
	}
	return err
}

func (agent *stdBroadcastAgent) PeerSettings() *peer.PeerSettings {
	agent.peerSettingsRW.RLock()
	defer agent.peerSettingsRW.RUnlock()
	return agent.peerSettings
}

func (agent *stdBroadcastAgent) SetPeerSettings(settings *peer.PeerSettings) {
	agent.peerSettingsRW.Lock()
	defer agent.peerSettingsRW.Unlock()
	agent.peerSettings = settings

	agent.Reload()
}

func (agent *stdBroadcastAgent) Store() BroadcastStore {
	agent.storeRW.RLock()
	defer agent.storeRW.RUnlock()
	return agent.store
}

func (agent *stdBroadcastAgent) SetStore(store BroadcastStore) {
	agent.storeRW.Lock()
	defer agent.storeRW.Unlock()
	agent.store = store

	agent.Reload()
}

func (agent *stdBroadcastAgent) QuicCluster() quic.QuicCluster {
	agent.quicClusterRW.RLock()
	defer agent.quicClusterRW.RUnlock()
	return agent.quicCluster
}

func (agent *stdBroadcastAgent) SetQuicCluster(cluster quic.QuicCluster) {
	agent.quicClusterRW.Lock()
	defer agent.quicClusterRW.Unlock()
	agent.quicCluster = cluster
}

// DeliverOnline sends the online message to the peers.
//
// The message is sent to the peers whose addresses are in the parameter list.
// If the parameter list is empty, the message is sent to all online peers.
//
// Returns an error if the broadcast agent is unavailable or if there is an issue delivering the online message.
func (agent *stdBroadcastAgent) DeliverOnline(deliverAddrs ...string) error {

	store := agent.Store()
	cluster := agent.cluster
	settings := agent.PeerSettings()
	if settings == nil || cluster == nil || store == nil {
		return ErrBroadcastAgentDeliverOnlineUnavailable
	}

	if len(deliverAddrs) <= 0 {
		deliverAddrs = cluster.DeliverAddrs()
	}

	if len(deliverAddrs) <= 0 {
		return nil
	}

	publicAddrs := agent.PublicAddrs()
	if len(publicAddrs) <= 0 {
		return nil
	}

	info := BroadcastInfo{}
	info.PeerID = settings.PeerID()
	info.Heightest = 0
	info.UpdatedAt = time.Now()
	info, err := store.SelectOrCreate(info)
	if err != nil {
		return err
	}

	//
	addrsMap := make(map[uint8][]string)
	for _, addr := range deliverAddrs {
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			continue
		}
		multicastType := selectBroadcastMulticastType(udpAddr.IP)
		if multicastType != 0 {
			addrsMap[multicastType] = append(addrsMap[multicastType], addr)
		}
	}

	// deliver with public addrs
	for _, publicAddr := range publicAddrs {
		ip, _, err := net.SplitHostPort(publicAddr)
		if err != nil {
			continue
		}
		ipAddr, err := net.ResolveIPAddr("ip", ip)
		if err != nil || ipAddr.IP.IsMulticast() {
			continue
		}

		var addrs []string

		if ipAddr.IP.IsUnspecified() {
			addrs = deliverAddrs
		} else {
			addrsList := make([][]string, 0)
			if ipAddr.IP.IsLinkLocalUnicast() || ipAddr.IP.IsPrivate() {
				if ipAddr.IP.To4() != nil {
					addrsList = append(addrsList, addrsMap[BroadcastMulticastTypeLocal])
				} else {
					addrsList = append(addrsList, addrsMap[BroadcastMulticastTypeIPV6Local])
				}
			}
			if ipAddr.IP.IsGlobalUnicast() && !ipAddr.IP.IsPrivate() {
				if ipAddr.IP.To4() != nil {
					addrsList = append(addrsList, addrsMap[BroadcastMulticastTypeGlobal])
				} else {
					addrsList = append(addrsList, addrsMap[BroadcastMulticastTypeIPV6Global])
				}
			}
			addrs = slices.Concat(addrsList...)
		}

		if len(addrs) <= 0 {
			continue
		}

		// deliver online
		var msg PeerOnline
		msg.PeerId = settings.PeerID()
		msg.Addr = publicAddr
		msg.Heightest = info.Heightest
		data, err := proto.Marshal(&msg)
		if err != nil {
			continue
		}
		sig, err := peer.Sign(data, settings.PrivateKey())
		if err != nil {
			continue
		}
		payload := packAgentPayload(data, sig)
		payload = slices.Insert(payload, 0, BroadcastTypeOnline)
		err = cluster.Deliver(payload, addrs...)
		if err != nil {
			log.Default().Error("BroadcastAgent", "Deliver Error: "+err.Error())
		}
	}

	return nil
}

// AcceptOnline accepts an online message from a peer.
//
// The message is verified and the peer is routed to the quic peer module.
//
// Returns an error if the message is invalid or if there is an issue routing the peer.
func (agent *stdBroadcastAgent) AcceptOnline(payload []byte, addr string) error {
	quicCluster := agent.QuicCluster()
	store := agent.Store()
	settings := agent.PeerSettings()
	if settings == nil || store == nil || quicCluster == nil {
		return ErrBroadcastAgentDeliverOnlineUnavailable
	}

	data, sig, err := unpackAgentPayload(payload)
	if err != nil {
		return err
	}

	var msg PeerOnline
	err = proto.Unmarshal(data, &msg)
	if err == nil {
		if bytes.Equal(msg.PeerId, settings.PeerID()) {
			return err
		}
		err = peer.Verify(data, sig, msg.PeerId)
	}
	if err != nil {
		return err
	}

	addrIP, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	msgAddr := msg.Addr
	host, port, err := net.SplitHostPort(msgAddr)
	if err != nil {
		return err
	}
	if addrIP != host {
		ipAddr, err := net.ResolveIPAddr("ip", host)
		if err != nil || !ipAddr.IP.IsUnspecified() {
			return err
		}
		msgAddr = net.JoinHostPort(addrIP, port)
	}

	err = store.SaveHighest(msg.PeerId, msg.Heightest)
	if err != nil {
		return err
	}

	return quicCluster.Route(msg.PeerId, msgAddr)
}

func (agent *stdBroadcastAgent) Deliver(ctx context.Context) error {
	agent.logger.Debug("BroadcastAgent", "Deliver begin")
	defer agent.logger.Debug("BroadcastAgent", "Deliver end")

	var err error
	var closed bool

	var wg sync.WaitGroup
	var cancel context.CancelCauseFunc

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-agent.reloadChan:
			agent.reloadLock.Lock()
			agent.reload = false
			agent.reloadLock.Unlock()
		}

		if cancel != nil {
			cancel(ErrBroadcastAgentDeliverExit)
			cancel = nil
			wg.Wait()
		}

		if closed {
			break
		}

		peerSettings := agent.PeerSettings()
		store := agent.Store()
		if peerSettings == nil || store == nil {
			continue
		}

		err := store.Init(peerSettings.PeerID())
		if err != nil {
			agent.logger.Error("BroadcastAgent", " Init Store Error: "+err.Error())
			continue
		}

		causeCtx, causeCancel := context.WithCancelCause(ctx)
		cancel = causeCancel
		go func(causeCtx context.Context) {
		loop:
			for {
				err := agent.DeliverOnline()
				if err != nil {
					break
				}

				select {
				case <-causeCtx.Done():
					break loop
				case <-time.After(30 * time.Second):
				}
			}
		}(causeCtx)
	}
	return err

}

func packAgentPayload(payload []byte, sig []byte) []byte {
	sigSizeBuffer := make([]byte, 2)
	packBuffer := bytes.Join([][]byte{sigSizeBuffer, sig, payload}, nil)
	binary.BigEndian.PutUint16(packBuffer[:2], uint16(len(sig)))
	return packBuffer
}
func unpackAgentPayload(buffer []byte) ([]byte, []byte, error) {
	usize := binary.BigEndian.Uint16(buffer[:2])
	if usize+2 > uint16(len(buffer)) {
		return nil, nil, ErrBroadcastAgentPayloadInvalid
	}
	offset := 2 + usize
	sig := buffer[2:offset]
	payload := buffer[offset:]
	return payload, sig, nil
}

func selectBroadcastMulticastType(ip net.IP) uint8 {

	var isGlobal bool
	if ip.IsMulticast() {
		isGlobal = isGlobalMulticastIP(ip)
	} else {
		isGlobal = ip.IsGlobalUnicast() && !ip.IsPrivate()
	}

	var multicastType uint8
	if isGlobal {
		if ip.To16() == nil {
			multicastType = BroadcastMulticastTypeGlobal
		} else {
			multicastType = BroadcastMulticastTypeIPV6Global
		}
	} else {
		if ip.To4() != nil {
			multicastType = BroadcastMulticastTypeLocal
		} else {
			multicastType = BroadcastMulticastTypeIPV6Local
		}
	}
	return multicastType
}
