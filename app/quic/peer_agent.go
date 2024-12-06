package quic

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

const (
	QuicConnFlagSync = uint8(iota)
	QuicConnFlagGreet
	QuicConnFlagInvite
	QuicConnFlagDeclineInvite
)

func acceptQuicConn(conn QuicConn, agent *quicPeerAgent, ch chan struct{}) error {
	ch <- struct{}{}
	var err error
	for {
		stream, err := conn.AcceptUniStream(context.Background())
		if err != nil {
			break
		}

		// read flag
		flags := make([]byte, 1)
		_, err = stream.Read(flags)
		if err != nil {
			break
		}
		//

		switch flags[0] {
		case QuicConnFlagSync:
			go agent.AcceptSync(stream, conn)
		case QuicConnFlagGreet:
			go agent.AcceptGreet(stream, conn)
		case QuicConnFlagInvite:
			go agent.AcceptInvite(stream, conn)
		case QuicConnFlagDeclineInvite:
			go agent.AcceptDeclineInvite(stream, conn)
		}

	}
	return err
}

func doQuicConn(conn QuicConn, flag uint8, reader io.Reader) error {
	var dataReader io.Reader
	if reader == nil {
		dataReader = reader
	} else {
		dataReader = io.MultiReader(bytes.NewReader([]byte{flag}), reader)
	}

	stream, err := conn.OpenUniStream()
	if err != nil {
		return err
	}
	defer stream.Close()

	_, err = io.Copy(stream, dataReader)
	return err
}

type quicPeerAgent struct {
	quicPeerModule QuicPeerModule
	matrix         [][]*quicReception
	locker         sync.Mutex
}

func (agent *quicPeerAgent) Follow(conn QuicConn) error {

	ch := make(chan struct{})
	go acceptQuicConn(conn, agent, ch)
	<-ch

	agent.locker.Lock()
	receptions := searchReceptions(agent.matrix, conn.PeerID())
	if len(receptions) <= 0 {
		agent.locker.Unlock()
		return nil
	}
	agent.matrix = deleteReception(agent.matrix, receptions[0])
	agent.locker.Unlock()
	receptions[0].ch <- conn

	return nil
}

func (agent *quicPeerAgent) Sync(conn QuicConn) error {
	// TODO: implement

	return nil
}

func (agent *quicPeerAgent) AcceptSync(stream quic.ReceiveStream, conn QuicConn) error {
	defer stream.CancelRead(quic.StreamErrorCode(quic.NoError))

	// TODO: implement
	return nil
}

func (agent *quicPeerAgent) Greet(conn QuicConn) error {
	addrs := agent.quicPeerModule.PublicAddrs()
	if len(addrs) <= 0 {
		return nil
	}

	greetMsg := QuicGreet{
		Addrs: addrs,
	}
	data, err := proto.Marshal(&greetMsg)
	if err != nil {
		return err
	}

	return doQuicConn(conn, QuicConnFlagGreet, bytes.NewReader(data))
}

func (agent *quicPeerAgent) AcceptGreet(stream quic.ReceiveStream, conn QuicConn) error {
	defer stream.CancelRead(quic.StreamErrorCode(quic.NoError))

	data, err := io.ReadAll(stream)
	if err != nil {
		return err
	}

	var greetMsg QuicGreet
	err = proto.Unmarshal(data, &greetMsg)
	if err == nil {

		remoteAddr := conn.RemoteAddr()
		ip, _, remoteAddrErr := net.SplitHostPort(remoteAddr.String())
		if remoteAddrErr != nil {
			return remoteAddrErr
		}

		for _, addr := range greetMsg.Addrs {

			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				continue
			}

			if ip != host {
				ipAddr, err := net.ResolveIPAddr("ip", host)
				if err != nil || !ipAddr.IP.IsUnspecified() {
					continue
				}
				addr = net.JoinHostPort(ip, port)
			}

			agent.quicPeerModule.Route(conn.PeerID(), addr, false)
		}
	}

	return err
}

func (agent *quicPeerAgent) cancelInvite(reception *quicReception) {
	agent.locker.Lock()
	defer agent.locker.Unlock()

	agent.matrix = deleteReception(agent.matrix, reception)
}

func (agent *quicPeerAgent) Invite(conn QuicConn) (QuicConn, error) {

	var reception quicReception
	reception.ch = make(chan QuicConn, 2)
	defer close(reception.ch)

	agent.locker.Lock()
	agent.matrix = storeReception(agent.matrix, &reception)
	agent.locker.Unlock()

	err := doQuicConn(conn, QuicConnFlagInvite, nil)
	if err != nil {
		agent.cancelInvite(&reception)
		return nil, err
	}

	var conn_ QuicConn
	select {
	case conn_ = <-reception.ch:
	case <-time.After(6 * time.Second):
		agent.cancelInvite(&reception)
		err = errors.New("invite timeout")
	}

	return conn_, err
}

func (agent *quicPeerAgent) AcceptInvite(stream quic.ReceiveStream, conn QuicConn) error {
	defer stream.CancelRead(quic.StreamErrorCode(quic.NoError))
	_, err := agent.quicPeerModule.Dial(context.Background(), conn.PeerID())
	if err != nil {
		agent.DeclineInvite(conn)
	}
	return err
}

func (agent *quicPeerAgent) DeclineInvite(conn QuicConn) error {
	return doQuicConn(conn, QuicConnFlagDeclineInvite, nil)
}

func (agent *quicPeerAgent) AcceptDeclineInvite(stream quic.ReceiveStream, conn QuicConn) error {
	defer stream.CancelRead(quic.StreamErrorCode(quic.NoError))
	agent.locker.Lock()
	matrix, receptions := takeOutReception(agent.matrix, conn.PeerID())
	agent.matrix = matrix
	agent.locker.Unlock()

	for _, reception := range receptions {
		reception.ch <- nil
	}
	return nil
}
