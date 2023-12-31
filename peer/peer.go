package peer

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"

	"io"
	"net"
	"pan/core"
	"pan/memory"
	"sync"

	"github.com/google/uuid"
)

const (
	QUICNodeType = uint8(iota)
	TCPNodeType
)

const (
	NormalAuthenticateMode = uint8(iota)
	TestOnlyAuthenticateMode
	OpenAuthenticateMode
)

const (
	OnlinePeerState = PeerState(iota)
	OfflinePeerState
	UnknownPeerState
)

type PeerId uuid.UUID

type PeerState uint

// compareUUIDBytes ...
func compareUUID(prev, next uuid.UUID) int {
	prevHigh := binary.BigEndian.Uint64(prev[:8])
	nextHigh := binary.BigEndian.Uint64(next[:8])
	if prevHigh > nextHigh {
		return 1
	}
	if prevHigh < nextHigh {
		return -1
	}

	prevLow := binary.BigEndian.Uint64(prev[8:])
	nextLow := binary.BigEndian.Uint64(next[8:])
	if prevLow > nextLow {
		return 1
	}
	if prevLow < nextLow {
		return -1
	}
	return 0
}

// comparePeerId ...
func comparePeerId(prev, next PeerId) int {
	return compareUUID(uuid.UUID(prev), uuid.UUID(next))
}

type Peer interface {
	Stat(id PeerId) PeerState
	Connect(dialerType uint8, addr []byte) (Node, error)
	Attach(dialer NodeDialer) error
	Detach(dialer NodeDialer)
	Authenticate(node Node, mode uint8) (PeerId, error)
	AcceptAuthenticate(ctx context.Context, node Node)
	Open(id PeerId) (Node, error)
	Request(node Node, body io.Reader, method []byte, headers ...*HeaderSegment) (*Response, error)
	AcceptServe(ctx context.Context, serve NodeServe)
	Accept(ctx context.Context, node Node, peerId PeerId)
}

type PeerIdGenerator interface {
	Generate(baseId []byte, node Node) (PeerId, error)
}

type PeerPassport struct {
	IsPeerId bool
}

type SimplePeerIdGenerator struct {
	*memory.Bucket[*PeerPassport, uuid.UUID]
	defaultDeny bool
}

// Generate ...
func (pg *SimplePeerIdGenerator) Generate(baseId []byte, node Node) (peerId PeerId, err error) {
	space, err := uuid.FromBytes(baseId)
	if err != nil {
		return
	}

	cert := node.Certificate()
	pubKey, err := core.ExtractPublicKeyFromCert(cert)
	if err != nil {
		return
	}

	id := uuid.NewSHA1(space, pubKey)

	pass := !pg.defaultDeny
	items := pg.FindBlockItems(id)
	for _, item := range items {
		if item.Expired() {
			continue
		}
		passport := item.Value()
		if passport.IsPeerId {
			pass = pg.defaultDeny
			break
		}
	}

	if pass == pg.defaultDeny {
		peerId = PeerId(id)
		return
	}

	items = pg.FindBlockItems(space)
	for _, item := range items {
		if item.Expired() {
			continue
		}
		passport := item.Value()
		if !passport.IsPeerId {
			pass = pg.defaultDeny
			break
		}
	}

	if !pass {
		err = errors.New("Deny Peer Id")
		return
	}

	peerId = PeerId(id)
	return
}

// NewPeerIdGenerator ...
func NewPeerIdGenerator(defaultDeny bool) *SimplePeerIdGenerator {

	bucket := memory.NewBucket[*PeerPassport, uuid.UUID](compareUUID)

	generator := new(SimplePeerIdGenerator)
	generator.Bucket = bucket
	generator.defaultDeny = defaultDeny

	return generator
}

type peerRoute struct {
	rw        *sync.RWMutex
	NodeType  uint8
	Addr      []byte
	FailedNum uint8
}

type peerDialer struct {
	dialerMap map[uint8]NodeDialer
	rw        *sync.RWMutex
}

// Connect ...
func (pd *peerDialer) Connect(dialerType uint8, addr []byte) (node Node, err error) {
	pd.rw.RLock()
	dialer, ok := pd.dialerMap[dialerType]
	pd.rw.RUnlock()

	if ok == false {
		err = errors.New("Not Found node dialer")
		return
	}

	node, err = dialer.Connect(addr)
	return
}

