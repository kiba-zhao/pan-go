package peer

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"

	"net"

	"github.com/quic-go/quic-go"
)

type quicNodeServeSt struct {
	conn     *net.UDPConn
	listener *quic.Listener
}

// Accept ...
func (ns *quicNodeServeSt) Accept(ctx context.Context) (Node, error) {
	conn, err := ns.listener.Accept(ctx)
	if err != nil {
		return nil, err
	}
	return &quicNodeSt{conn: conn}, nil
}

// Close ...
func (ns *quicNodeServeSt) Close() error {
	err := ns.listener.Close()
	if err != nil {
		return err
	}
	err = ns.conn.Close()
	return err
}

type quicNodeSt struct {
	conn quic.Connection
}

// Type ...
func (n *quicNodeSt) Type() NodeType {
	return QUICNodeType
}

// Addr ...
func (n *quicNodeSt) Addr() []byte {
	addr := n.conn.RemoteAddr()
	udpAddr, _ := net.ResolveUDPAddr(addr.Network(), addr.String())
	return MarshalQUICAddr(udpAddr)
}

// Certificate ...
func (n *quicNodeSt) Certificate() *x509.Certificate {
	state := n.conn.ConnectionState()
	return state.TLS.PeerCertificates[0]
}

func (n *quicNodeSt) AcceptNodeStream(ctx context.Context) (NodeStream, error) {
	conn := n.conn
	quicStream, err := conn.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}

	stream := new(quicNodeStreamSt)
	stream.Stream = quicStream
	return stream, nil
}

func (n *quicNodeSt) OpenNodeStream() (NodeStream, error) {

	conn := n.conn
	quicStream, err := conn.OpenStream()
	if err != nil {
		return nil, err
	}

	stream := new(quicNodeStreamSt)
	stream.Stream = quicStream
	return stream, nil
}

func (n *quicNodeSt) Close() error {
	conn := n.conn
	return conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), quic.NoError.Message())
}

type quicNodeStreamSt struct {
	quic.Stream
}

// CloseWrite ...
func (ns *quicNodeStreamSt) CloseWrite() error {
	ns.Stream.CancelWrite(quic.StreamErrorCode(quic.NoError))
	return nil
}

// CloseRead ...
func (ns *quicNodeStreamSt) CloseRead() error {
	ns.Stream.CancelRead(quic.StreamErrorCode(quic.NoError))
	return nil
}

type quicNodeDialerSt struct {
	tls       *tls.Config
	ctx       context.Context
	handshake []byte
}

// Type ...
func (nd *quicNodeDialerSt) Type() NodeType {
	return QUICNodeType
}

// Connect ...
func (nd *quicNodeDialerSt) Connect(addr []byte) (node Node, err error) {
	addrs := bytes.Split(addr, []byte(","))
	errs := make([]error, 0)
	for _, addr := range addrs {
		quicAddr, quicErr := UnmarshalQUICAddr(addr)
		if quicErr == nil {
			node, quicErr = DialQUICNode(quicAddr, nd.tls, nd.ctx)
		}

		if quicErr != nil {
			errs = append(errs, quicErr)
			continue
		}
		break
	}

	if len(errs) > 0 {
		err = errors.Join(errs...)
	}
	return
}

// Handshake ...
func (nd *quicNodeDialerSt) Handshake() []byte {
	return nd.handshake
}

// NewNodeDialer ...
func NewNodeDialer(tls *tls.Config, ctx context.Context, addrs ...[]byte) NodeDialer {
	dialer := new(quicNodeDialerSt)
	dialer.tls = tls
	dialer.ctx = ctx
	if len(addrs) > 0 {
		dialer.handshake = bytes.Join(addrs, []byte(","))
	}
	return dialer
}

// ServeQUICNode ...
func ServeQUICNode(addr *net.UDPAddr, tls *tls.Config) (NodeServe, error) {
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	quicConf := &quic.Config{}
	ln, err := quic.Listen(udpConn, tls, quicConf)
	if err != nil {
		return nil, err
	}

	return &quicNodeServeSt{conn: udpConn, listener: ln}, err
}

func DialQUICNode(addr *net.UDPAddr, tls *tls.Config, ctx context.Context) (Node, error) {
	quicConf := &quic.Config{}
	conn, err := quic.DialAddr(ctx, addr.String(), tls, quicConf)
	if err != nil {
		return nil, err
	}

	return &quicNodeSt{conn: conn}, err
}

// MarshalAddr ...
func MarshalQUICAddr(addr *net.UDPAddr) []byte {
	addrstr := addr.String()
	return []byte(addrstr)
}

// UnmarshalAddr ...
func UnmarshalQUICAddr(payload []byte) (addr *net.UDPAddr, err error) {
	addrstr := string(payload)
	addr, err = net.ResolveUDPAddr("udp", addrstr)
	return
}
