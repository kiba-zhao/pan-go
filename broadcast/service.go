package broadcast

import (
	"bytes"
	"crypto/rand"
	"net"
	"strconv"

	"pan/peer"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

type Service struct {
	serveInfos []*ServeInfo
	repo       Repo
	seq        int64
	token      []byte
	rw         *sync.RWMutex
	pr         peer.Peer
}

// TokenRefresh ...
func (s *Service) RefreshToken() {
	s.rw.Lock()
	s.seq = time.Now().Unix()
	rand.Read(s.token)
	s.rw.Unlock()
}

// GenerateAliveMessage ...
func (s *Service) GenerateAliveMessage() (payload []byte, err error) {

	msg := new(Alive)
	s.rw.RLock()
	msg.Seq = s.seq
	msg.Token = s.token
	s.rw.RUnlock()

	msg.ServeInfos = s.serveInfos
	payload, err = proto.Marshal(msg)

	return
}

// RecvAliveMessage ...
func (s *Service) RecvAliveMessage(addr []byte, payload []byte) (err error) {
	msg := new(Alive)
	err = proto.Unmarshal(payload, msg)
	if err != nil {
		return
	}
	// Implement: Intercept invalid messages

	rd, err := s.repo.FindOneWithAddrAndSeq(addr, msg.Seq)
	if err != nil || (rd != nil && rd.Seq >= msg.Seq) {
		return
	}
	ip, _, err := net.SplitHostPort(string(addr))
	if err != nil {
		return
	}
	if rd == nil {
		rd = new(Record)
	}

	// Implement: verify alive message
	for _, serveInfo := range msg.ServeInfos {
		switch serveInfo.Type[0] {
		case peer.QUICNodeType:
			addrString := net.JoinHostPort(ip, strconv.Itoa(int(serveInfo.Port)))
			udpAddr, addrErr := net.ResolveUDPAddr("udp", addrString)
			if addrErr != nil {
				err = addrErr
				return
			}
			quicAddr := peer.MarshalQUICAddr(udpAddr)
			node, connErr := s.pr.Connect(peer.QUICNodeType, quicAddr)
			if connErr != nil {
				err = addrErr
				return
			}
			peerId, authErr := s.pr.Authenticate(node, peer.NormalAuthenticateMode)
			if authErr != nil {
				err = authErr
				node.Close()
				return
			}
			rd.PeerId = peerId[:]
		}
	}

	// Implement: set node online
	rd.Seq = msg.Seq
	rd.Token = msg.Token
	rd.Addr = addr
	err = s.repo.Save(rd)

	return err
}

func (s *Service) GenerateDeadMessage() (payload []byte, err error) {

	msg := new(Death)
	s.rw.RLock()
	msg.Seq = s.seq
	msg.Token = s.token
	s.rw.RUnlock()

	payload, err = proto.Marshal(msg)
	return
}

// RecvDeadMessage ...
func (s *Service) RecvDeadMessage(addr []byte, payload []byte) error {
	msg := new(Death)
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		return err
	}
	// Implement: Intercept invalid messages
	rd, err := s.repo.FindOneWithAddrAndSeq(addr, msg.Seq)
	if err != nil || rd == nil || rd.DeathTime > 0 || rd.Seq != msg.Seq || !bytes.Equal(rd.Token, msg.Token) {
		return err
	}

	// Implement: verify death message
	if rd.PeerId != nil {
		peerId := peer.PeerId(rd.PeerId)
		state := s.pr.Stat(peerId)
		if state != peer.OfflinePeerState {
			// TODO: clean peer node ?
		}
	}

	// Implement: set node offline
	rd.Seq = msg.Seq
	rd.Token = msg.Token
	rd.Addr = addr
	rd.DeathTime = time.Now().Unix()
	err = s.repo.Save(rd)

	return err
}

// NewService ...
func NewService(repo Repo, pr peer.Peer, serveInfos ...*ServeInfo) *Service {
	service := new(Service)
	service.serveInfos = serveInfos
	service.repo = repo
	service.token = make([]byte, 32)
	service.rw = new(sync.RWMutex)
	service.pr = pr

	service.RefreshToken()
	return service
}
