package peer

const (
	QUICRouteType = uint8(iota)
	TCPRouteType
)

type PeerPassport struct {
	PeerId PeerId
	Enable bool
}

type PeerRepository interface {
	FindOne(peerId PeerId) (*PeerPassport, error)
}

type PeerRoute struct {
	PeerId PeerId
	Type   uint8
	Addr   []byte
	Delay  int
	Rate   int
	Enable bool
}

type PeerRouteRepository interface {
	FindByPeerIdAndEnable(peerId PeerId, enable bool) ([]*PeerRoute, error)
	UpdateOne(route *PeerRoute) error
	CreateOne(route *PeerRoute) error
}
