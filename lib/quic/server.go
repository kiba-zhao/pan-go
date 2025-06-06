// Define peer server for quic
package quic

import (
	"context"
	"crypto/tls"
	"errors"
	"pan/lib/log"
	"slices"
	"sync"

	"github.com/quic-go/quic-go"
)

var ErrQuicServerUnavailable = errors.New("quic.QuicServer Error: Unavailable")

type stdQuicServer struct {
	logger  log.Logger
	cluster *stdQuicCluster

	addrs   []string
	addrsRW sync.RWMutex

	certificate   tls.Certificate
	certificateRW sync.RWMutex

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool
}

func (qs *stdQuicServer) Addrs() []string {
	qs.addrsRW.RLock()
	defer qs.addrsRW.RUnlock()
	return qs.addrs
}

func (qs *stdQuicServer) SetAddrs(addrs []string) {
	qs.logger.Debug("QuicServer", "SetAddrs")

	qs.addrsRW.Lock()
	defer qs.addrsRW.Unlock()
	if slices.Equal(qs.addrs, addrs) {
		return
	}
	qs.addrs = addrs
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
	var servers []*quic.Listener

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-qs.reloadChan:
			qs.reloadLock.Lock()
			qs.reload = false
			qs.reloadLock.Unlock()
		}

		if len(servers) > 0 {
			for _, server := range servers {
				server.Close()
			}
			wg.Wait()
		}

		if closed {
			break
		}

		certificate := qs.Certificate()
		addrs := qs.Addrs()
		if len(addrs) <= 0 || len(certificate.Certificate) <= 0 {
			continue
		}

		servers = make([]*quic.Listener, 0)
		cluster := qs.cluster
		tlsConf := &tls.Config{ClientAuth: tls.RequireAnyClientCert, Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
		quicConf := &quic.Config{}
		for _, addr := range addrs {
			ln, lnErr := quic.ListenAddr(addr, tlsConf, quicConf)
			if lnErr != nil {
				qs.logger.Error("QuicServer", "quic.ListenAddr Error: "+addr)
				continue
			} else {
				qs.logger.Info("QuicServer", "quic.ListenAddr Success: "+addr)
			}

			servers = append(servers, ln)
			wg.Add(1)
			go func(ln *quic.Listener) {
				defer wg.Done()
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
			}(ln)

		}

	}
	return err
}
