package net

import (
	"bytes"
	"context"
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

type quicNodeStream struct {
	quic.Stream
	node    *quicNode
	closed  bool
	isServe bool
}

func (qs *quicNodeStream) Read(p []byte) (n int, err error) {
	n, err = qs.Stream.Read(p)
	if (err != nil || n == 0) && !qs.isServe {
		qs.Close()
	}
	return
}

func (qs *quicNodeStream) Close() error {
	if qs.closed {
		return nil
	}
	qs.closed = true
	qs.node.decreaseStream()
	var err error
	if qs.isServe {
		err = qs.Stream.Close()
	} else {
		qs.CancelRead(quic.StreamErrorCode(quic.NoError))
	}

	return err
}

type quicNode struct {
	resourceId  appNode.NodeResourceID
	nodeId      appNode.NodeID
	conn        quic.Connection
	quicModule  QuicModule
	mgr         appNode.NodeManager
	streamCount int8
	rw          sync.RWMutex
}

func (qn *quicNode) ID() appNode.NodeID {
	return qn.nodeId
}

func (qn *quicNode) Type() appNode.NodeType {
	return appNode.NodeTypeAlive
}

func (qn *quicNode) Do(ctx context.Context, reader io.Reader) (io.ReadCloser, error) {
	qn.increaseStream()
	stream, err := qn.quicModule.Do(ctx, qn.conn, reader)
	if err != nil {
		qn.decreaseStream()
		return nil, err
	}
	return &quicNodeStream{Stream: stream, node: qn}, err
}

func (qn *quicNode) Close() error {
	qn.rw.Lock()
	qn.streamCount = -1
	qn.rw.Unlock()

	qn.mgr.Delete(qn)
	return qn.conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
}

func (qn *quicNode) ResourceID() appNode.NodeResourceID {
	return qn.resourceId
}

func (qn *quicNode) IsIdle() bool {
	qn.rw.RLock()
	defer qn.rw.RUnlock()
	return qn.streamCount >= 0 && qn.streamCount < 2
}

func (qn *quicNode) increaseStream() {
	qn.rw.Lock()
	defer qn.rw.Unlock()
	qn.streamCount++
}

func (qn *quicNode) decreaseStream() {
	qn.rw.Lock()
	defer qn.rw.Unlock()
	if qn.streamCount < 0 {
		return
	}
	qn.streamCount--
}

type quicRoute struct {
	resourceId    appNode.NodeResourceID
	nodeId        appNode.NodeID
	address       string
	quicModule    *quicModule
	failures      uint8
	failureLocker sync.RWMutex
	routeId       []byte
	closed        bool
	closedRW      sync.RWMutex
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
		nodeId, err := qr.quicModule.parseNodeID(conn)
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

func (qr *quicRoute) Do(ctx context.Context, reader io.Reader) (io.ReadCloser, error) {
	conn, err := qr.Dial(ctx)
	if err != nil {
		return nil, err
	}

	node, err := qr.quicModule.Serve(conn)
	if err != nil {
		return nil, err
	}
	return node.Do(ctx, reader)
}

func (qr *quicRoute) Close() error {
	qr.closedRW.Lock()
	qr.closed = true
	qr.closedRW.Unlock()

	qr.quicModule.purgeRoute(qr)
	return nil
}

func (qr *quicRoute) ResourceID() appNode.NodeResourceID {
	return qr.resourceId
}

func (qr *quicRoute) IsIdle() bool {
	qr.closedRW.RLock()
	defer qr.closedRW.RUnlock()
	return !qr.closed
}

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

func (qs *quicServer) ListenAndServe(ctx context.Context) error {

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
		conn, err := ln.Accept(ctx)
		if err != nil && conn != nil {
			conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
		}
		if err == nil {
			_, err = qs.quicModule.Serve(conn)
		}
		if errors.Is(err, quic.ErrServerClosed) || errors.Is(err, context.Canceled) {
			break
		}

	}

	return err
}

type QuicModule interface {
	Serve(quic.Connection) (appNode.Node, error)
	Do(context.Context, quic.Connection, io.Reader) (quic.Stream, error)
	Dial(context.Context, string) (quic.Connection, error)
	Route(appNode.NodeID, string) error
}

type quicModule struct {
	Broadcast  Broadcast
	NodeModule appNode.NodeModule

	publicAddrs     []string
	addrs           []string
	locker          sync.RWMutex
	reloadChan      chan struct{}
	reloadOnce      sync.Once
	reloadServe     bool
	reloadBroadcast bool

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

func (qm *quicModule) ReloadChan() chan struct{} {
	qm.reloadOnce.Do(func() {
		qm.reloadChan = make(chan struct{}, 1)
	})
	return qm.reloadChan
}

func (qm *quicModule) OnNodeSettingsUpdated(settings appNode.NodeSettings) {
	qm.locker.Lock()
	defer qm.locker.Unlock()

	if qm.reloadServe {
		return
	}
	qm.reloadServe = true
	qm.ReloadChan() <- struct{}{}
}

func (qm *quicModule) OnConfigUpdated(settings config.AppSettings) {

	qm.locker.Lock()
	defer qm.locker.Unlock()

	var reloadServe bool
	if reloadServe = !slices.Equal(qm.addrs, settings.NodeAddress); reloadServe {
		qm.addrs = settings.NodeAddress
	}

	if !qm.reloadServe && reloadServe {
		qm.reloadServe = reloadServe
	}

	var reloadBroadcast bool
	if reloadBroadcast = !slices.Equal(qm.publicAddrs, settings.PublicAddress); reloadBroadcast {
		qm.publicAddrs = settings.PublicAddress
	}

	if !qm.reloadBroadcast && reloadBroadcast {
		qm.reloadBroadcast = reloadBroadcast
	}

	if !qm.reloadServe && !qm.reloadBroadcast {
		return
	}

	qm.ReloadChan() <- struct{}{}
}

func (qm *quicModule) Do(ctx context.Context, conn quic.Connection, reader io.Reader) (quic.Stream, error) {
	stream, err := conn.OpenStream()
	if err != nil {
		return nil, err
	}

	errChan := make(chan error)
	go func() {
		_, err = io.Copy(stream, reader)
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

func (qm *quicModule) parseNodeID(conn quic.Connection) (appNode.NodeID, error) {
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

func (qm *quicModule) createNode(conn quic.Connection) (*quicNode, error) {

	nodeId, err := qm.parseNodeID(conn)
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

	return qnode, err
}

func (qm *quicModule) Route(nodeId appNode.NodeID, addr string) error {
	nodeModule := qm.NodeModule

	routeId := make([]byte, 0)
	routeId = append(routeId, nodeId...)
	routeId = append(routeId, []byte(addr)...)

	// check conflict
	qm.routeLocker.RLock()
	_, ok := slices.BinarySearchFunc(qm.routes, routeId, qm.compareWithQuicRoute)
	if ok {
		qm.routeLocker.RUnlock()
		return constant.ErrConflict
	}
	qm.routeLocker.RUnlock()
	//

	// try connect
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := qm.Dial(ctx, addr)
	if err == nil {
		err = conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
	}
	if err != nil {
		return err
	}
	//

	route := &quicRoute{
		quicModule: qm,
		nodeId:     nodeId,
		address:    addr,
		routeId:    routeId,
	}

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

	return err
}

func (qm *quicModule) purgeRoute(route *quicRoute) {
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

func (qm *quicModule) serve(node *quicNode) error {
	defer node.Close()

	conn := node.conn

	var err error
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			break
		}
		node.increaseStream()
		go qm.NodeModule.Serve(&quicNodeStream{Stream: stream, node: node, isServe: true}, node)
	}

	return err
}

func (qm *quicModule) Serve(conn quic.Connection) (appNode.Node, error) {
	node, err := qm.createNode(conn)
	if err == nil {
		go qm.serve(node)
	}
	return node, err

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

	// Try  route
	if err == nil {
		err = qm.Route(nodeId, address)
	}

	return err
}

func (qm *quicModule) Components() []bootstrap.Component {
	return []bootstrap.Component{
		bootstrap.NewComponent(qm, bootstrap.ComponentNoneScope),
		bootstrap.NewComponent[QuicModule](qm, bootstrap.ComponentExternalScope),
	}
}

func (qm *quicModule) Ready(ctx context.Context) error {

	var servers []*quicServer
	defer qm.shutdownForQuic(servers)

	var cancel context.CancelCauseFunc
	defer qm.shutdownForBroadcast(cancel)

	var err error
	closed := false

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-qm.ReloadChan():
		}

		qm.locker.Lock()
		reloadServe := qm.reloadServe
		qm.reloadServe = false

		reloadBroadcast := qm.reloadBroadcast
		qm.reloadBroadcast = false
		qm.locker.Unlock()

		if closed {
			break
		}

		if reloadServe {
			qm.shutdownForQuic(servers)
			servers = qm.serveForQuic(ctx)
		}

		if reloadBroadcast {
			qm.shutdownForBroadcast(cancel)
			causeCtx, causeCancel := context.WithCancelCause(ctx)
			cancel = causeCancel
			qm.serveForBroadcast(causeCtx)
		}

	}
	return err
}
func (qm *quicModule) shutdownForQuic(servers []*quicServer) {
	if len(servers) > 0 {
		for _, server := range servers {
			server.Shutdown()
		}
		qm.wg.Wait()
	}
}

func (qm *quicModule) serveForQuic(ctx context.Context) []*quicServer {
	servers := make([]*quicServer, 0)
	addrs := qm.Addrs()
	for _, addr := range addrs {
		server := &quicServer{
			address:    addr,
			quicModule: qm,
		}

		servers = append(servers, server)
		qm.wg.Add(1)
		go func(s *quicServer, c context.Context) {
			defer qm.wg.Done()
			_ = s.ListenAndServe(c)
			// TODO: write error into log
		}(server, ctx)
	}
	return servers
}

func (qm *quicModule) shutdownForBroadcast(cancel context.CancelCauseFunc) {
	if cancel != nil {
		cancel(constant.ErrUnavailable)
	}
}

func (qm *quicModule) serveForBroadcast(ctx context.Context) {

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
