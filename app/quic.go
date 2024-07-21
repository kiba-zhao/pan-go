package app

import (
	"bytes"
	"cmp"
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"pan/app/cache"
	"pan/runtime"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

type quicNode struct {
	seq  uint
	id   NodeID
	conn quic.Connection
}

func (qn *quicNode) ID() NodeID {
	return qn.id
}

func (qn *quicNode) Close() error {
	return qn.conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
}

func (qn *quicNode) HashCode() uint {
	return qn.seq
}

type quicNestNodeBucket = *cache.NestBucket[NodeID, uint, *quicNode]
type quicNodeBucket = cache.Bucket[NodeID, quicNestNodeBucket]

type quicNodeConnManager struct {
	bucket quicNodeBucket
	seq    uint
	lock   sync.Mutex
}

func (mgr *quicNodeConnManager) Search(id NodeID) (*quicNode, bool) {
	bucket, ok := mgr.bucket.Search(id)
	if !ok {
		return nil, false
	}

	size := bucket.Size()
	if size <= 0 {
		return nil, false
	}

	idx, ok := bucket.Index(mgr.seq)
	if ok {
		return bucket.At(idx)
	}

	if idx > 0 {
		return bucket.At(idx - 1)
	}
	return bucket.At(size - 1)
}

func (mgr *quicNodeConnManager) Store(conn quic.Connection) (*quicNode, error) {

	state := conn.ConnectionState()
	cert := state.TLS.PeerCertificates[0]
	pubKey, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return nil, err
	}

	bucket, ok := mgr.bucket.Search(pubKey)
	if !ok {
		bucket = &cache.NestBucket[NodeID, uint, *quicNode]{}
		bucket.Code = pubKey
		simpleBucket := cache.NewBucket[uint, *quicNode](cmp.Compare[uint])
		bucket.Bucket = cache.WrapSyncBucket(simpleBucket)
		bucket, _ = mgr.bucket.SearchOrStore(bucket)
	}

	item := &quicNode{}
	item.id = pubKey
	item.conn = conn
	for {
		mgr.lock.Lock()
		item.seq = mgr.seq
		mgr.seq++
		mgr.lock.Unlock()

		_, ok := bucket.SearchOrStore(item)
		if !ok {
			break
		}
	}

	return item, err
}

func (mgr *quicNodeConnManager) Delete(node *quicNode) {
	bucket, ok := mgr.bucket.Search(node.ID())
	if ok {
		bucket.Delete(node)
	}
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
		return ErrUnavailable
	}

	qs.locker.Lock()
	qs.ln = nil
	qs.locker.Unlock()
	return ln.Close()
}

func (qs *quicServer) ListenAndServe() error {

	if qs.quicModule == nil || qs.quicModule.NodeModule == nil {
		return ErrUnavailable
	}

	nodeSettings := qs.quicModule.NodeModule.NodeSettings()
	if !nodeSettings.Available() {
		return ErrUnavailable
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

func (qs *quicServer) HashCode() string {
	return qs.address
}

type QuicModule interface {
	Serve(quic.Connection) error
}

type quicModule struct {
	Broadcast     Broadcast
	NodeModule    NodeModule
	connMgr       *quicNodeConnManager
	connMgrLocker sync.RWMutex

	addresses []string
	locker    sync.RWMutex
	sigChan   chan bool
	sigOnce   sync.Once
	hasSig    bool

	modules []interface{}
	once    sync.Once
}

func (qm *quicModule) Addresses() []string {
	qm.locker.RLock()
	defer qm.locker.RUnlock()
	return qm.addresses
}

func (qm *quicModule) setSig(sig bool) {
	if qm.hasSig {
		return
	}

	qm.sigOnce.Do(func() {
		qm.sigChan = make(chan bool, 1)
	})

	qm.hasSig = true
	qm.sigChan <- sig
}

func (qm *quicModule) OnNodeSettingsUpdated(settings NodeSettings) {
	qm.locker.Lock()
	defer qm.locker.Unlock()

	qm.setSig(true)
}

func (qm *quicModule) OnConfigUpdated(settings AppSettings) {

	qm.locker.Lock()
	defer qm.locker.Unlock()

	if slices.Equal(qm.addresses, settings.NodeAddress) {
		return
	}

	qm.addresses = settings.NodeAddress
	qm.setSig(true)
}

func (qm *quicModule) serve(stream quic.Stream, node *quicNode) error {
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

	qm.connMgrLocker.RLock()
	mgr := qm.connMgr
	qm.connMgrLocker.RUnlock()

	if mgr == nil {
		return ErrUnavailable
	}

	node_, err := mgr.Store(conn)
	if err != nil {
		return conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
	}
	defer mgr.Delete(node_)
	defer node_.Close()

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			break
		}

		go qm.serve(stream, node_)
	}

	return err
}