// Attach ...
func (pd *peerDialer) Attach(dialer NodeDialer) (err error) {

	t := dialer.Type()

	pd.rw.Lock()
	cdialer, ok := pd.dialerMap[t]
	if ok == true && cdialer != dialer {
		err = errors.New("Duplicate node dialer")
	} else if ok == false {
		pd.dialerMap[t] = dialer
	}
	pd.rw.Unlock()

	return
}

// Detach ...
func (pd *peerDialer) Detach(dialer NodeDialer) {
	t := dialer.Type()

	pd.rw.Lock()
	cdialer, ok := pd.dialerMap[t]
	if ok == true && cdialer == dialer {
		delete(pd.dialerMap, t)
	}
	pd.rw.Unlock()
}

type peerSt struct {
	*peerDialer
	generator    PeerIdGenerator
	router       *memory.Bucket[*peerRoute, PeerId]
	bucket       *memory.Bucket[Node, PeerId]
	baseId       uuid.UUID
	app          core.App[Context]
	maxFailedNum uint8
}

// Stat ...
func (p *peerSt) Stat(id PeerId) PeerState {
	bucketItem := p.bucket.FindBlockItem(id)
	if bucketItem != nil {
		return OnlinePeerState
	}
	routeItem := p.router.FindBlockItem(id)
	if routeItem != nil {
		return UnknownPeerState
	}
	return OfflinePeerState
}

// Authenticate ...
func (p *peerSt) Authenticate(node Node, mode uint8) (peerId PeerId, err error) {

	body := bytes.NewReader(p.baseId[:])
	header := NewHeaderSegment([]byte("Mode"), []byte{mode})
	res, err := p.Request(node, body, []byte("Authenticate"), header)
	if err != nil {
		return
	}

	resBody, err := io.ReadAll(res.Body())
	if err != nil {
		return
	}

	if res.IsError() {
		err = NewReponseError(res.Code(), string(resBody))
		return
	}

	peerId, err = p.generator.Generate(resBody, node)
	if err != nil {
		return
	}

	if mode != TestOnlyAuthenticateMode {
		go p.Accept(context.Background(), node, peerId)
	}

	if mode == OpenAuthenticateMode {
		return
	}

	route := new(peerRoute)
	route.Addr = node.Addr()
	route.NodeType = node.Type()
	route.FailedNum = 0
	route.rw = new(sync.RWMutex)
	items := p.router.FindBlockItems(peerId)

	notFound := true
	if items != nil && len(items) > 0 {

		for _, item := range items {
			if item.Expired() {
				continue
			}
			r := item.Value()
			if bytes.Equal(r.Addr, route.Addr) && r.NodeType == route.NodeType {
				r.rw.Lock()
				r.FailedNum = 0
				r.rw.Unlock()
				notFound = false
			}
		}
	} else {
		notFound = false
	}
	if notFound {
		p.router.PutItem(peerId, route)
	}

	return
}

// AcceptAuthenticate ...
func (p *peerSt) AcceptAuthenticate(ctx context.Context, node Node) {

	goto AcceptAuthenticate

NextAcceptAuthenticate:

	p.AcceptAuthenticate(ctx, node)
	return

AcceptAuthenticate:

	stream, err := node.AcceptNodeStream(ctx)
	if err != nil {
		if errors.Is(err, net.ErrClosed) == false {
			node.Close()
		}
		return
	}

	authCtx, err := NewContext(stream, PeerId(p.baseId))
	if err != nil {
		_ = authCtx.ThrowError(BadRequestErrorCode, "Bad Request")
		goto NextAcceptAuthenticate
	}
	if bytes.Equal([]byte("Authenticate"), authCtx.Method()) == false {
		_ = authCtx.ThrowError(UnauthorizedErrorCode, "Unauthorized")
		goto NextAcceptAuthenticate
	}

	body, err := io.ReadAll(authCtx.Body())

	if err != nil {
		_ = authCtx.ThrowError(BadRequestErrorCode, "Bad Request")
		goto NextAcceptAuthenticate
	}

	peerId, err := p.generator.Generate(body, node)

	if err != nil {
		_ = authCtx.ThrowError(ForbiddenErrorCode, "Forbidden")
		goto NextAcceptAuthenticate
	}

	_ = authCtx.Respond(bytes.NewReader(p.baseId[:]))

	value := authCtx.Header([]byte("Mode"))
	if value != nil && value[0] != TestOnlyAuthenticateMode {
		p.Accept(ctx, node, peerId)
	}

	// TODO: try salute

}

