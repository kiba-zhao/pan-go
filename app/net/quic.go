package net

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"pan/app/config"
	"pan/app/constant"
	"pan/app/node"
	"pan/runtime"
	"slices"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

type quicNode struct {
	resourceId node.NodeResourceID
	nodeId     node.NodeID
	conn       quic.Connection
	quicModule QuicModule
	mgr        node.NodeManager
}

func (qn *quicNode) ID() node.NodeID {
	return qn.nodeId
}

func (qn *quicNode) Type() node.NodeType {
	return node.NodeTypeAlive
}

func (qn *quicNode) Do(ctx context.Context, reader io.Reader) (io.Reader, error) {
	return qn.quicModule.Do(ctx, qn.conn, reader)
}

func (qn *quicNode) Greet(ctx context.Context) error {
	return qn.quicModule.Greet(ctx, qn.conn)
}

func (qn *quicNode) Close() error {
	qn.mgr.Delete(qn)
	return qn.conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
}

func (qn *quicNode) ResourceID() node.NodeResourceID {
	return qn.resourceId
}

type quicRoute struct {
	resourceId    node.NodeResourceID
	nodeId        node.NodeID
	address       string
	quicModule    QuicModule
	mgr           node.NodeManager
	failures      uint8
	failureLocker sync.RWMutex
}

func (qr *quicRoute) ID() node.NodeID {
	return qr.nodeId
}

func (qr *quicRoute) Type() node.NodeType {
	return node.NodeTypeReachable
}

func (qr *quicRoute) Dial(ctx context.Context) (quic.Connection, error) {

	qr.failureLocker.RLock()
	if qr.failures >= 3 {
		qr.Close()
		return nil, constant.ErrInvalidNode
	}
	qr.failureLocker.RUnlock()

	conn, err := qr.quicModule.Dial(ctx, qr.address)
	if err == nil {
		nodeId, err := qr.quicModule.ParseNodeID(conn)
		if err == nil && !bytes.Equal(nodeId, qr.nodeId) {
			defer qr.Close()
			err = constant.ErrInvalidNode
		}
		if err != nil {
			conn = nil
			defer conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
		}
	}

	qr.failureLocker.Lock()
	defer qr.failureLocker.Unlock()
	if err != nil {
		qr.failures++
		failures := qr.failures
		if failures >= 3 {
			qr.Close()
		}
	} else {
		qr.failures = 0
	}

	return conn, err
}

func (qr *quicRoute) Do(ctx context.Context, reader io.Reader) (io.Reader, error) {
	conn, err := qr.Dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
	go qr.quicModule.Serve(conn)
	return qr.quicModule.Do(ctx, conn, reader)
}

func (qr *quicRoute) Greet(ctx context.Context) error {
	conn, err := qr.Dial(ctx)
	if err == nil {
		defer conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
		err = qr.quicModule.Greet(ctx, conn)
	}
	return err
}

func (qr *quicRoute) Close() error {
	qr.mgr.Delete(qr)
	return nil
}

func (qr *quicRoute) ResourceID() node.NodeResourceID {
	return qr.resourceId
}

const (
	QuicNodeStream = uint8(iota)
	QuicGreetStream
)

type quicServer struct {
	quicModule *quicModule
	locker     sync.RWMutex
	ln         *quic.Listener
	address    string
}

func (qs *quicServer) Shutdown() error {
	qs.locker.RLock()
	ln := qs.ln
	qs.locker.RUnlock()
	if ln == nil {
		return constant.ErrUnavailable
	}

	qs.locker.Lock()
	qs.ln = nil
	qs.locker.Unlock()
	return ln.Close()
}

