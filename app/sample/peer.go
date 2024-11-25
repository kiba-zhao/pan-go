package sample

import (
	"bytes"
	"io"
	"pan/app/peer"

	"google.golang.org/protobuf/proto"
)

type SamplePeer interface {
	Do(peer.PeerID, *peer.Request, ...peer.PeerDoContextUpdater) (*peer.Response, error)
	Request(peer.PeerID, peer.RequestName, proto.Message, ...peer.PeerDoContextUpdater) (*peer.Response, error)
	RequestWithProto(peer.PeerID, peer.RequestName, proto.Message, proto.Message, ...peer.PeerDoContextUpdater) error
}

func (s *sample[T]) Do(peerId peer.PeerID, request *peer.Request, updaters ...peer.PeerDoContextUpdater) (*peer.Response, error) {
	scope := s.PeerScope()
	peer.SetRequestScope(request, scope)
	return s.PeerModule.Do(peerId, request, updaters...)
}

func (s *sample[T]) Request(peerId peer.PeerID, name peer.RequestName, body proto.Message, updaters ...peer.PeerDoContextUpdater) (*peer.Response, error) {
	var request *peer.Request
	if body != nil {
		requestBytes, err := proto.Marshal(body)
		if err != nil {
			return nil, err
		}
		request = peer.NewRequest(name, bytes.NewReader(requestBytes))
	} else {
		request = peer.NewRequest(name, nil)
	}

	return s.Do(peerId, request, updaters...)
}

func (s *sample[T]) RequestWithProto(peerId peer.PeerID, name peer.RequestName, resp proto.Message, body proto.Message, updaters ...peer.PeerDoContextUpdater) error {
	res, err := s.Request(peerId, name, body, updaters...)
	if err != nil {
		return err
	}
	defer res.Close()

	var content []byte
	content, err = io.ReadAll(res)
	if err != nil {
		return err
	}
	return proto.Unmarshal(content, resp)
}
