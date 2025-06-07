package feature

import (
	"bytes"
	"context"
	"errors"
	"io"
	"pan/lib/app"
	"pan/lib/peer"
	"sync"

	"google.golang.org/protobuf/proto"
)

var ErrFeatureBrokerUnAvailable = errors.New("feature.Broker Error: UnAvailable")

type BrokerHelper interface {
	Do(ctx context.Context, peerId peer.PeerID, request peer.PeerRequest) (peer.PeerResponse, error)
}

var _ = (BrokerHelper)((*stdFeatureHelper)(nil))

func (helper *stdFeatureHelper) Do(ctx context.Context, peerId peer.PeerID, request peer.PeerRequest) (peer.PeerResponse, error) {
	cluster := helper.PeerCluster
	if cluster == nil {
		return nil, ErrFeatureHelperPeerClusterNotFound
	}

	app.SetRequestScope(request, helper.PeerScope())
	return cluster.Do(ctx, peerId, request)
}

type Broker interface {
	InitBroker(helper BrokerHelper)
}

type BaseBroker struct {
	helper BrokerHelper
	rw     sync.RWMutex
}

var _ = (Broker)((*BaseBroker)(nil))

func (broker *BaseBroker) InitBroker(helper BrokerHelper) {
	broker.rw.Lock()
	defer broker.rw.Unlock()

	broker.helper = helper
}

func (broker *BaseBroker) Do(ctx context.Context, peerId peer.PeerID, request peer.PeerRequest) (peer.PeerResponse, error) {

	broker.rw.RLock()
	defer broker.rw.RUnlock()

	helper := broker.helper
	if helper == nil {
		return nil, ErrFeatureBrokerUnAvailable
	}

	return helper.Do(ctx, peerId, request)
}

func (broker *BaseBroker) Request(ctx context.Context, peerId peer.PeerID, name peer.PeerRequestName, body proto.Message, headerItems ...peer.PeerHeaderItem) (peer.PeerResponse, error) {

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

func (broker *BaseBroker) RequestWithProto(ctx context.Context, peerId peer.PeerID, name peer.PeerRequestName, resp proto.Message, body proto.Message, headerItems ...peer.PeerHeaderItem) error {
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

type BrokerMeta = Metadata[Broker]

func NewBrokerMeta[T any](broker Broker) BrokerMeta {
	return newMetadata[T, Broker](broker)
}

type BrokerMetaProvider interface {
	BrokerMetaList() []BrokerMeta
}

func getBrokerMetaList(feature interface{}) []BrokerMeta {
	if metaProvider, ok := feature.(BrokerMetaProvider); ok {
		return metaProvider.BrokerMetaList()
	}
	return nil
}
