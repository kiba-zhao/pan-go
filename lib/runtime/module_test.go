package runtime_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"pan/lib/peer"
	"pan/lib/runtime"
	"sync"
	"testing"

	"github.com/quic-go/quic-go"
	"github.com/stretchr/testify/assert"
)

func TestSimpleModule(t *testing.T) {

	setup := func() (engine *runtime.Engine) {
		engine, _ = runtime.New()
		return
	}

	t.Run("simple module", func(t *testing.T) {
		e := setup()

		module := &TestModule{}
		simpleModule := runtime.NewModule(module)
		err := e.Mount(simpleModule)
		assert.Nil(t, err)
	})

	t.Run("test", func(t *testing.T) {
		peerSettings, err := peer.GeneratePeerSettings()
		if err != nil {
			t.Fatal(err)
		}

		settings, err := peer.GeneratePeerSettings()
		if err != nil {
			t.Fatal(err)
		}
		servtlsConf := &tls.Config{ClientAuth: tls.RequireAnyClientCert, Certificates: []tls.Certificate{settings.Certificate()}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
		serverUdpConn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 9001})
		if err != nil {
			t.Fatal(err)
		}
		defer serverUdpConn.Close()
		serverTr := &quic.Transport{Conn: serverUdpConn}
		server, err := serverTr.Listen(servtlsConf, &quic.Config{})
		if err != nil {
			t.Fatal(err)
		}
		defer server.Close()
		var wg sync.WaitGroup
		wg.Add(1)
		go func(server *quic.Listener, serverTr *quic.Transport) {
			wg.Done()
			conn, err := server.Accept(context.Background())
			if err == nil {
				defer conn.CloseWithError(quic.ApplicationErrorCode(0), "")
				fmt.Println("server accept :", conn.RemoteAddr().String())
				servtlsConf := &tls.Config{Certificates: []tls.Certificate{settings.Certificate()}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
				serverConn, err := serverTr.Dial(conn.Context(), conn.RemoteAddr(), servtlsConf, &quic.Config{})
				if err == nil {
					defer serverConn.CloseWithError(quic.ApplicationErrorCode(0), "")
					fmt.Println("server conn accept")
					stream, err := conn.OpenStream()
					if err == nil {
						defer stream.Close()
						stream.Write([]byte("hello"))
						fmt.Println("server stream open", stream.StreamID())
					}
				}

			}

			<-time.After(time.Second * 10)
			fmt.Println("server accept", err)
		}(server, serverTr)

		udpConn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 9000})
		tlsConf := &tls.Config{ClientAuth: tls.RequireAnyClientCert, Certificates: []tls.Certificate{peerSettings.Certificate()}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
		tr := &quic.Transport{Conn: udpConn}
		ln, err := tr.Listen(tlsConf, &quic.Config{})
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		wg.Add(1)
		go func(ln *quic.Listener) {
			wg.Done()
			conn, err := ln.Accept(context.Background())
			if err == nil {
				fmt.Println("ln accept :", conn.RemoteAddr().String())
				conn.CloseWithError(quic.ApplicationErrorCode(0), "")
			}
			fmt.Println("ln accept", err)
		}(ln)
		wg.Wait()

		clienttlsConf := &tls.Config{Certificates: []tls.Certificate{peerSettings.Certificate()}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
		conn, err := tr.Dial(context.Background(), &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9001}, clienttlsConf, &quic.Config{})
		if err != nil {
			t.Fatal(err)
		}
		defer conn.CloseWithError(quic.ApplicationErrorCode(0), "")

		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*10))
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer stream.Close()
		fmt.Println("conn stream accept", stream.StreamID())
		<-time.After(time.Second * 10)
	})
}