func (qs *quicServer) ListenAndServe() error {

	if qs.quicModule == nil || qs.quicModule.NodeModule == nil {
		return constant.ErrUnavailable
	}

	nodeSettings := qs.quicModule.NodeModule.NodeSettings()
	if !nodeSettings.Available() {
		return constant.ErrUnavailable
	}

	certificate := nodeSettings.Certificate()
	tlsConf := &tls.Config{ClientAuth: tls.RequireAnyClientCert, Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
	quicConf := &quic.Config{}
	ln, err := quic.ListenAddr(qs.address, tlsConf, quicConf)
	if err != nil {
		return err
	}
	qs.locker.Lock()
	qs.ln = ln
	qs.locker.Unlock()
	defer qs.Shutdown()

	for {
		conn, err := ln.Accept(context.Background())
		if errors.Is(err, quic.ErrServerClosed) {
			break
		}
		if err != nil {
			go conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
			continue
		}
		go qs.quicModule.Serve(conn)
	}

	return err
}

type QuicModule interface {
	Serve(quic.Connection) error
	Do(context.Context, quic.Connection, io.Reader) (io.Reader, error)
	Greet(context.Context, quic.Connection) error
	Dial(context.Context, string) (quic.Connection, error)
	ParseNodeID(quic.Connection) (node.NodeID, error)
	CreateNode(quic.Connection) (node.Node, error)
	CreateRoute(node.NodeID, string) (node.Node, error)
}

type quicModule struct {
	Broadcast  Broadcast
	NodeModule node.NodeModule

	publicAddrs []string
	addrs       []string
	locker      sync.RWMutex
	sigChan     chan []bool
	sigOnce     sync.Once
	hasSig      bool
}

func (qm *quicModule) PublicAddrs() []string {
	qm.locker.RLock()
	defer qm.locker.RUnlock()
	return qm.publicAddrs
}

func (qm *quicModule) Addrs() []string {
	qm.locker.RLock()
	defer qm.locker.RUnlock()
	return qm.addrs
}

func (qm *quicModule) setSig(sig []bool) {
	if qm.hasSig {
		return
	}

	qm.sigOnce.Do(func() {
		qm.sigChan = make(chan []bool, 1)
	})

	qm.hasSig = true
	qm.sigChan <- sig
}

func (qm *quicModule) OnNodeSettingsUpdated(settings node.NodeSettings) {
	qm.locker.Lock()
	defer qm.locker.Unlock()

	qm.setSig([]bool{true, false})
}

func (qm *quicModule) OnConfigUpdated(settings config.AppSettings) {

	qm.locker.Lock()
	defer qm.locker.Unlock()

	sigArr := make([]bool, 2)
	if sigArr[0] = !slices.Equal(qm.addrs, settings.NodeAddress); sigArr[0] {
		qm.addrs = settings.NodeAddress
	}

	if sigArr[1] = !slices.Equal(qm.publicAddrs, settings.PublicAddress); sigArr[1] {
		qm.publicAddrs = settings.PublicAddress
	}

	for _, sig := range sigArr {
		if sig {
			qm.setSig(sigArr)
			break
		}
	}
}

func (qm *quicModule) doRequest(ctx context.Context, conn quic.Connection, reader io.Reader, flag byte) (io.Reader, error) {
	stream, err := conn.OpenStream()
	if err != nil {
		return nil, err
	}

	errChan := make(chan error)
	go func() {
		_, err = stream.Write([]byte{flag})
		if err == nil {
			_, err = io.Copy(stream, reader)
		}
		if err == nil {
			err = stream.Close()
		}

		errChan <- err
	}()

	select {
	case err = <-errChan:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return stream, err
}

func (qm *quicModule) Do(ctx context.Context, conn quic.Connection, reader io.Reader) (io.Reader, error) {
	return qm.doRequest(ctx, conn, reader, QuicNodeStream)
}

func (qm *quicModule) Greet(ctx context.Context, conn quic.Connection) error {
	body := make([]byte, 128)
	_, err := rand.Read(body)
	if err != nil {
		return err
	}
	reader, err := qm.doRequest(ctx, conn, bytes.NewReader(body), QuicGreetStream)
	if err != nil {
		return err
	}
	readerBytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if bytes.Equal(body, readerBytes) {
		return nil
	}
	return constant.ErrInvalidNode

}

func (qm *quicModule) ParseNodeID(conn quic.Connection) (node.NodeID, error) {
	state := conn.ConnectionState()
	certificate := state.TLS.PeerCertificates[0]
	return x509.MarshalPKIXPublicKey(certificate.PublicKey)
}

func (qm *quicModule) Dial(ctx context.Context, addr string) (quic.Connection, error) {
	if qm.NodeModule == nil {
		return nil, constant.ErrUnavailable
	}
	settings := qm.NodeModule.NodeSettings()
	if !settings.Available() {
		return nil, constant.ErrUnavailable
	}
	certificate := settings.Certificate()
	tlsConf := &tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
	quicConf := &quic.Config{}
	return quic.DialAddr(ctx, addr, tlsConf, quicConf)
}

func (qm *quicModule) CreateNode(conn quic.Connection) (node.Node, error) {

	nodeId, err := qm.ParseNodeID(conn)
	if err != nil {
		return nil, err
	}

	nodeMgr := qm.NodeModule.NodeManager()

	qnode := &quicNode{
		quicModule: qm,
		conn:       conn,
		nodeId:     nodeId,
		mgr:        nodeMgr,
	}

	for {
		qnode.resourceId = nodeMgr.NewResourceID(node.NodeTypeAlive)
		_, ok := nodeMgr.SearchOrStore(node.Node(qnode))
		if !ok {
			break
		}
	}

	return node.Node(qnode), nil
}

func (qm *quicModule) CreateRoute(nodeId node.NodeID, addr string) (node.Node, error) {
	nodeMgr := qm.NodeModule.NodeManager()
	route := &quicRoute{
		quicModule: qm,
		nodeId:     nodeId,
		address:    addr,
	}

	for {
		route.resourceId = nodeMgr.NewResourceID(node.NodeTypeReachable)
		_, ok := nodeMgr.SearchOrStore(node.Node(route))
		if !ok {
			break
		}
	}

	return node.Node(route), nil
}

func (qm *quicModule) serve(stream quic.Stream, node node.Node) error {
	flags := make([]byte, 1)
	_, err := stream.Read(flags)
	if err != nil {
		return err
	}

	if flags[0] == QuicNodeStream {
		return qm.NodeModule.Serve(stream, node)
	}
	_, err = io.Copy(stream, stream)
	return err
}

func (qm *quicModule) Serve(conn quic.Connection) error {
	node, err := qm.CreateNode(conn)
	defer node.Close()

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			break
		}

		go qm.serve(stream, node)
	}

	return err
}

