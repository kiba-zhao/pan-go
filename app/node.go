package app

import (
	"io"
	"pan/app/node"
	"pan/runtime"
	"reflect"
	"sync"
)

type NodeApp = *node.App
type NodeRouter = node.AppHandleGroup

type NodeModule interface {
	SetupToNode(NodeRouter) error
}

type NodeModuleProvider interface {
	NodeModules() []NodeModule
}

type NodeStreamCloser interface {
	io.Closer
	CloseRead() error
	CloseWrite() error
}

type NodeStream interface {
	io.Reader
	io.Writer
	NodeStreamCloser
}

type NodeID = []byte
type Node interface {
	ID() NodeID
}

type NodeServer interface {
	Serve(NodeStream, Node) error
}

var (
	ContextNode = []byte("NODE")
)

type nodeServer struct {
	app NodeApp
}

func (n *nodeServer) Serve(stream NodeStream, target Node) error {
	ctx := node.NewAppContext()
	ctx.Set(ContextNode, target)
	err := node.UnmarshalRequest(stream, ctx.Request())
	if err == nil {
		err = n.app.Run(ctx, nil)
	}

	if err != nil {
		ctx.ThrowError(CodeInternalError, err)
	}

	if ctx.Code() < 0 {
		ctx.ThrowError(CodeNotFound, ErrNotFound)
	}

	reader := node.MarshalResponse(&ctx.Response)
	_, resErr := io.Copy(stream, reader)
	if resErr == nil {
		resErr = stream.Close()
	}

	if err == nil && resErr != nil {
		err = resErr
	}
	return err
}

type nodeApp struct {
	app  NodeApp
	once sync.Once
}

func (n *nodeApp) Init(registry runtime.Registry) error {

	app := n.app
	err := runtime.TraverseRegistry(registry, func(module NodeModule) error {
		return module.SetupToNode(app)
	})
	if err != nil {
		return err
	}

	return runtime.TraverseRegistry(registry, func(module NodeModuleProvider) error {
		var setupErr error
		for _, m := range module.NodeModules() {
			setupErr = m.SetupToNode(app)
			if err != nil {
				break
			}
		}
		return setupErr
	})
}

func (n *nodeApp) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[NodeModule](),
		reflect.TypeFor[NodeModuleProvider](),
	}
}

func (n *nodeApp) Components() []runtime.Component {
	n.once.Do(func() {
		n.app = node.NewApp()
	})

	server := &nodeServer{app: n.app}
	return []runtime.Component{
		runtime.NewComponent[NodeServer](server, runtime.ComponentExternalScope),
	}
}

func NewNodeModule() interface{} {
	return runtime.NewModule(&nodeApp{})
}

// type NodeConnection = quic.Connection

// type NodeServeHandler interface {
// 	ServeConn(conn NodeConnection) error
// }

// type nodeServer struct {
// 	Config   AppConfig
// 	once     sync.Once
// 	h3Server http3.Server
// 	handlers []NodeServeHandler
// 	rw       sync.RWMutex
// }

// func (n *nodeServer) NodeApp() NodeApp {
// 	return n.h3Server.Handler.(NodeApp)
// }

// func (n *nodeServer) Init(registry runtime.Registry) error {
// 	n.rw.Lock()
// 	defer n.rw.Unlock()
// 	n.handlers = runtime.ModulesForType[NodeServeHandler](registry)
// 	return nil
// }

// func (n *nodeServer) ServeConn(conn NodeConnection) error {
// 	return n.h3Server.ServeQUICConn(conn)
// }

// func (n *nodeServer) onAccept(conn quic.Connection) error {
// 	n.rw.RLock()
// 	defer n.rw.RUnlock()

// 	var err error
// 	for _, handler := range n.handlers {
// 		err = handler.ServeConn(conn)
// 		if err != nil {
// 			break
// 		}
// 	}
// 	return err
// }

// func (n *nodeServer) Ready() error {
// 	settings, err := n.Config.Read()
// 	if err != nil {
// 		return err
// 	}

// 	certificate, err := tls.LoadX509KeyPair(settings.CertificatePath, settings.PrivateKeyPath)
// 	if err != nil {
// 		return err
// 	}

// 	addr, err := net.ResolveUDPAddr("udp", settings.Host+":"+strconv.Itoa(settings.Port))
// 	if err != nil {
// 		return err
// 	}

// 	udpConn, err := net.ListenUDP("udp", addr)
// 	if err != nil {
// 		return err
// 	}

// 	tr := quic.Transport{Conn: udpConn}
// 	tlsConf := http3.ConfigureTLSConfig(&tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13})
// 	quicConf := &quic.Config{}
// 	ln, _ := tr.ListenEarly(tlsConf, quicConf)
// 	ctx := context.Background()
// 	for {

// 		conn, err := ln.Accept(ctx)
// 		if err == quic.ErrServerClosed {
// 			break
// 		}
// 		if err != nil {
// 			go conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
// 			continue
// 		}

// 		switch conn.ConnectionState().TLS.NegotiatedProtocol {
// 		case http3.NextProtoH3:
// 			go n.onAccept(conn)
// 		default:
// 			go conn.CloseWithError(quic.ApplicationErrorCode(quic.ProtocolViolation), "unknown protocol")
// 		}
// 	}
// 	return err
// }

// func (n *nodeServer) Components() []runtime.Component {

// 	n.once.Do(func() {
// 		n.h3Server = http3.Server{Handler: NewNodeApp()}
// 	})

// 	return []runtime.Component{
// 		runtime.NewComponent(n, runtime.ComponentNoneScope),
// 	}
// }

// type NodeDialHandler interface {
// 	DialConn(conn NodeConnection) error
// }

// type nodeClient struct {
// 	handlers []NodeDialHandler
// 	rw       sync.RWMutex
// }

// func (n *nodeClient) Init(registry runtime.Registry) error {
// 	n.rw.Lock()
// 	defer n.rw.Unlock()
// 	n.handlers = runtime.ModulesForType[NodeDialHandler](registry)
// 	return nil
// }

// func (n *nodeClient) DialConn(conn NodeConnection) error {
// 	return nil
// }

// func extractPublicKeyFromConn(conn NodeConnection) ([]byte, error) {
// 	state := conn.ConnectionState()
// 	cert := state.TLS.PeerCertificates[0]
// 	return x509.MarshalPKIXPublicKey(cert.PublicKey)
// }
