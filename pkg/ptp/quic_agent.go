package ptp

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"pan/pkg/log"
	"path"
	"slices"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"google.golang.org/protobuf/proto"
)

var ErrQuicAgentUnavailable = errors.New("net.QuicAgent Error: Unavailable")
var ErrQuicAgentInvalidBroadcastPayload = errors.New("net.QuicAgent Error: Invalid Broadcast Payload")
var ErrQuicAgentInvalidBroadcastType = errors.New("net.QuicAgent Error: Invalid Broadcast Type")
var ErrQuicAgentServeBroadcastDenied = errors.New("net.QuicAgent Error: Serve Broadcast Denied")
var ErrQuicAgentDeliverExit = errors.New("net.QuicAgent Error: Deliver Exit")
var ErrQuicAgentNoUDPConn = errors.New("net.QuicAgent Error: No UDPConn")

const QUIC_AGENT_BROADCAST_DELIVER_TIMEOUT = 15 * time.Second
const QUIC_AGENT_OPTIMIZE_NETWORK_INTERVAL = 6 * time.Second

const (
	QuicBroadcastType = uint8(iota + 1)
)

var (
	QuicGreetRequestName = []byte("QuicGreet")
)

type stdQuicAgent struct {
	logger           log.Logger
	quicNetwork      *stdQuicNetwork
	broadcastNetwork *stdBroadcastNetwork
	server           *stdQuicServer
	guard            *stdPeerGuard

	broadcastCache *expirable.LRU[string, uint64]

	peerId           PeerID
	privateKey       crypto.PrivateKey
	broadcastEnabled bool

	reloadChan chan struct{}
	reloadLock sync.RWMutex
	reload     bool
}

func (agent *stdQuicAgent) setup(config QuicConfig) {
	agent.logger.Debug("net.QuicAgent", "Setup begin")
	defer agent.logger.Debug("net.QuicAgent", "Setup end")

	agent.reloadLock.Lock()
	defer agent.reloadLock.Unlock()

	changed := false
	peerId := config.PeerID()
	privateKey := config.PrivateKey()
	broadcastEnabled := config.BroadcastEnabled()

	if agent.broadcastEnabled != broadcastEnabled {
		agent.broadcastEnabled = broadcastEnabled
		changed = true
	}

	if !bytes.Equal(agent.peerId, peerId) {
		agent.peerId = peerId
		changed = true
	}

	if !equalPrivateKey(agent.privateKey, privateKey) {
		agent.privateKey = privateKey
		changed = true
	}

	if !changed || agent.reload {
		return
	}
	agent.reload = true
	agent.reloadChan <- struct{}{}
}

func (agent *stdQuicAgent) optimizeNetwork(ctx context.Context) error {
	agent.logger.Debug("net.QuicAgent", "OptimizeNetwork begin")
	defer agent.logger.Debug("net.QuicAgent", "OptimizeNetwork end")

	var err error
optimize_loop:
	for {
		agent.quicNetwork.optimize(ctx)

		select {
		case <-ctx.Done():
			err = ctx.Err()
			break optimize_loop
		case <-time.After(QUIC_AGENT_OPTIMIZE_NETWORK_INTERVAL):
			continue
		}
	}
	return err
}

var _ = (BroadcastServeModule)((*stdQuicAgent)(nil))

