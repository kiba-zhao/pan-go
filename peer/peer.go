package peer

import (
	"bytes"
	"context"
	"errors"

	"io"
	"net"
	"pan/core"
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

type PeerId uuid.UUID

type Peer interface {
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

type peerPassport struct {
	id       []byte
	isPeerId bool
}

type SimplePeerIdGenerator struct {
	passports   []*peerPassport
	defaultDeny bool
	rw          *sync.RWMutex
}

// Contains ...
func (pg *SimplePeerIdGenerator) Contains(id []byte, isPeerId bool) bool {
	return false
}

// AddItem ...
func (pg *SimplePeerIdGenerator) Add(id []byte, isPeerId bool) error {

	if pg.Contains(id, isPeerId) {
		return errors.New("Duplicate peer id item")
	}

	passport := new(peerPassport)
	passport.id = id
	passport.isPeerId = isPeerId

	pg.rw.Lock()
	pg.passports = append(pg.passports, passport)
	pg.rw.Unlock()

	return nil
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
	idBytes := id[:]
	pass := !pg.defaultDeny

	pg.rw.RLock()
	for _, passport := range pg.passports {
		if passport.isPeerId == true && bytes.Equal(passport.id, idBytes) {
			pass = pg.defaultDeny
			break
		}
		if passport.isPeerId == false && bytes.Equal(passport.id, baseId) {
			pass = pg.defaultDeny
			break
		}
	}
	pg.rw.RUnlock()

	if !pass {
		err = errors.New("Deny Peer Id")
		return
	}

	peerId = PeerId(id)
	return
}

// NewPeerIdGenerator ...
func NewPeerIdGenerator(defaultDeny bool) *SimplePeerIdGenerator {
	generator := new(SimplePeerIdGenerator)
	generator.passports = make([]*peerPassport, 0)
	generator.defaultDeny = defaultDeny
	generator.rw = new(sync.RWMutex)

	return generator
}

type peerNodeItem struct {
	peerId PeerId
	node   Node
	idx    int
}

type peerBucket struct {
	items []*peerNodeItem
	rw    *sync.RWMutex
}

func (pb *peerBucket) selectOne(peerId PeerId) (node Node) {

	pb.rw.RLock()

	for idx := len(pb.items) - 1; idx >= 0; idx-- {
		item := pb.items[idx]
		if bytes.Equal(peerId[:], item.peerId[:]) {
			node = item.node
			break
		}
	}

	pb.rw.RUnlock()
	return
}

func (pb *peerBucket) add(node Node, peerId PeerId) *peerNodeItem {

	item := new(peerNodeItem)
	item.peerId = peerId
	item.node = node

	pb.rw.Lock()

	item.idx = len(pb.items)
	pb.items = append(pb.items, item)

	pb.rw.Unlock()

	return item
}

// delete ...
func (pb *peerBucket) del(item *peerNodeItem) {

	pb.rw.Lock()

	lastIdx := len(pb.items) - 1
	if lastIdx != item.idx {
		lastItem := pb.items[lastIdx]
		lastItem.idx = item.idx
		pb.items[item.idx] = lastItem
	}
	pb.items = pb.items[:lastIdx]

	pb.rw.Unlock()
}

type peerRoute struct {
	PeerId    PeerId
	NodeType  uint8
	Addr      []byte
	FailedNum uint8
	idx       int
}

type peerRouter struct {
	routes []*peerRoute
	rw     *sync.RWMutex
}

// find ...
func (pr *peerRouter) find(peerId PeerId) (routes []*peerRoute) {
	if len(pr.routes) <= 0 {
		return
	}
	pr.rw.RLock()
	rs := make([]*peerRoute, 0)
	for _, route := range pr.routes {
		if bytes.Equal(route.PeerId[:], peerId[:]) {
			rs = append(rs, route)
			break
		}
	}
	pr.rw.RUnlock()

	if len(rs) > 0 {
		routes = rs
	}
	return
}

// findOne ...
func (pr *peerRouter) findOne(peerId PeerId, node Node) (route *peerRoute) {

	pr.rw.RLock()
	for _, r := range pr.routes {
		if bytes.Equal(r.PeerId[:], peerId[:]) && r.NodeType == node.Type() && bytes.Equal(r.Addr, node.Addr()) {
			route = r
			break
		}
	}
	pr.rw.RUnlock()
	return
}

// addRoute ...
func (pr *peerRouter) addRoute(peerId PeerId, node Node) *peerRoute {

	route := pr.findOne(peerId, node)
	if route != nil {
		return nil
	}
	route = new(peerRoute)
	route.PeerId = peerId
	route.NodeType = node.Type()
	route.Addr = node.Addr()
	route.FailedNum = 0

	pr.rw.Lock()

	route.idx = len(pr.routes)
	pr.routes = append(pr.routes, route)
	pr.rw.Unlock()

	return route
}

// removeRoute ...
func (pr *peerRouter) removeRoute(route *peerRoute) {

	pr.rw.Lock()
	lastIdx := len(pr.routes) - 1
	if lastIdx != route.idx {
		lastRoute := pr.routes[lastIdx]
		lastRoute.idx = route.idx
		pr.routes[route.idx] = lastRoute
	}
	pr.routes = pr.routes[:lastIdx]
	pr.rw.Unlock()
}

type peerDialer struct {
	dialerMap map[uint8]NodeDialer
	rw        *sync.RWMutex
}

// Connect ...
func (pd *peerDialer) connect(dialerType uint8, addr []byte) (node Node, err error) {
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
	router       *peerRouter
	bucket       *peerBucket
	baseId       []byte
	app          core.App[Context]
	maxFailedNum uint8
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

	if mode != OpenAuthenticateMode {
		p.router.addRoute(peerId, node)
	}

	if mode == NormalAuthenticateMode || mode == OpenAuthenticateMode {
		go p.Accept(context.Background(), node, peerId)
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
	if value != nil && value[0] == NormalAuthenticateMode {
		p.Accept(ctx, node, peerId)
	}

}

// Open ...
func (p *peerSt) Open(peerId PeerId) (node Node, err error) {

	node = p.bucket.selectOne(peerId)
	if node == nil {
		return
	}

	routes := p.router.find(peerId)
	if routes == nil {
		return
	}
	if len(routes) <= 0 {
		err = errors.New("Not Found peer node")
		return
	}

	for _, route := range routes {
		node, err = p.peerDialer.connect(route.NodeType, route.Addr)
		if err == nil {
			peerId, err := p.Authenticate(node, OpenAuthenticateMode)
			if err == nil && bytes.Equal(route.PeerId[:], peerId[:]) {
				route.FailedNum = 0
				break
			}
		}

		route.FailedNum++
		if route.FailedNum > p.maxFailedNum {
			p.router.removeRoute(route)
		}
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
	item := p.bucket.add(node, peerId)
	defer p.bucket.del(item)

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
				p.app.Run(c)
			}
		}()

	}

}

// New ...
func New(baseId []byte, app core.App[Context], generator PeerIdGenerator, maxFailedNum uint8) Peer {
	bucket := new(peerBucket)
	bucket.items = make([]*peerNodeItem, 0)
	bucket.rw = new(sync.RWMutex)

	router := new(peerRouter)
	router.routes = make([]*peerRoute, 0)
	router.rw = new(sync.RWMutex)

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
