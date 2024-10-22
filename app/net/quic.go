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
	"pan/app/bootstrap"
	"pan/app/config"
	"pan/app/constant"
	appNode "pan/app/node"
	"slices"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

type quicNode struct {
	resourceId appNode.NodeResourceID
	nodeId     appNode.NodeID
	conn       quic.Connection
	quicModule QuicModule
	mgr        appNode.NodeManager
}

func (qn *quicNode) ID() appNode.NodeID {
	return qn.nodeId
}

func (qn *quicNode) Type() appNode.NodeType {
	return appNode.NodeTypeAlive
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

func (qn *quicNode) ResourceID() appNode.NodeResourceID {
	return qn.resourceId
}

type quicRoute struct {
	resourceId    appNode.NodeResourceID
	nodeId        appNode.NodeID
	address       string
	quicModule    *quicModule
	failures      uint8
	failureLocker sync.RWMutex
	routeId       []byte
}

func (qr *quicRoute) ID() appNode.NodeID {
	return qr.nodeId
}

func (qr *quicRoute) Type() appNode.NodeType {
	return appNode.NodeTypeReachable
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
	qr.quicModule.destroyRoute(qr)
	return nil
}

func (qr *quicRoute) ResourceID() appNode.NodeResourceID {
	return qr.resourceId
}

const (
	QuicNodeStream byte = iota + 1
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
	ParseNodeID(quic.Connection) (appNode.NodeID, error)
	CreateNode(quic.Connection) (appNode.Node, error)
	CreateRoute(appNode.NodeID, string) (appNode.Node, error)
}

type quicModule struct {
	Broadcast  Broadcast
	NodeModule appNode.NodeModule

	publicAddrs []string
	addrs       []string
	locker      sync.RWMutex
	sigChan     chan []bool
	sigOnce     sync.Once
	hasSig      bool

	routes      []*quicRoute
	routeLocker sync.RWMutex

	wg sync.WaitGroup
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

func (qm *quicModule) SigChan() chan []bool {
	qm.sigOnce.Do(func() {
		qm.sigChan = make(chan []bool, 1)
	})
	return qm.sigChan
}

func (qm *quicModule) setSig(sig []bool) {
	if qm.hasSig {
		return
	}

	qm.hasSig = true
	qm.SigChan() <- sig
}

func (qm *quicModule) OnNodeSettingsUpdated(settings appNode.NodeSettings) {
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
		break
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
	n, err := rand.Read(body)
	if err != nil && n != 128 {
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

func (qm *quicModule) ParseNodeID(conn quic.Connection) (appNode.NodeID, error) {
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

func (qm *quicModule) CreateNode(conn quic.Connection) (appNode.Node, error) {

	nodeId, err := qm.ParseNodeID(conn)
	if err != nil {
		return nil, err
	}

	nodeModule := qm.NodeModule

	qnode := &quicNode{
		quicModule: qm,
		conn:       conn,
		nodeId:     nodeId,
		mgr:        nodeModule.NodeManager(),
	}

	for {
		qnode.resourceId = nodeModule.NewResourceID(appNode.NodeTypeAlive)
		err = nodeModule.Control(appNode.Node(qnode))
		if err == nil || err != constant.ErrConflict {
			break
		}
	}

	return appNode.Node(qnode), err
}

func (qm *quicModule) CreateRoute(nodeId appNode.NodeID, addr string) (appNode.Node, error) {
	nodeModule := qm.NodeModule

	routeId := make([]byte, 0)
	routeId = append(routeId, nodeId...)
	routeId = append(routeId, []byte(addr)...)

	// check conflict
	qm.routeLocker.RLock()
	_, ok := slices.BinarySearchFunc(qm.routes, routeId, qm.compareWithQuicRoute)
	if ok {
		qm.routeLocker.RUnlock()
		return nil, constant.ErrConflict
	}
	qm.routeLocker.RUnlock()

	route := &quicRoute{
		quicModule: qm,
		nodeId:     nodeId,
		address:    addr,
		routeId:    routeId,
	}

	var err error
	for {
		route.resourceId = nodeModule.NewResourceID(appNode.NodeTypeReachable)
		err = nodeModule.Control(appNode.Node(route))
		if err == nil || err != constant.ErrConflict {
			break
		}
	}

	// add to routes
	qm.routeLocker.Lock()
	idx, _ := slices.BinarySearchFunc(qm.routes, routeId, qm.compareWithQuicRoute)
	qm.routes = slices.Insert(qm.routes, idx, route)
	qm.routeLocker.Unlock()

	return appNode.Node(route), err
}

func (qm *quicModule) destroyRoute(route *quicRoute) {
	nodeModule := qm.NodeModule
	if nodeModule != nil {
		mgr := nodeModule.NodeManager()
		mgr.Delete(route)
	}

	qm.routeLocker.Lock()
	defer qm.routeLocker.Unlock()
	idx, ok := slices.BinarySearchFunc(qm.routes, route.routeId, qm.compareWithQuicRoute)
	if !ok {
		return
	}

	for i := idx; i < len(qm.routes); i++ {
		routeItem := qm.routes[i]
		if routeItem == route {
			qm.routes = slices.Delete(qm.routes, i, i+1)
			break
		}
		if !bytes.Equal(routeItem.routeId, route.routeId) {
			break
		}
	}

}

func (qm *quicModule) compareWithQuicRoute(route *quicRoute, routeId []byte) int {
	return bytes.Compare(route.routeId, routeId)
}

func (qm *quicModule) serve(stream quic.Stream, node appNode.Node) error {
	defer stream.Close()
	flags := make([]byte, 1)
	n, err := stream.Read(flags)
	if err != nil && n != 1 {
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
	if err != nil {
		conn.CloseWithError(quic.ApplicationErrorCode(quic.ConnectionRefused), "")
		return err
	}
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

	var nodeSettings appNode.NodeSettings
	if qm.NodeModule != nil {
		nodeSettings = qm.NodeModule.NodeSettings()
	}
	if nodeSettings == nil {
		return constant.ErrUnavailable
	}

	payloadLen := len(payload)
	offset := 0
	nextOffset := offset + 2

	nodeIdLen := int(binary.BigEndian.Uint16(payload[offset:nextOffset]))
	offset = nextOffset
	nextOffset += nodeIdLen

	if nextOffset+2 > payloadLen {
		return nil
	}
	nodeId := appNode.NodeID(payload[offset:nextOffset])

	// ingore self by node id
	if bytes.Equal(nodeSettings.NodeID(), nodeId) {
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

	// Try add route
	route, err := qm.CreateRoute(nodeId, address)
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeoutCause(context.Background(), 5*time.Second, constant.ErrTimeout)
	err = route.Greet(ctx)
	if err != nil {
		route.Close()
	}
	return err
}

func (qm *quicModule) Components() []bootstrap.Component {
	return []bootstrap.Component{
		bootstrap.NewComponent(qm, bootstrap.ComponentNoneScope),
		bootstrap.NewComponent[QuicModule](qm, bootstrap.ComponentExternalScope),
	}
}

func (qm *quicModule) Ready() error {

	var servers []*quicServer
	defer qm.shutdownForQuic(servers)

	var cancel context.CancelCauseFunc
	defer qm.shutdownForBroadcast(cancel)

	for {
		sig := <-qm.SigChan()

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
		qm.wg.Wait()
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
		qm.wg.Add(1)
		go func(s *quicServer) {
			defer qm.wg.Done()
			_ = s.ListenAndServe()
			// TODO: write error into log
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
	broadcastLoop:
		for {
			err := qm.deliverForBroadcast()
			if errors.Is(err, constant.ErrUnavailable) {
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
