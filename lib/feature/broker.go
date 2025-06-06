package feature

import (
	"bytes"
	"context"
	"io"
	"pan/lib/app"
	"pan/lib/peer"

	"google.golang.org/protobuf/proto"
)

type BrokerHelper interface {
	Do(ctx context.Context, peerId peer.PeerID, request peer.PeerRequest) (peer.PeerResponse, error)
	Request(ctx context.Context, peerId peer.PeerID, name peer.PeerRequestName, body proto.Message, headerItems ...peer.PeerHeaderItem) (peer.PeerResponse, error)
	RequestWithProto(ctx context.Context, peerId peer.PeerID, name peer.PeerRequestName, resp proto.Message, body proto.Message, headerItems ...peer.PeerHeaderItem) error
}

type stdBrokerHelper struct {
	PeerCluster   peer.PeerCluster
	featureHelper *stdFeatureHelper
}

var _ = (BrokerHelper)((*stdBrokerHelper)(nil))

func (broker *stdBrokerHelper) Do(ctx context.Context, peerId peer.PeerID, request peer.PeerRequest) (peer.PeerResponse, error) {

	featureHelper := broker.featureHelper
	peerCluster := broker.PeerCluster
	if peerCluster == nil || featureHelper == nil {
		return nil, ErrFeatureHelperPeerClusterNotFound
	}

	app.SetRequestScope(request, featureHelper.PeerScope())
	return peerCluster.Do(ctx, peerId, request)
}

func (broker *stdBrokerHelper) Request(ctx context.Context, peerId peer.PeerID, name peer.PeerRequestName, body proto.Message, headerItems ...peer.PeerHeaderItem) (peer.PeerResponse, error) {

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

func (broker *stdBrokerHelper) RequestWithProto(ctx context.Context, peerId peer.PeerID, name peer.PeerRequestName, resp proto.Message, body proto.Message, headerItems ...peer.PeerHeaderItem) error {
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
