package net

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"pan/internal/log"
	"slices"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"google.golang.org/protobuf/proto"
)

var ErrQuicAgentUnavailable = errors.New("net.QuicAgent Error: Unavailable")
var ErrQuicAgentInvalidBroadcastPayload = errors.New("net.QuicAgent Error: Invalid Broadcast Payload")
var ErrQuicAgentInvalidBroadcastType = errors.New("net.QuicAgent Error: Invalid Broadcast Type")
var ErrQuicAgentGreetForbidden = errors.New("net.QuicAgent Error: Greet Forbidden")
var ErrQuicAgentDeliverExit = errors.New("net.QuicAgent Error: Deliver Exit")
var ErrQuicAgentNoUDPConn = errors.New("net.QuicAgent Error: No UDPConn")

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

	cache *expirable.LRU[string, uint64]

	peerId           PeerID
	privateKey       crypto.PrivateKey
	broadcastEnabled bool

	guards  []PeerGuard
	guardRW sync.RWMutex

	reloadChan chan struct{}
	reloadLock sync.RWMutex
	reload     bool
}

var _ = (PeerTopic)((*stdQuicAgent)(nil))

func (agent *stdQuicAgent) SetupToPeer(router PeerServletRouter) error {
	router.Handle(QuicGreetRequestName, agent.HandleGreet)
	return nil
}

func (agent *stdQuicAgent) Setup(config QuicConfig) {
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
	cache := agent.cache
	if quicNetwork == nil || cache == nil || len(peerId) <= 0 {
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

	cacheKey := EncodePeerID(online.PeerId)
	cacheValue, ok := cache.Get(cacheKey)
	if ok && cacheValue >= online.Heightest {
		return nil
	}

	if ok && cacheValue < online.Heightest {
		cache.Remove(cacheKey)
	}
	if !ok || cacheValue < online.Heightest {
		cache.Add(cacheKey, online.Heightest)
	}

	if quicNetwork.HasRoute(online.PeerId, addr) {
		return nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := quicNetwork.ConnectAddr(ctx, addr)
	if err == nil {
		err = agent.Greet(conn)
	}

	if err == nil {
		quicNetwork.Route(conn)
		quicNetwork.Reuse(conn)
	}
	return err
}

func (agent *stdQuicAgent) Greet(conn *stdQuicConn) error {
	// TODO: Implement greet logic
	reqReader := NewRequest(QuicGreetRequestName, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := conn.OpenStream(ctx)
	if err != nil {
		return err
	}
	defer stream.Close()

	_, reader, err := DoAction(ctx, stream, reqReader)
	if err != nil {
		return err
	}
	defer reader.Close()

	// TODO:
	return nil
}

func (agent *stdQuicAgent) HandleGreet(ctx PeerServletContext, next PeerServletNext) error {

	conn, ok := ctx.Session(PeerConnSessionKey)
	if !ok {
		return next()
	}
	quicConn, ok := conn.(*stdQuicConn)
	if !ok {
		return next()
	}

	peerId, ok := ctx.Session(PeerIDSessionKey)
	if !ok {
		return next()
	}
	peerId_, ok := peerId.(PeerID)
	if !ok {
		return next()
	}

	agent.guardRW.RLock()
	guards := agent.guards
	agent.guardRW.RUnlock()

	isAllow := true
	if len(guards) > 0 {
		for _, guard := range guards {
			isAllow = guard.AllowAccess(peerId_)
			if !isAllow {
				break
			}
		}
	}

	if !isAllow {
		return ErrQuicAgentGreetForbidden
	}

	// TODO:
	agent.quicNetwork.Route(quicConn)
	return nil
}

func (agent *stdQuicAgent) DeliverBroadcast(ctx context.Context) error {
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
				if err != nil {
					agent.logger.Error("net.QuicAgent", " deliverOnline Error: "+err.Error())
					break
				}

				select {
				case <-causeCtx.Done():
					break loop
				case <-time.After(15 * time.Second):
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
