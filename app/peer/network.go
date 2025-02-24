package peer

import (
	"context"
	"errors"
	"io"
	"pan/logger"
)

var ErrPeerNetworkUnavailable = errors.New("peer.PeerNetwork Error: Unavailable")

type PeerNetwork interface {
	RoundTrip(context.Context, PeerID, io.Reader) (io.ReadCloser, error)
	CanReach(PeerID) bool
	Purge(PeerID) error
}

type peerNetwork struct {
	peerModule *peerModule
}

func (pn *peerNetwork) CanReach(peerId PeerID) bool {

	modules := pn.peerModule.PeerNetworks()
	if len(modules) <= 0 {
		return false
	}

	for _, module := range modules {
		peerNetwork := module.(PeerNetwork)
		if peerNetwork.CanReach(peerId) {
			return true
		}
	}
	return false
}

func (pn *peerNetwork) Purge(peerId PeerID) error {
	modules := pn.peerModule.PeerNetworks()
	if len(modules) <= 0 {
		return nil
	}

	for _, module := range modules {
		peerNetwork := module.(PeerNetwork)
		err := peerNetwork.Purge(peerId)
		if err != nil {
			logger.Default().Log(context.Background(), logger.LevelError, err.Error())
		}
	}
	return nil
}

func (pn *peerNetwork) RoundTrip(ctx context.Context, peerId PeerID, reader io.Reader) (io.ReadCloser, error) {

	modules := pn.peerModule.PeerNetworks()
	if len(modules) <= 0 {
		return nil, ErrPeerNetworkUnavailable
	}

	var resReader io.ReadCloser
	var err error
	reqReader := &peerReader{Reader: reader}
	for _, module := range modules {
		peerNetwork := module.(PeerNetwork)
		if !peerNetwork.CanReach(peerId) {
			continue
		}
		resReader, err = peerNetwork.RoundTrip(ctx, peerId, reqReader)
		if err == nil || ctx.Err() == err || reqReader.haveRead {
			break
		}
	}

	return resReader, err
}
