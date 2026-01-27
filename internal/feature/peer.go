package feature

import (
	"bytes"
	"context"
	"errors"
	"io"
	"pan/internal/app"
	"pan/internal/peer"

	"google.golang.org/protobuf/proto"
)

var ErrPeerBrokerPeerNetworkNotFound = errors.New("feature.PeerBroker Error: Peer Network Not Found")

type PeerAppModule interface {
	PeerScope() []byte
}

type PeerBroker struct {
	PeerNetwork   peer.PeerNetwork
	PeerAppModule PeerAppModule
}

func (broker *PeerBroker) Do(ctx context.Context, peerId peer.PeerID, request peer.PeerRequest) (peer.PeerResponse, error) {
	network := broker.PeerNetwork
	if network == nil {
		return nil, ErrPeerBrokerPeerNetworkNotFound
	}

	if broker.PeerAppModule != nil {
		scope := broker.PeerAppModule.PeerScope()
		if len(scope) > 0 {
			app.SetRequestScope(request, scope)
		}
	}
	return network.Do(ctx, peerId, request)
}

func (broker *PeerBroker) Request(ctx context.Context, peerId peer.PeerID, name peer.PeerRequestName, body proto.Message, headerItems ...peer.PeerHeaderItem) (peer.PeerResponse, error) {

	var request peer.PeerRequest
	if body != nil {
		requestBytes, err := proto.Marshal(body)
		if err != nil {
			return nil, err
		}
		request = app.NewRequest(name, bytes.NewReader(requestBytes))
	} else {
		request = app.NewRequest(name, nil)
	}

	if len(headerItems) > 0 {
		for _, headerItem := range headerItems {
			request.SetHeader(headerItem.Key, headerItem.Value)
		}
	}

	return broker.Do(ctx, peerId, request)
}

func (broker *PeerBroker) RequestWithProto(ctx context.Context, peerId peer.PeerID, name peer.PeerRequestName, resp proto.Message, body proto.Message, headerItems ...peer.PeerHeaderItem) error {
	res, err := broker.Request(ctx, peerId, name, body, headerItems...)
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

type PeerTopic = peer.PeerAppModule

type PeerTopicProvider interface {
	PeerTopics() []PeerTopic
}

func getPeerTopics(feature interface{}) []PeerTopic {
	if provider, ok := feature.(PeerTopicProvider); ok {
		return provider.PeerTopics()
	}
	return nil
}
