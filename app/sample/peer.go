package sample

import (
	"bytes"
	"context"
	"io"
	"pan/app/peer"

	"google.golang.org/protobuf/proto"
)

type SamplePeer interface {
	Do(context.Context, peer.PeerID, *peer.Request) (*peer.Response, error)
	Request(context.Context, peer.PeerID, peer.RequestName, proto.Message, ...peer.HeaderItem) (*peer.Response, error)
	RequestWithProto(context.Context, peer.PeerID, peer.RequestName, proto.Message, proto.Message, ...peer.HeaderItem) error
}

func (s *sample[T]) Do(ctx context.Context, peerId peer.PeerID, request *peer.Request) (*peer.Response, error) {
	scope := s.PeerScope()
	peer.SetRequestScope(request, scope)
	return s.PeerModule.Do(ctx, peerId, request)
}

func (s *sample[T]) Request(ctx context.Context, peerId peer.PeerID, name peer.RequestName, body proto.Message, headerItems ...peer.HeaderItem) (*peer.Response, error) {
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

	if len(headerItems) > 0 {
		for _, headerItem := range headerItems {
			request.SetHeader(headerItem.Key, headerItem.Value)
		}
	}

	return s.Do(ctx, peerId, request)
}

func (s *sample[T]) RequestWithProto(ctx context.Context, peerId peer.PeerID, name peer.RequestName, resp proto.Message, body proto.Message, headerItems ...peer.HeaderItem) error {
	res, err := s.Request(ctx, peerId, name, body, headerItems...)
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
