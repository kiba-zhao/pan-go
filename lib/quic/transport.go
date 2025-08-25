package quic

import (
	"net"

	"github.com/quic-go/quic-go"
)

type QuicTransport struct {
	*quic.Transport
	ipNet *net.IPNet
}
