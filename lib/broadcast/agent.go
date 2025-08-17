// Define broadcast agent for broadcast
//
// The broadcast agent is used to serve and deliver broadcast messages.
package broadcast

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
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

type BroadcastAgentRuntime interface {
	QuicCluster() quic.QuicCluster
	PeerSettings() *peer.PeerSettings
	Store() BroadcastStore
}

type stdBroadcastAgent struct {
	logger log.Logger

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool

	runtime BroadcastAgentRuntime
	cluster BroadcastCluster
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

// DeliverOnline sends the online message to the peers.
//
// The message is sent to the peers whose addresses are in the parameter list.
// If the parameter list is empty, the message is sent to all online peers.
//
// Returns an error if the broadcast agent is unavailable or if there is an issue delivering the online message.
func (agent *stdBroadcastAgent) DeliverOnline(deliverAddrs ...string) error {

	agentRuntime := agent.runtime
	if agentRuntime == nil {
		return ErrBroadcastAgentUnavailable
	}
	store := agentRuntime.Store()
	cluster := agent.cluster
	settings := agentRuntime.PeerSettings()
	if settings == nil || cluster == nil || store == nil {
		return ErrBroadcastAgentUnavailable
	}

	info := BroadcastInfo{}
	info.PeerID = settings.PeerID()
	info.Heightest = 0
	info.UpdatedAt = time.Now()
	info, err := store.SelectOrCreate(info)
	if err != nil {
		return err
	}

	// deliver online
	var msg PeerOnline
	msg.PeerId = settings.PeerID()
	msg.Heightest = info.Heightest
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	sig, err := peer.Sign(data, settings.PrivateKey())
	if err != nil {
		return err
	}
	payload := packAgentPayload(data, sig)
	payload = slices.Insert(payload, 0, BroadcastTypeOnline)

	return cluster.Deliver(payload, deliverAddrs...)
}

// AcceptOnline accepts an online message from a peer.
//
// The message is verified and the peer is routed to the quic peer module.
//
// Returns an error if the message is invalid or if there is an issue routing the peer.
func (agent *stdBroadcastAgent) AcceptOnline(payload []byte, addr string) error {
	agentRuntime := agent.runtime
	if agentRuntime == nil {
		return ErrBroadcastAgentUnavailable
	}

	quicCluster := agentRuntime.QuicCluster()
	store := agentRuntime.Store()
	settings := agentRuntime.PeerSettings()
	if settings == nil || store == nil || quicCluster == nil {
		return ErrBroadcastAgentUnavailable
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
	if err == nil {
		err = store.SaveHighest(msg.PeerId, msg.Heightest)
	}
	if err != nil {
		return err
	}

	return quicCluster.Route(msg.PeerId, addr)
}

func (agent *stdBroadcastAgent) Deliver(ctx context.Context) error {
	agent.logger.Debug("BroadcastAgent", "Deliver begin")
	defer agent.logger.Debug("BroadcastAgent", "Deliver end")

	var err error
	var closed bool

	var wg sync.WaitGroup
	var cancel context.CancelCauseFunc
	var timer <-chan time.Time

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-agent.reloadChan:
			if timer != nil {
				<-timer
				timer = nil
			}
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

		agentRuntime := agent.runtime
		if ctx == nil {
			timer = time.After(time.Second * 5)
			agent.Reload()
			continue
		}

		peerSettings := agentRuntime.PeerSettings()
		store := agentRuntime.Store()
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
					agent.logger.Error("BroadcastAgent", " DeliverOnline Error: "+err.Error())
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
