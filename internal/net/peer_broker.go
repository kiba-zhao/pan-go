package net

import (
	"context"
	"io"
	"pan/internal/servlet"
)

type PeerBroker struct {
	PeerNetwork     PeerNetwork
	PeerRouteModule PeerRouteModule
}

func (broker *PeerBroker) DoAction(ctx context.Context, peerId PeerID, reader io.Reader) (ActionSession, io.ReadCloser, error) {
	stream, err := broker.PeerNetwork.RoundTrip(ctx, peerId)
	if err != nil {
		return nil, nil, err
	}
	return DoAction(ctx, stream, reader)
}

func (broker *PeerBroker) NewRequest(name servlet.RequestName, bodyReader io.Reader, headerItems ...HeaderItem) io.Reader {
	scope := broker.PeerRouteModule.PeerRouteScope()
	realName := servlet.GenerateRouteName(scope, name)
	reqReader := NewRequest(realName, bodyReader, headerItems...)
	return reqReader
}