// Open ...
func (p *peerSt) Open(peerId PeerId) (node Node, err error) {

	item := p.bucket.FindBlockItem(peerId)
	if item != nil && !item.Expired() {
		node = item.Value()
		return
	}

	ritems := p.router.FindBlockItems(peerId)
	if ritems == nil || len(ritems) <= 0 {
		err = errors.New("Not Found peer node")
		return
	}

	for _, ritem := range ritems {
		if ritem.Expired() {
			continue
		}
		route := ritem.Value()
		node, err = p.peerDialer.Connect(route.NodeType, route.Addr)
		if err == nil {
			authPeerId, err := p.Authenticate(node, OpenAuthenticateMode)
			if err == nil && bytes.Equal(authPeerId[:], peerId[:]) {
				route.rw.Lock()
				route.FailedNum = 0
				route.rw.Unlock()
				break
			}
		}

		route.rw.Lock()
		route.FailedNum++
		route.rw.Unlock()

		route.rw.RLock()
		if route.FailedNum > p.maxFailedNum {
			p.router.RemoveItem(ritem)
		}
		route.rw.RUnlock()

		_ = node.Close()
	}

	if node == nil && err == nil {
		err = errors.New("Unable to open")
	}

	return
}

// Request ...
func (p *peerSt) Request(node Node, body io.Reader, method []byte, headers ...*HeaderSegment) (res *Response, err error) {

	stream, err := node.OpenNodeStream()
	if err != nil {
		return
	}

	request := NewRequest(method, body, headers...)
	req, err := MarshalRequest(request)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer stream.Close()
		defer wg.Done()
		_, err = io.Copy(stream, req)
	}()

	var resErr error
	go func() {
		defer stream.Close()
		defer wg.Done()
		response := new(Response)
		resErr = UnmarshalResponse(stream, response)
		if resErr == nil {
			res = response
		}
	}()

	wg.Wait()
	if res != nil {
		err = nil
	} else if err == nil && resErr != nil {
		err = resErr
	}

	return
}

// AcceptServe ...
func (p *peerSt) AcceptServe(ctx context.Context, serve NodeServe) {
	for {
		node, err := serve.Accept(ctx)
		if errors.Is(err, net.ErrClosed) {
			break
		}
		if err != nil {
			panic(err)
		}

		go p.AcceptAuthenticate(ctx, node)

	}
}

// Accept
func (p *peerSt) Accept(ctx context.Context, node Node, peerId PeerId) {

	if len(peerId) <= 0 {
		panic(errors.New("Missing peer id"))
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			stream, err := node.AcceptNodeStream(ctx)
			if err != nil {
				break
			}
			go func() {
				defer stream.Close()
				c, err := NewContext(stream, peerId)
				if err == nil {
					// TODO: close with error
					_ = p.app.Run(c)
				}
			}()

		}
	}()

	item := p.bucket.PutItem(peerId, node)
	defer p.bucket.RemoveItem(item)

	wg.Wait()
}

// New ...
func New(baseId uuid.UUID, app core.App[Context], generator PeerIdGenerator, maxFailedNum uint8) Peer {

	bucket := memory.NewBucket[Node, PeerId](comparePeerId)
	router := memory.NewBucket[*peerRoute, PeerId](comparePeerId)

	dialer := new(peerDialer)
	dialer.dialerMap = make(map[uint8]NodeDialer)
	dialer.rw = new(sync.RWMutex)

	peer := new(peerSt)
	peer.baseId = baseId
	peer.app = app
	peer.generator = generator
	peer.maxFailedNum = maxFailedNum
	peer.bucket = bucket
	peer.peerDialer = dialer
	peer.router = router

	return peer
}
