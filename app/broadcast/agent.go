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
	"pan/app/discovery"
	"pan/app/injection"
	"pan/app/peer"
	"pan/app/quic"
	"pan/logger"
	"pan/runtime"
	"reflect"
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

type broadcastAgent struct {
	PeerModule     peer.PeerModule
	QuicPeerModule quic.QuicPeerModule
	Store          BroadcastStore

	provider injection.ComponentStoreProvider
	module   *broadcastModule

	registry   runtime.Registry
	registryRW sync.RWMutex

	reloadLocker sync.Mutex
	reloadChan   chan struct{}
	reloadOnce   sync.Once
	needReload   bool
}

// ComponentStore returns the ComponentStore associated with the broadcast agent.
//
// This method retrieves the ComponentStore from the agent's provider,
// allowing access to the map of components that are injected into other components.
func (agent *broadcastAgent) ComponentStore() injection.ComponentStore {
	return agent.provider.ComponentStore()
}

func (agent *broadcastAgent) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(agent, injection.ComponentNoneScope),
		injection.NewComponent[discovery.Broadcast](agent, injection.ComponentExternalScope),
	}
}

func (agent *broadcastAgent) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[BroadcastPublicAddrsProvider](),
	}
}

// PublicAddrs returns the list of public addresses that the broadcast agent can use.
//
// The addresses are used to broadcast messages publicly.
//
// The addresses are retrieved from the registry, and the list is deduplicated
// before being returned.
func (agent *broadcastAgent) PublicAddrs() []string {
	agent.registryRW.RLock()
	registry := agent.registry
	agent.registryRW.RUnlock()

	if registry == nil {
		return nil
	}

	addrs := make([]string, 0)
	runtime.TraverseRegistry(registry, func(provider BroadcastPublicAddrsProvider) error {
		publicAddrs := provider.BroadcastPublicAddrs()
		if len(publicAddrs) <= 0 {
			return nil
		}
		for _, addr := range publicAddrs {
			if idx, ok := slices.BinarySearch(addrs, addr); !ok {
				addrs = slices.Insert(addrs, idx, addr)
			}
		}
		return nil
	})
	return addrs
}

func (agent *broadcastAgent) Init(registry runtime.Registry) error {
	agent.registryRW.Lock()
	defer agent.registryRW.Unlock()
	agent.registry = registry
	return nil
}

func (agent *broadcastAgent) ReloadChan() chan struct{} {

	agent.reloadOnce.Do(func() {
		agent.reloadChan = make(chan struct{}, 1)
	})

	return agent.reloadChan
}

// Reload reloads the broadcast agent.
//
// It is called by the runtime to reload the broadcast agent after the application has finished initializing.
//
// The function first checks if the registry is available, and if it is not, an error is returned.
// If the registry is available, the function calls ReloadModules to reload the broadcast agent.
//
// ReloadModules is a noop if the registry is not available.
func (agent *broadcastAgent) Reload() {
	agent.reloadLocker.Lock()
	defer agent.reloadLocker.Unlock()
	if agent.needReload {
		return
	}
	agent.needReload = true
	agent.ReloadChan() <- struct{}{}
}

// ServeBroadcast serves the broadcast message to the peer.
// Returns an error if the broadcast agent is unavailable or if there is an issue serving the broadcast.
func (agent *broadcastAgent) ServeBroadcast(payload []byte, addr string) error {
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

// DeliverOnline sends the online message to the peers.
//
// The message is sent to the peers whose addresses are in the parameter list.
// If the parameter list is empty, the message is sent to all online peers.
//
// Returns an error if the broadcast agent is unavailable or if there is an issue delivering the online message.
func (agent *broadcastAgent) DeliverOnline(deliverAddrs ...string) error {

	settings := agent.PeerModule.PeerSettings()
	if !settings.Available() {
		return ErrBroadcastAgentDeliverOnlineUnavailable
	}

	if len(deliverAddrs) <= 0 {
		deliverAddrs = agent.module.DeliverAddrs()
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
	info, err := agent.Store.SelectOrCreate(info)
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
		sig, err := peer.Sign(data, settings.PrivKey())
		if err != nil {
			continue
		}
		payload := packAgentPayload(data, sig)
		payload = slices.Insert(payload, 0, BroadcastTypeOnline)
		err = agent.module.Deliver(payload, addrs...)
		if err != nil {
			logger.Default().Log(context.Background(), logger.LevelError, "broadcast agent deliver online failed: "+err.Error())
		}
	}

	return nil
}

// AcceptOnline accepts an online message from a peer.
//
// The message is verified and the peer is routed to the quic peer module.
//
// Returns an error if the message is invalid or if there is an issue routing the peer.
func (agent *broadcastAgent) AcceptOnline(payload []byte, addr string) error {
	settings := agent.PeerModule.PeerSettings()
	if !settings.Available() {
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

	err = agent.Store.SaveHighest(msg.PeerId, msg.Heightest)
	if err != nil {
		return err
	}

	return agent.QuicPeerModule.Route(msg.PeerId, msgAddr, true)
}

// Ready prepares the broadcast agent to operate within the given context.
//
// The method manages the lifecycle of the broadcast agent, handling context cancellation,
// reload signals, and delivery of online messages to peers. It stops any ongoing delivery
// when the context is done or when a reload is triggered and restarts it if necessary.
//
// The method initializes the broadcast store with the peer ID and ensures that the peer
// settings are available before proceeding. It runs the StartDelivery method in a
// separate goroutine, which continues to deliver online messages until the context is canceled.
//
// Returns an error if the context is canceled or if there is an issue initializing the store.

func (agent *broadcastAgent) Ready(ctx context.Context) error {

	var cancel context.CancelCauseFunc
	var err error
	closed := false

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-agent.ReloadChan():
		}

		agent.reloadLocker.Lock()
		agent.needReload = false
		agent.reloadLocker.Unlock()

		agent.StopDelivery(cancel)
		if closed {
			break
		}

		settings := agent.PeerModule.PeerSettings()
		if !settings.Available() {
			continue
		}

		err := agent.Store.Init(settings.PeerID())
		if err != nil {
			logger.Default().Log(context.Background(), logger.LevelError, "broadcast agent init failed: %s", err.Error())
		}

		causeCtx, causeCancel := context.WithCancelCause(ctx)
		cancel = causeCancel
		go agent.StartDelivery(causeCtx)
	}

	return err
}

// StartDelivery starts delivering online messages to peers.
//
// The method is designed to be run in its own goroutine and will
// continue to deliver online messages until the context is canceled.
//
// The method will wait for 30 seconds between each delivery, unless
// the context is canceled before the next delivery can be made.
//
// If an error occurs while delivering the online message, the method
// will stop and return the error.
func (agent *broadcastAgent) StartDelivery(ctx context.Context) {
loop:
	for {
		err := agent.DeliverOnline()
		if err != nil {
			break
		}

		select {
		case <-ctx.Done():
			break loop
		case <-time.After(30 * time.Second):
		}
	}
}

// StopDelivery stops the delivery of online messages to peers.
//
// The method cancels the context that was passed to the StartDelivery method.
// If the context is already canceled, the method does nothing.
//
// The method is safe to call multiple times.
func (agent *broadcastAgent) StopDelivery(cancel context.CancelCauseFunc) {
	if cancel != nil {
		cancel(ErrBroadcastAgentDeliverExit)
	}
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
