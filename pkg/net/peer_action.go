package net

import (
	"context"
	"io"
	"pan/pkg/servlet"
)

type HeaderItem = servlet.HeaderItem

type ActionSession interface {
	Header(key []byte) ([]byte, bool)
}

func DoAction(ctx context.Context, stream PeerStream, reader io.Reader) (ActionSession, io.ReadCloser, error) {

	ctx_, cancel := context.WithCancelCause(ctx)
	go func(cancel context.CancelCauseFunc) {
		_, err := io.Copy(stream, reader)
		cancel(err)
	}(cancel)

	<-ctx_.Done()
	err := ctx_.Err()
	if err != nil {
		stream.Close()
		return nil, nil, err
	}

	res := &servlet.Response{}
	servlet.InitResponse(res)
	err = servlet.UnmarshalResponse(stream, res)
	if err == nil && res.Code() != CodeOK {
		var content []byte
		content, err = io.ReadAll(res)
		if err == nil {
			err = &PeerError{code: res.Code(), err: string(content)}
		}
	}
	if err != nil {
		stream.Close()
		return nil, nil, err
	}

	return res, res, err
}

func NewRequest(name servlet.RequestName, bodyReader io.Reader, headerItems ...HeaderItem) io.Reader {
	req := servlet.NewRequest(name, bodyReader)
	if len(headerItems) > 0 {
		for _, item := range headerItems {
			req.SetHeader(item.Key, item.Value)
		}
	}

	return servlet.MarshalRequest(req)
}