func (qm *quicModule) Do(nodeId NodeID, reader io.Reader, ctx context.Context) (io.Reader, error) {
	qm.connMgrLocker.RLock()
	mgr := qm.connMgr
	qm.connMgrLocker.RUnlock()

	if mgr == nil {
		return nil, ErrUnavailable
	}

	node, ok := mgr.Search(nodeId)
	if !ok {
		return nil, ErrNotFound
	}

	stream, err := node.conn.OpenStream()
	defer func() {
		if err != nil && !errors.Is(err, ctx.Err()) {
			mgr.Delete(node)
		}
	}()
	if err != nil {
		return nil, err
	}

	errChan := make(chan error)
	go func() {
		_, err = stream.Write([]byte{QuicNodeStream})
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

func (qm *quicModule) Lookup(nodeId NodeID, ctx context.Context) error {
	qm.connMgrLocker.RLock()
	mgr := qm.connMgr
	qm.connMgrLocker.RUnlock()
	_, ok := mgr.Search(nodeId)
	if !ok {
		return ErrNotFound
	}

	return nil
}

func (qm *quicModule) Components() []runtime.Component {
	return []runtime.Component{
		runtime.NewComponent(qm, runtime.ComponentNoneScope),
		runtime.NewComponent[QuicModule](qm, runtime.ComponentExternalScope),
	}
}

func (qm *quicModule) Ready() error {
	var serverMgr cache.Bucket[string, *quicServer]
	qm.sigOnce.Do(func() {
		qm.sigChan = make(chan bool, 1)
	})

	qm.connMgrLocker.Lock()
	if qm.connMgr == nil {
		nodeBucket := cache.NewBucket[NodeID, quicNestNodeBucket](bytes.Compare)
		qm.connMgr = &quicNodeConnManager{bucket: cache.WrapSyncBucket(nodeBucket)}
	}
	qm.connMgrLocker.Unlock()

	for {
		sig := <-qm.sigChan

		qm.locker.Lock()
		qm.hasSig = false
		qm.locker.Unlock()

		if serverMgr != nil {
			for _, item := range serverMgr.Items() {
				item.Shutdown()
			}
		}

		if !sig {
			break
		}

		serverMgr_ := cache.NewBucket[string, *quicServer](cmp.Compare[string])
		serverMgr = cache.WrapSyncBucket(serverMgr_)
		addresses := qm.Addresses()
		for _, address := range addresses {
			server := &quicServer{
				address:    address,
				quicModule: qm,
			}

			serverMgr.Store(server)
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

	}
	return nil
}

func (qm *quicModule) Modules() []interface{} {

	qm.once.Do(func() {
		bucket_ := cache.NewBucket[NodeID, quicNestRouteBucket](bytes.Compare)
		network := &quicNetwork{routeBucket: cache.WrapSyncBucket(bucket_), quicModule: qm}

		qm.modules = []interface{}{
			network,
		}
	})
	return qm.modules
}

type quicRoute struct {
	id      NodeID
	address string
}

func (qr *quicRoute) HashCode() string {
	return qr.address
}

type quicNestRouteBucket = *cache.NestBucket[NodeID, string, *quicRoute]
type quicRouteBucket = cache.Bucket[NodeID, quicNestRouteBucket]

type quicNetwork struct {
	routeBucket quicRouteBucket
	quicModule  *quicModule

	quicPorts []int
	locker    sync.RWMutex
	sigChan   chan bool
	sigOnce   sync.Once
	hasSig    bool
}

func (qn *quicNetwork) ServeBroadcast(payload []byte, ip string) error {
	// extract 1 byte version, 2 byte nodeIdLen, nodeId, 2 byte addrLen, address
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
	if (payloadLen-nextOffset)%2 != 0 {
		return nil
	}
	nodeId := NodeID(payload[offset:nextOffset])

	for nextOffset < payloadLen {
		offset = nextOffset
		nextOffset += 2

		port := int(binary.BigEndian.Uint16(payload[offset:nextOffset]))
		address := net.JoinHostPort(ip, strconv.Itoa(port))

		// check if route exists
		bucket, ok := qn.routeBucket.Search(nodeId)
		if ok {
			_, ok := bucket.Search(address)
			if ok {
				continue
			}
		}

		// add route
		body := make([]byte, 128)
		_, err := rand.Read(body)
		if err != nil {
			return err
		}

		route := &quicRoute{
			id:      nodeId,
			address: address,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer func(cancel context.CancelFunc) {
			cancel()
		}(cancel)
		reader, err := qn.doRequest(route, bytes.NewReader(body), ctx, QuicGreetStream)
		if err != nil {
			continue
		}

		resBody, err := io.ReadAll(reader)
		if err != nil || !bytes.Equal(resBody, body) {
			continue
		}

		bucket, ok = qn.routeBucket.Search(nodeId)
		if !ok {
			bucket_ := &cache.NestBucket[NodeID, string, *quicRoute]{}
			bucket_.Code = nodeId
			simpleBucket_ := cache.NewBucket[string, *quicRoute](cmp.Compare[string])
			bucket_.Bucket = cache.WrapSyncBucket(simpleBucket_)
			bucket, _ = qn.routeBucket.SearchOrStore(bucket_)
		}

		bucket.SearchOrStore(route)
	}

	return nil
}

func (qn *quicNetwork) doRequest(route *quicRoute, reader io.Reader, ctx context.Context, flag byte) (io.Reader, error) {
	if qn.quicModule == nil || qn.quicModule.NodeModule == nil {
		return nil, ErrUnavailable
	}

	nodeSettings := qn.quicModule.NodeModule.NodeSettings()
	if !nodeSettings.Available() {
		return nil, ErrUnavailable
	}

	certificate := nodeSettings.Certificate()
	tlsConf := &tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
	quicConf := &quic.Config{}
	conn, err := quic.DialAddr(ctx, route.address, tlsConf, quicConf)
	if err != nil {
		return nil, err
	}

	if flag == QuicNodeStream {
		go qn.quicModule.Serve(conn)
	} else {
		defer conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
	}

	state := conn.ConnectionState()
	cert := state.TLS.PeerCertificates[0]
	pubKey, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(pubKey, route.id) {
		return nil, ErrInvalidNode
	}

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

func (qn *quicNetwork) Do(nodeId NodeID, reader io.Reader, ctx context.Context) (io.Reader, error) {
	bucket, ok := qn.routeBucket.Search(nodeId)
	if !ok {
		return nil, ErrNotFound
	}

	route, ok := bucket.At(0)
	if !ok {
		return nil, ErrNotFound
	}

	return qn.doRequest(route, reader, ctx, QuicNodeStream)
}

func (qn *quicNetwork) Greet(nodeId NodeID, ctx context.Context) error {
	bucket, ok := qn.routeBucket.Search(nodeId)
	if !ok {
		return ErrNotFound
	}

	body := make([]byte, 128)
	_, err := rand.Read(body)
	if err != nil {
		return err
	}
	err = ErrNotFound
	for _, route := range bucket.Items() {
		for failedCount := 0; failedCount < 3; failedCount++ {
			resReader, reqErr := qn.doRequest(route, bytes.NewReader(body), ctx, QuicGreetStream)
			err = reqErr
			if err != nil {
				if errors.Is(err, ctx.Err()) {
					break
				}
				continue
			}
			resBody, resErr := io.ReadAll(resReader)
			err = resErr
			if err != nil {
				continue
			}
			if bytes.Equal(resBody, body) {
				break
			}
			err = ErrNotFound
		}

		if err == nil || errors.Is(err, ctx.Err()) {
			break
		}

	}

	return err
}

func (qn *quicNetwork) setSig(sig bool) {
	if qn.hasSig {
		return
	}

	qn.sigOnce.Do(func() {
		qn.sigChan = make(chan bool, 1)
	})

	qn.hasSig = true
	qn.sigChan <- sig
}

func (qn *quicNetwork) OnConfigUpdated(settings AppSettings) {
	qn.locker.Lock()
	defer qn.locker.Unlock()

	if slices.Equal(qn.quicPorts, settings.BroadcasQuicPorts) {
		return
	}

	qn.quicPorts = settings.BroadcasQuicPorts
	qn.setSig(true)
}

func (qn *quicNetwork) Ready() error {

	var ctx context.Context
	var cancel context.CancelCauseFunc
	qn.sigOnce.Do(func() {
		qn.sigChan = make(chan bool, 1)
	})

	for {
		sig := <-qn.sigChan

		if !sig {
			break
		}

		qn.locker.Lock()
		qn.hasSig = false
		qn.locker.Unlock()

		// TODO: do broadcast every 15s
		if cancel != nil {
			cancel(nil)
		}

		ctx, cancel = context.WithCancelCause(context.Background())
		go func(ctx context.Context) {
		innerLoop:
			for {
				err := qn.Deliver()
				if errors.Is(err, ErrUnavailable) {
					return
				}

				select {
				case <-ctx.Done():
					break innerLoop
				case <-time.After(15 * time.Second):
					continue
				}
			}
		}(ctx)
	}
	return nil
}

func (qn *quicNetwork) QuicPorts() []int {

	qn.locker.RLock()
	defer qn.locker.RUnlock()

	return qn.quicPorts
}

func (qn *quicNetwork) Deliver() error {
	ports := qn.QuicPorts()
	portCount := len(ports)
	if portCount <= 0 {
		return ErrUnavailable
	}

	if qn.quicModule == nil || qn.quicModule.NodeModule == nil {
		return ErrUnavailable
	}

	settings := qn.quicModule.NodeModule.NodeSettings()
	if !settings.Available() {
		return ErrUnavailable
	}

	nodeId := settings.NodeID()
	nodeIdLen := len(nodeId)
	buffer := make([]byte, nodeIdLen+2+2*portCount)

	binary.BigEndian.PutUint16(buffer, uint16(len(nodeId)))
	copy(buffer[2:], nodeId)
	offset := 2 + nodeIdLen
	for i := 0; i < portCount; i++ {
		port := ports[i]
		binary.BigEndian.PutUint16(buffer[offset:], uint16(port))
		offset += 2
	}

	return qn.quicModule.Broadcast.Deliver(buffer)
}