func (agent *stdQuicAgent) ServeBroadcast(payload []byte, addr string) error {
	if len(payload) <= 0 {
		return ErrQuicAgentInvalidBroadcastPayload
	}

	if payload[0] != QuicBroadcastType {
		return ErrQuicAgentInvalidBroadcastType
	}
	payload_ := payload[1:]

	agent.reloadLock.RLock()
	peerId := agent.peerId
	agent.reloadLock.RUnlock()

	quicNetwork := agent.quicNetwork
	broadcastCache := agent.broadcastCache
	if quicNetwork == nil || broadcastCache == nil || len(peerId) <= 0 {
		return ErrQuicAgentUnavailable
	}

	data, sig, err := unpackAgentPayload(payload_)
	if err != nil {
		return err
	}

	var online PeerOnline
	err = proto.Unmarshal(data, &online)
	if err == nil {
		if bytes.Equal(online.PeerId, peerId) {
			return err
		}
		err = Verify(data, sig, online.PeerId)
	}

	if err != nil {
		return err
	}

	cacheKey := path.Join(addr, EncodePeerID(peerId))
	cacheValue, ok := broadcastCache.Get(cacheKey)
	if ok && cacheValue >= online.Heightest {
		return nil
	}

	if ok && cacheValue < online.Heightest {
		broadcastCache.Remove(cacheKey)
	}
	if !ok || cacheValue < online.Heightest {
		broadcastCache.Add(cacheKey, online.Heightest)
	}

	if agent.guard.check(online.PeerId) != PeerGuardDeny {
		return ErrQuicAgentServeBroadcastDenied
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := quicNetwork.connectAddr(ctx, online.PeerId, addr)
	if err == nil {
		err = quicNetwork.route(conn)
		if err == nil {
			err = quicNetwork.reuse(conn)
		}
		if err != nil {
			conn.Close()
		}
	}

	return err
}

func (agent *stdQuicAgent) deliverBroadcast(ctx context.Context) error {
	agent.logger.Debug("net.QuicAgent", "DeliverBroadcast begin")
	defer agent.logger.Debug("net.QuicAgent", "DeliverBroadcast end")

	var err error

	var wg sync.WaitGroup
	var cancel context.CancelCauseFunc
	var enabled bool

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-agent.reloadChan:
			agent.reloadLock.Lock()
			agent.reload = false
			enabled = agent.broadcastEnabled
			agent.reloadLock.Unlock()
		}

		if cancel != nil {
			cancel(ErrQuicAgentDeliverExit)
			cancel = nil
			wg.Wait()
		}

		if err != nil {
			break
		}

		if !enabled {
			continue
		}

		causeCtx, causeCancel := context.WithCancelCause(ctx)
		cancel = causeCancel
		wg.Add(1)
		go func(causeCtx context.Context) {
			defer wg.Done()
		loop:
			for {
				err := agent.deliverOnline()
				if err != nil && err != ErrQuicAgentNoUDPConn {
					agent.logger.Error("net.QuicAgent", " deliverOnline Error: "+err.Error())
					break
				}

				select {
				case <-causeCtx.Done():
					break loop
				case <-time.After(QUIC_AGENT_BROADCAST_DELIVER_TIMEOUT):
				}
			}
		}(causeCtx)
	}
	return err
}

func (agent *stdQuicAgent) deliverOnline(deliverAddrs ...string) error {
	agent.reloadLock.RLock()
	peerId := agent.peerId
	privateKey := agent.privateKey
	agent.reloadLock.RUnlock()

	broadcastNetwork := agent.broadcastNetwork
	server := agent.server
	if server == nil || broadcastNetwork == nil || len(peerId) <= 0 || privateKey == nil {
		return ErrQuicAgentUnavailable
	}

	connList := server.TransportConnList()
	if len(connList) <= 0 {
		return ErrQuicAgentNoUDPConn
	}

	// deliver online
	var msg PeerOnline
	msg.PeerId = peerId
	msg.Heightest = uint64(time.Now().Unix())
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	sig, err := Sign(data, privateKey)
	if err != nil {
		return err
	}
	payload := packAgentPayload(data, sig)
	payload = slices.Insert(payload, 0, QuicBroadcastType)

	return broadcastNetwork.Deliver(connList, payload, deliverAddrs...)
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
		return nil, nil, ErrQuicAgentInvalidBroadcastPayload
	}
	offset := 2 + usize
	sig := buffer[2:offset]
	payload := buffer[offset:]
	return payload, sig, nil
}

func equalPrivateKey(target crypto.PrivateKey, source crypto.PrivateKey) bool {
	if target == nil && source == nil {
		return true
	}
	if target != nil && source != nil {

		if ecdsaPrivKey, ok := target.(*ecdsa.PrivateKey); ok {
			if ecdsaSource, ok := source.(*ecdsa.PrivateKey); ok {
				return ecdsaPrivKey.Equal(ecdsaSource)
			}
			return false
		}

	}
	return false
}
