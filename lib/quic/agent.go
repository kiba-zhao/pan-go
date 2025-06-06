// Define peer agent for quic
//
// It is used to manage the connections and streams between the local node and the remote node.
package quic

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"pan/lib/log"
	"slices"

	"sync"
	"time"

	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

const (
	QuicConnFlagGreet = uint8(iota)
	QuicConnFlagInvite
	QuicConnFlagDeclineInvite
)

func acceptQuicConn(conn QuicConn, agent *stdQuicAgent) error {
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
		dataReader = bytes.NewReader([]byte{flag})
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

type stdQuicAgent struct {
	logger log.Logger

	cluster *stdQuicCluster
	matrix  [][]*stdQuicReception
	locker  sync.Mutex

	publicAddrs   []string
	publicAddrsRW sync.RWMutex
}

func (agent *stdQuicAgent) SetPublicAddrs(addrs []string) {
	agent.logger.Debug("QuicAgent", "SetPublicAddrs")

	agent.publicAddrsRW.Lock()
	defer agent.publicAddrsRW.Unlock()
	if slices.Equal(agent.publicAddrs, addrs) {
		return
	}
	agent.publicAddrs = addrs
}

func (agent *stdQuicAgent) PublicAddrs() []string {
	agent.publicAddrsRW.RLock()
	defer agent.publicAddrsRW.RUnlock()
	return agent.publicAddrs
}

// Follow processes the given QuicConn by searching for associated receptions
// in the agent's matrix using the connection's peer ID. If any receptions are
// found, it deletes the first reception from the matrix and sends the connection
// through the reception's channel. The function then spawns a goroutine to
// handle the connection with acceptQuicConn. The method acquires a lock to
// ensure thread-safe access to the matrix.

func (agent *stdQuicAgent) Follow(conn QuicConn) error {

	agent.locker.Lock()
	receptions := searchReceptions(agent.matrix, conn.PeerID())
	if len(receptions) > 0 {
		agent.matrix = deleteReception(agent.matrix, receptions[0])
		receptions[0].ch <- conn
	}
	agent.locker.Unlock()

	go acceptQuicConn(conn, agent)
	return nil
}

// Greet sends a QuicConnFlagGreet message over the given connection
// containing the public addresses of the local node. If the local node
// has no public addresses, the function returns immediately. The method
// simply marshals a QuicGreet message and calls doQuicConn with the result.
func (agent *stdQuicAgent) Greet(conn QuicConn) error {
	addrs := agent.PublicAddrs()
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

// AcceptGreet handles a QuicConnFlagGreet message sent by the peer, by unmarshaling
// the message and adding the remote addresses to the peer routing table. The method
// returns an error if the message unmarshaling fails.
func (agent *stdQuicAgent) AcceptGreet(stream quic.ReceiveStream, conn QuicConn) error {
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

			agent.cluster.Route(conn.PeerID(), addr)
		}
	}

	return err
}

func (agent *stdQuicAgent) cancelInvite(reception *stdQuicReception) {
	agent.locker.Lock()
	defer agent.locker.Unlock()

	agent.matrix = deleteReception(agent.matrix, reception)
}

// Invite sends a QuicConnFlagInvite message over the given connection.
// The message signals that the local node is willing to connect to the
// remote node. The method then waits for a QuicConnFlagDeclineInvite or
// QuicConnFlagGreet message from the remote node. If a QuicConnFlagGreet
// message is received, the method returns a new connection to the remote
// node. If a QuicConnFlagDeclineInvite message is received or if the
// remote node does not respond within 6 seconds, the method returns an
// error.
func (agent *stdQuicAgent) Invite(conn QuicConn) (QuicConn, error) {

	var reception stdQuicReception
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

// AcceptInvite handles a QuicConnFlagInvite message sent by the peer, by
// dialing to the peer using the quicPeerModule. If the dialing fails, the
// method sends a QuicConnFlagDeclineInvite message to the peer and returns
// the error.
func (agent *stdQuicAgent) AcceptInvite(stream quic.ReceiveStream, conn QuicConn) error {
	defer stream.CancelRead(quic.StreamErrorCode(quic.NoError))
	_, err := agent.cluster.Dial(context.Background(), conn.PeerID())
	if err != nil {
		agent.DeclineInvite(conn)
	}
	return err
}

// DeclineInvite sends a QuicConnFlagDeclineInvite message over the given connection.
// This message tells the remote node that the local node is not willing to connect
// to the remote node.
func (agent *stdQuicAgent) DeclineInvite(conn QuicConn) error {
	return doQuicConn(conn, QuicConnFlagDeclineInvite, nil)
}

// AcceptDeclineInvite handles a QuicConnFlagDeclineInvite message sent by the peer, by
// sending a nil value to all the channels that are waiting for a QuicConn object
// from the peer. This method ensures that the state of the agent is consistent
// with the peer's decision not to connect to the local node.
func (agent *stdQuicAgent) AcceptDeclineInvite(stream quic.ReceiveStream, conn QuicConn) error {
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
