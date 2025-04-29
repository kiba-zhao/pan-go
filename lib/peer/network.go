// Define peer network
package peer

import (
	"context"
	"errors"
	"io"
	"pan/lib/log"
)

var ErrPeerNetworkUnavailable = errors.New("peer.PeerNetwork Error: Unavailable")

type PeerNetwork interface {
	RoundTrip(context.Context, PeerID, io.Reader) (io.ReadCloser, error)
	CanReach(PeerID) bool
}

type PeerNetworkPurgeable interface {
	Purge(PeerID) error
}

type peerNetwork struct {
	peerModule *peerModule
}

// CanReach checks if the given peerId is reachable by any of the peer networks
// associated with the peer module. It returns true if at least one peer network
// can reach the specified peerId, otherwise it returns false.

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

// Purge is a no-op if there are no peer networks associated with the peer module.
//
// It iterates over the peer networks and calls their Purge method with the given
// peerId. If any error occurs during the iteration, it is logged at the error level.
//
// The function does not return an error.
func (pn *peerNetwork) Purge(peerId PeerID) error {
	modules := pn.peerModule.PeerNetworks()
	if len(modules) <= 0 {
		return nil
	}

	for _, module := range modules {
		peerNetwork, ok := module.(PeerNetworkPurgeable)
		if !ok {
			continue
		}
		err := peerNetwork.Purge(peerId)
		if err != nil {
			log.Default().Log(context.Background(), log.LevelError, err.Error())
		}
	}
	return nil
}

// RoundTrip sends a request to a peer and returns the response.
//
// It iterates over the peer networks and calls their RoundTrip method with the given
// peerId and reader. If any error occurs during the iteration, it is returned.
//
// If at least one peer network can reach the specified peerId, the response is
// returned. Otherwise, ErrPeerNetworkUnavailable is returned.
//
// If the given context is canceled or the request reader has been read before, the
// iteration is stopped and the last error is returned.
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

		// break with success or request reader has been read
		if resReader != nil || reqReader.haveRead {
			break
		}

		// break with context is done
		if errors.Is(err, ctx.Err()) {
			break
		}

		// break access denied error
		if errors.Is(err, ErrPeerModuleAccessDenied) {
			break
		}

	}

	return resReader, err
}
