package quic

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"net"
	"pan/app/broadcast"
	"pan/app/peer"
	"time"
)

var ErrQuicPeerBroadcastUnavailable = errors.New("quic.PeerBroadcast Error: Unavailable")
var ErrQuicPeerBroadcastPeerSettingsUnavailable = errors.New("quic.PeerBroadcast Error: Peer Settings Unavailable")
var ErrQuicPeerBroadcastNoPublicAddress = errors.New("quic.PeerBroadcast Error: No Public Address")

type quicPeerBroadcast struct {
	Broadcast      broadcast.BroadcastModule
	quicPeerModule QuicPeerModule
}

func (qb *quicPeerBroadcast) ServeBroadcast(payload []byte, ip string) error {
	settings := qb.quicPeerModule.PeerSettings()
	if !settings.Available() {
		return ErrQuicPeerBroadcastUnavailable
	}

	payloadLen := len(payload)
	offset := 0
	nextOffset := offset + 2

	peerIdLen := int(binary.BigEndian.Uint16(payload[offset:nextOffset]))
	offset = nextOffset
	nextOffset += peerIdLen

	if nextOffset+2 > payloadLen {
		return nil
	}
	peerId := peer.PeerID(payload[offset:nextOffset])

	// ingore self by peer id
	if bytes.Equal(settings.PeerID(), peerId) {
		return nil
	}

	offset = nextOffset
	nextOffset += 2
	addressLen := int(binary.BigEndian.Uint16(payload[offset:nextOffset]))
	offset = nextOffset
	nextOffset += addressLen
	if nextOffset != payloadLen {
		return nil
	}

	address := string(payload[offset:nextOffset])
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	if host != ip {
		ipAddr, err := net.ResolveIPAddr("ip", host)
		if err != nil || !ipAddr.IP.IsUnspecified() {
			return err
		}
		address = net.JoinHostPort(ip, port)
	}

	// Try  route
	if err == nil {
		err = qb.quicPeerModule.Route(peerId, address)
	}

	return err
}

func (qb *quicPeerBroadcast) DeliverBroadcast() error {
	addrs := qb.quicPeerModule.PublicAddrs()
	addrsCount := len(addrs)
	if addrsCount <= 0 {
		return ErrQuicPeerBroadcastNoPublicAddress
	}

	settings := qb.quicPeerModule.PeerSettings()
	if !settings.Available() {
		return ErrQuicPeerBroadcastPeerSettingsUnavailable
	}

	peerId := settings.PeerID()
	peerIdLen := len(peerId)

	var buffer []byte
	var offset int
	errs := make([]error, 0)
	for _, address := range addrs {
		addressLen := len(address)
		bufferSize := 4 + peerIdLen + addressLen
		if len(buffer) != bufferSize {
			buffer = make([]byte, bufferSize)
			binary.BigEndian.PutUint16(buffer, uint16(len(peerId)))
			copy(buffer[2:], peerId)
			offset = 2 + peerIdLen
			binary.BigEndian.PutUint16(buffer[offset:], uint16(addressLen))
			offset += 2
		}
		copy(buffer[offset:], []byte(address))
		err := qb.Broadcast.Deliver(buffer)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) <= 0 {
		return nil
	}
	return errors.Join(errs...)
}

func (qb *quicPeerBroadcast) Ready(ctx context.Context) error {
	go func(ctx context.Context) {
	broadcastLoop:
		for {
			err := qb.DeliverBroadcast()
			if errors.Is(err, ErrQuicPeerBroadcastPeerSettingsUnavailable) || errors.Is(err, ErrQuicPeerBroadcastNoPublicAddress) {
				return
			}

			select {
			case <-ctx.Done():
				break broadcastLoop
			case <-time.After(15 * time.Second):
				continue
			}
		}
	}(ctx)
	return nil
}