func (qm *quicModule) ServeBroadcast(payload []byte, ip string) error {

	if payload[0] != 1 || len(payload) <= 5 {
		return nil
	}
	payloadLen := len(payload)
	offset := 1
	nextOffset := offset + 2

	nodeIdLen := int(binary.BigEndian.Uint16(payload[offset:nextOffset]))
	offset = nextOffset
	nextOffset += nodeIdLen

	if nextOffset+2 > payloadLen {
		return nil
	}
	nodeId := node.NodeID(payload[offset:nextOffset])

	offset = nextOffset
	nextOffset += 2
	addressLen := int(binary.BigEndian.Uint16(payload[offset:nextOffset]))
	offset = nextOffset
	nextOffset += addressLen
	if nextOffset != payloadLen {
		return nil
	}

	address := string(payload[offset:nextOffset])
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	if host != ip {
		ipAddr, err := net.ResolveIPAddr("ip", host)
		if err != nil || !ipAddr.IP.IsUnspecified() {
			return err
		}
	}

	// TODO: check if route exists

	// Try add route
	ctx, _ := context.WithTimeoutCause(context.Background(), 3*time.Second, constant.ErrTimeout)
	conn, err := qm.Dial(ctx, address)
	if err == nil {
		defer conn.CloseWithError(0, "")
		err = qm.Greet(context.Background(), conn)
	}

	if err == nil {
		nodeId_, err := qm.ParseNodeID(conn)
		if err == nil && !bytes.Equal(nodeId, nodeId_) {
			return constant.ErrInvalidNode
		}
	}

	if err == nil {
		_, err = qm.CreateRoute(nodeId, address)
	}

	return err
}

