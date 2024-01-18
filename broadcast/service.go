package broadcast

import (
	"bytes"
	"net"
	"strconv"

	"pan/memory"
	"pan/peer"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type aliveItem struct {
	*memory.BucketItem[[]byte]
	seq     int64
	peerId  peer.PeerId
	expried bool
}

type Service struct {
	serveInfos []*ServeInfo
	seq        int64
	baseId     uuid.UUID
	rw         *sync.RWMutex
	pr         peer.Peer
	store      *memory.Bucket[[]byte, *aliveItem]
}

// TokenRefresh ...
func (s *Service) RefreshToken() {
	s.rw.Lock()
	s.seq = time.Now().Unix()
	s.rw.Unlock()
}

// GenerateAliveMessage ...
func (s *Service) GenerateAliveMessage() (payload []byte, err error) {

	msg := new(Alive)
	s.rw.RLock()
	msg.Seq = s.seq
	msg.BaseId = s.baseId[:]
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

	item := s.store.GetItem(msg.BaseId)
	if item != nil && msg.Seq <= item.seq {
		return
	}

	ip, _, err := net.SplitHostPort(string(addr))
	if err != nil {
		return
	}

	item = new(aliveItem)
	item.BucketItem = memory.NewBucketItem[[]byte](msg.BaseId)

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
			peerId, authErr := s.pr.TwowayAuthenticate(node)
			if authErr != nil {
				err = authErr
				node.Close()
				return
			}
			item.peerId = peerId
		}
	}

	// Implement: set node online
	item.seq = msg.Seq
	item.expried = false
	s.store.SetItem(item)

	return
}

func (s *Service) GenerateDeadMessage() (payload []byte, err error) {

	msg := new(Death)
	s.rw.RLock()
	msg.Seq = s.seq
	msg.BaseId = s.baseId[:]
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

	item := s.store.GetItem(msg.BaseId)
	if item == nil || msg.Seq < item.seq || item.expried {
		return err
	}

	peerId := item.peerId
	state := s.pr.Stat(peerId)
	if state != peer.OfflinePeerState {
		// TODO: clean peer node ?
	}

	item = new(aliveItem)
	item.BucketItem = memory.NewBucketItem[[]byte](msg.BaseId)
	item.seq = msg.Seq
	item.peerId = peerId
	item.expried = true

	s.store.SetItem(item)

	return err
}

// NewService ...
func NewService(baseId uuid.UUID, pr peer.Peer, serveInfos ...*ServeInfo) *Service {
	service := new(Service)
	service.serveInfos = serveInfos
	service.store = memory.NewBucket[[]byte, *aliveItem](bytes.Compare)
	service.baseId = baseId
	service.rw = new(sync.RWMutex)
	service.pr = pr

	service.RefreshToken()
	return service
}
