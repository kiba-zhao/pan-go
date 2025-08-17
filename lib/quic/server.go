// Define peer server for quic
package quic

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"pan/lib/log"
	libNet "pan/lib/net"
	"strconv"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

var ErrQuicServerUnavailable = errors.New("quic.QuicServer Error: Unavailable")

type stdQuicServer struct {
	logger  log.Logger
	cluster *stdQuicCluster

	transport   *quic.Transport
	transportRW sync.RWMutex

	port   uint16
	portRW sync.RWMutex

	certificate   tls.Certificate
	certificateRW sync.RWMutex

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool
}

func (qs *stdQuicServer) Transport() *quic.Transport {
	qs.transportRW.RLock()
	defer qs.transportRW.RUnlock()
	return qs.transport
}

func (qs *stdQuicServer) setTransport(transport *quic.Transport) {
	qs.transportRW.Lock()
	defer qs.transportRW.Unlock()
	qs.transport = transport
}

func (qs *stdQuicServer) Port() uint16 {
	qs.portRW.RLock()
	defer qs.portRW.RUnlock()
	return qs.port
}

func (qs *stdQuicServer) SetPort(port uint16) {
	qs.logger.Debug("QuicServer", "SetPort")

	qs.portRW.Lock()
	defer qs.portRW.Unlock()
	if qs.port == port {
		return
	}
	qs.port = port
	qs.Reload()
}

func (qs *stdQuicServer) Certificate() tls.Certificate {
	qs.certificateRW.RLock()
	defer qs.certificateRW.RUnlock()
	return qs.certificate
}

func (qs *stdQuicServer) SetCertificate(certificate tls.Certificate) {
	qs.certificateRW.Lock()
	defer qs.certificateRW.Unlock()
	qs.certificate = certificate
	qs.Reload()
}

func (qs *stdQuicServer) Reload() {
	qs.logger.Debug("QuicServer", "Reload")

	qs.reloadLock.Lock()
	defer qs.reloadLock.Unlock()
	if qs.reload {
		return
	}

	qs.reload = true
	qs.reloadChan <- struct{}{}
}

func (qs *stdQuicServer) ListenAndServe(ctx context.Context) error {
	qs.logger.Debug("QuicServer", "ListenAndServe begin")
	defer qs.logger.Debug("QuicServer", "ListenAndServe end")

	var err error
	var closed bool

	var wg sync.WaitGroup
	var timer <-chan time.Time

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-qs.reloadChan:
			if timer != nil {
				<-timer
				timer = nil
			}
			qs.reloadLock.Lock()
			qs.reload = false
			qs.reloadLock.Unlock()
		}

		server := qs.Transport()
		if server != nil {
			server.Close()
			qs.setTransport(nil)
			wg.Wait()
		}

		if closed {
			break
		}

		certificate := qs.Certificate()
		port := qs.port
		if port <= 0 || len(certificate.Certificate) <= 0 {
			continue
		}

		addrStat, addrStatErr := libNet.StatAddr()
		if addrStatErr != nil {
			timer = time.After(time.Second * 5)
			qs.Reload()
			continue
		}

		var addr string
		if addrStat.IPv6Enabled {
			addr = "[::]:" + strconv.FormatUint(uint64(port), 10)
		} else {
			addr = "0.0.0.0:" + strconv.FormatUint(uint64(port), 10)
		}

		cluster := qs.cluster
		tlsConf := &tls.Config{ClientAuth: tls.RequireAnyClientCert, Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
		quicConf := &quic.Config{}

		udpAddr, udpAddrErr := net.ResolveUDPAddr("udp", addr)
		if udpAddrErr != nil {
			qs.logger.Error("QuicServer", "net.ResolveUDPAddr Error: "+addr)
			continue
		}

		udpConn, udpConnErr := net.ListenUDP("udp", udpAddr)
		if udpConnErr != nil {
			qs.logger.Error("QuicServer", "net.ListenUDP Error: "+addr)
			continue
		}

		tr := quic.Transport{
			Conn: udpConn,
		}
		qs.setTransport(&tr)

		ln, lnErr := tr.Listen(tlsConf, quicConf)
		if lnErr != nil {
			qs.logger.Error("QuicServer", "quic.ListenAddr Error: "+addr)
			continue
		} else {
			qs.logger.Info("QuicServer", "quic.ListenAddr Success: "+addr)
		}

		wg.Add(1)
		go func(ln *quic.Listener, udpConn *net.UDPConn) {
			defer wg.Done()
			defer udpConn.Close()
			for {
				conn, err := ln.Accept(ctx)
				if err != nil && conn != nil {
					conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
				}
				if err == nil {
					_, err = cluster.Serve(conn, nil)
				}
				if errors.Is(err, quic.ErrServerClosed) || errors.Is(err, context.Canceled) {
					break
				}
				if ctxErr := ctx.Err(); ctxErr != nil {
					break
				}
				qs.logger.Error("QuicServer", "quic.Listener.Accept Error: "+err.Error())
			}
		}(ln, udpConn)

	}

	return err
}
