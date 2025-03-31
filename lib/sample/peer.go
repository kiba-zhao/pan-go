package sample

import (
	"bytes"
	"context"
	"io"
	"pan/lib/peer"

	"google.golang.org/protobuf/proto"
)

// SamplePeer is sample interface for peer
type SamplePeer interface {
	// Do sends a request to a peer with the provided context, peer ID, and request.
	// It returns the response from the peer and any error encountered during the process.
	Do(context.Context, peer.PeerID, *peer.Request) (*peer.Response, error)
	// Request sends a request to a peer with the provided context, peer ID, request name, and request body.
	// The request body is a proto.Message and is marshaled into a *peer.Request before being sent.
	// The request headers may be set by providing a list of peer.HeaderItem.
	// It returns the response from the peer and any error encountered during the process.
	//
	Request(context.Context, peer.PeerID, peer.RequestName, proto.Message, ...peer.HeaderItem) (*peer.Response, error)
	// RequestWithProto sends a request to a peer with the provided context, peer ID, request name,
	// and request body, and expects a response. The request body and expected response are
	// proto.Messages. The request body is marshaled into a *peer.Request before being sent.
	// The request headers may be set by providing a list of peer.HeaderItem. It returns an error
	// if the request fails or the response cannot be unmarshaled into the provided proto.Message.
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
