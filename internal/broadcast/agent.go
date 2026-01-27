// Define broadcast agent for broadcast
//
// The broadcast agent is used to serve and deliver broadcast messages.
package broadcast

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"pan/internal/log"
	"pan/internal/peer"
	"slices"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"google.golang.org/protobuf/proto"
)

var ErrBroadcastAgentDeliverExit = errors.New("broadcast.BroadcastAgent Error: Deliver Exit")
var ErrBroadcastAgentPayloadInvalid = errors.New("broadcast.BroadcastAgent Error: Payload Invalid")
var ErrBroadcastAgentUnavailable = errors.New("broadcast.BroadcastAgent Error: Unavailable")

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
	logger   log.Logger
	network  *stdBroadcastNetwork
	provider *stdBroadcastProvider
	cache    *expirable.LRU[string, uint64]

	peerId     peer.PeerID
	privateKey crypto.PrivateKey

	reloadChan chan struct{}
	reloadLock sync.RWMutex
	reload     bool
}

func (agent *stdBroadcastAgent) Setup(config BroadcastConfig) {
	agent.logger.Debug("BroadcastAgent", "Setup")

	agent.reloadLock.Lock()
	defer agent.reloadLock.Unlock()

	changed := false
	peerId := config.PeerID()
	privateKey := config.PrivateKey()

	if !bytes.Equal(agent.peerId, peerId) {
		agent.peerId = peerId
		changed = true
	}

	if !equalPrivateKey(privateKey, agent.privateKey) {
		agent.privateKey = privateKey
	}

	if !changed || agent.reload {
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

// DeliverOnline sends the online message to the peers.
//
// The message is sent to the peers whose addresses are in the parameter list.
// If the parameter list is empty, the message is sent to all online peers.
//
// Returns an error if the broadcast agent is unavailable or if there is an issue delivering the online message.
func (agent *stdBroadcastAgent) DeliverOnline(deliverAddrs ...string) error {
	agent.reloadLock.RLock()
	peerId := agent.peerId
	privateKey := agent.privateKey
	agent.reloadLock.RUnlock()

	network := agent.network
	if network == nil || len(peerId) <= 0 || privateKey == nil {
		return ErrBroadcastAgentUnavailable
	}

	// deliver online
	var msg PeerOnline
	msg.PeerId = peerId
	msg.Heightest = uint64(time.Now().Unix())
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	sig, err := peer.Sign(data, privateKey)
	if err != nil {
		return err
	}
	payload := packAgentPayload(data, sig)
	payload = slices.Insert(payload, 0, BroadcastTypeOnline)

	return network.Deliver(payload, deliverAddrs...)
}

// AcceptOnline accepts an online message from a peer.
//
// The message is verified and the peer is routed to the quic peer module.
//
// Returns an error if the message is invalid or if there is an issue routing the peer.
func (agent *stdBroadcastAgent) AcceptOnline(payload []byte, addr string) error {
	provider := agent.provider
	network := agent.network
	if provider == nil || network == nil {
		return ErrBroadcastAgentUnavailable
	}

	quicNetwork := provider.QuicNetwork()
	if quicNetwork == nil {
		return ErrBroadcastAgentUnavailable
	}

	agent.reloadLock.RLock()
	peerId := agent.peerId
	privateKey := agent.privateKey
	agent.reloadLock.RUnlock()

	cache := agent.cache
	if cache == nil || len(peerId) <= 0 || privateKey == nil {
		return ErrBroadcastAgentUnavailable
	}

	data, sig, err := unpackAgentPayload(payload)
	if err != nil {
		return err
	}

	var msg PeerOnline
	err = proto.Unmarshal(data, &msg)
	if err == nil {
		if bytes.Equal(msg.PeerId, peerId) {
			return err
		}
		err = peer.Verify(data, sig, msg.PeerId)
	}

	if err != nil {
		return err
	}

	cacheKey := peer.EncodePeerID(msg.PeerId)
	cacheValue, ok := cache.Get(cacheKey)
	if ok && cacheValue >= msg.Heightest {
		return nil
	}

	if ok && cacheValue < msg.Heightest {
		cache.Remove(cacheKey)
	}
	if !ok || cacheValue < msg.Heightest {
		cache.Add(cacheKey, msg.Heightest)
	}

	return quicNetwork.Route(msg.PeerId, addr)
}

func (agent *stdBroadcastAgent) Deliver(ctx context.Context) error {
	agent.logger.Debug("BroadcastAgent", "Deliver begin")
	defer agent.logger.Debug("BroadcastAgent", "Deliver end")

	var err error

	var wg sync.WaitGroup
	var cancel context.CancelCauseFunc

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
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

		if err != nil {
			break
		}

		causeCtx, causeCancel := context.WithCancelCause(ctx)
		cancel = causeCancel
		wg.Add(1)
		go func(causeCtx context.Context) {
			defer wg.Done()
		loop:
			for {
				err := agent.DeliverOnline()
				if err != nil {
					agent.logger.Error("BroadcastAgent", " DeliverOnline Error: "+err.Error())
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
