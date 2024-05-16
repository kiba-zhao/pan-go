package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"pan/runtime"
	"reflect"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type NodeApp = *gin.Engine

type NodeAppProvider interface {
	NodeApp() NodeApp
}

type NodeModule interface {
	SetupToNode(app NodeApp) error
}

type NodeModuleProvider interface {
	NodeModules() []NodeModule
}

func NewNodeApp() NodeApp {
	return gin.New()
}

type nodeEngine struct {
	modules []interface{}
	once    sync.Once
}

func (n *nodeEngine) Init(registry runtime.Registry) error {

	providers := runtime.ModulesForType[NodeAppProvider](registry)
	if len(providers) <= 0 {
		return errors.New("no providers for NodeAppProvider")
	}
	err := runtime.TraverseRegistry(registry, func(module NodeModule) error {
		var setupErr error
		for _, provider := range providers {
			setupErr = module.SetupToNode(provider.NodeApp())
			if setupErr != nil {
				break
			}
		}
		return setupErr
	})
	if err != nil {
		return err
	}

	err = runtime.TraverseRegistry(registry, func(module NodeModuleProvider) error {
		var setupErr error
		for _, m := range module.NodeModules() {
			for _, provider := range providers {
				setupErr = m.SetupToNode(provider.NodeApp())
				if err != nil {
					break
				}
			}
		}
		return setupErr
	})
	return err
}

func (n *nodeEngine) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[NodeModule](),
		reflect.TypeFor[NodeModuleProvider](),
		reflect.TypeFor[NodeAppProvider](),
	}
}

func (n *nodeEngine) Modules() []interface{} {
	n.once.Do(func() {
		n.modules = append(n.modules,
			&nodeServer{},
			&nodeClient{},
		)
	})
	return n.modules
}

type NodeConnection = quic.Connection

type NodeServeHandler interface {
	ServeConn(conn NodeConnection) error
}

type nodeServer struct {
	Config   AppConfig
	once     sync.Once
	h3Server http3.Server
	handlers []NodeServeHandler
	rw       sync.RWMutex
}

func (n *nodeServer) NodeApp() NodeApp {
	return n.h3Server.Handler.(NodeApp)
}

func (n *nodeServer) Init(registry runtime.Registry) error {
	n.rw.Lock()
	defer n.rw.Unlock()
	n.handlers = runtime.ModulesForType[NodeServeHandler](registry)
	return nil
}

func (n *nodeServer) ServeConn(conn NodeConnection) error {
	return n.h3Server.ServeQUICConn(conn)
}

func (n *nodeServer) onAccept(conn quic.Connection) error {
	n.rw.RLock()
	defer n.rw.RUnlock()

	var err error
	for _, handler := range n.handlers {
		err = handler.ServeConn(conn)
		if err != nil {
			break
		}
	}
	return err
}

func (n *nodeServer) Ready() error {
	settings, err := n.Config.Read()
	if err != nil {
		return err
	}

	certificate, err := tls.LoadX509KeyPair(settings.CertificatePath, settings.PrivateKeyPath)
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", settings.Host+":"+strconv.Itoa(settings.Port))
	if err != nil {
		return err
	}

	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	tr := quic.Transport{Conn: udpConn}
	tlsConf := http3.ConfigureTLSConfig(&tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13})
	quicConf := &quic.Config{}
	ln, _ := tr.ListenEarly(tlsConf, quicConf)
	ctx := context.Background()
	for {

		conn, err := ln.Accept(ctx)
		if err == quic.ErrServerClosed {
			break
		}
		if err != nil {
			go conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
			continue
		}

		switch conn.ConnectionState().TLS.NegotiatedProtocol {
		case http3.NextProtoH3:
			go n.onAccept(conn)
		default:
			go conn.CloseWithError(quic.ApplicationErrorCode(quic.ProtocolViolation), "unknown protocol")
		}
	}
	return err
}

func (n *nodeServer) Components() []runtime.Component {

	n.once.Do(func() {
		n.h3Server = http3.Server{Handler: NewNodeApp()}
	})

	return []runtime.Component{
		runtime.NewComponent(n, runtime.ComponentNoneScope),
	}
}

type NodeDialHandler interface {
	DialConn(conn NodeConnection) error
}

type nodeClient struct {
	handlers []NodeDialHandler
	rw       sync.RWMutex
}

func (n *nodeClient) Init(registry runtime.Registry) error {
	n.rw.Lock()
	defer n.rw.Unlock()
	n.handlers = runtime.ModulesForType[NodeDialHandler](registry)
	return nil
}

func (n *nodeClient) DialConn(conn NodeConnection) error {
	return nil
}

func extractPublicKeyFromConn(conn NodeConnection) ([]byte, error) {
	state := conn.ConnectionState()
	cert := state.TLS.PeerCertificates[0]
	return x509.MarshalPKIXPublicKey(cert.PublicKey)
}