func (qm *quicModule) Components() []runtime.Component {
	return []runtime.Component{
		runtime.NewComponent(qm, runtime.ComponentNoneScope),
		runtime.NewComponent[QuicModule](qm, runtime.ComponentExternalScope),
	}
}

func (qm *quicModule) Ready() error {

	qm.sigOnce.Do(func() {
		qm.sigChan = make(chan []bool, 1)
	})

	var servers []*quicServer
	defer qm.shutdownForQuic(servers)

	var cancel context.CancelCauseFunc
	defer qm.shutdownForBroadcast(cancel)

	for {
		sig := <-qm.sigChan

		qm.locker.Lock()
		qm.hasSig = false
		qm.locker.Unlock()

		if len(sig) <= 0 {
			break
		}

		if sig[0] {
			qm.shutdownForQuic(servers)
			servers = qm.serveForQuic()
		}

		if sig[1] {
			qm.shutdownForBroadcast(cancel)
			cancel = qm.serveForBroadcast()
		}

	}
	return nil
}
func (qm *quicModule) shutdownForQuic(servers []*quicServer) {
	if len(servers) > 0 {
		for _, server := range servers {
			server.Shutdown()
		}
	}
}

func (qm *quicModule) serveForQuic() []*quicServer {
	servers := make([]*quicServer, 0)
	addrs := qm.Addrs()
	for _, addr := range addrs {
		server := &quicServer{
			address:    addr,
			quicModule: qm,
		}

		servers = append(servers, server)
		go func(s *quicServer) {
			for {
				err := s.ListenAndServe()
				if errors.Is(err, quic.ErrServerClosed) {
					break
				}
				time.Sleep(6 * time.Second)
			}

		}(server)
	}
	return servers
}

func (qm *quicModule) shutdownForBroadcast(cancel context.CancelCauseFunc) {
	if cancel != nil {
		cancel(constant.ErrUnavailable)
	}
}

func (qm *quicModule) serveForBroadcast() context.CancelCauseFunc {

	ctx, cancel := context.WithCancelCause(context.Background())
	go func(ctx context.Context) {
		for {
			err := qm.deliverForBroadcast()
			if errors.Is(err, constant.ErrUnavailable) {
				return
			}

			select {
			case <-ctx.Done():
				break
			case <-time.After(15 * time.Second):
				continue
			}
		}
	}(ctx)
	return cancel
}

func (qm *quicModule) deliverForBroadcast() error {
	addrs := qm.PublicAddrs()
	addrsCount := len(addrs)
	if addrsCount <= 0 {
		return constant.ErrUnavailable
	}

	settings := qm.NodeModule.NodeSettings()
	if !settings.Available() {
		return constant.ErrUnavailable
	}

	nodeId := settings.NodeID()
	nodeIdLen := len(nodeId)

	var buffer []byte
	var offset int
	errs := make([]error, 0)
	for _, address := range addrs {
		addressLen := len(address)
		bufferSize := 4 + nodeIdLen + addressLen
		if len(buffer) != bufferSize {
			buffer = make([]byte, bufferSize)
			binary.BigEndian.PutUint16(buffer, uint16(len(nodeId)))
			copy(buffer[2:], nodeId)
			offset = 2 + nodeIdLen
			binary.BigEndian.PutUint16(buffer[offset:], uint16(addressLen))
			offset += 2
		}
		copy(buffer[offset:], []byte(address))
		err := qm.Broadcast.Deliver(buffer)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) <= 0 {
		return nil
	}
	return errors.Join(errs...)
}
