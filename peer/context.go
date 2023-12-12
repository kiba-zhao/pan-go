package peer

import (
	"io"
	"pan/core"
	"strings"
)

type Context interface {
	core.Context
	PeerId() PeerId
	Headers() []*HeaderSegment
	Header(name []byte) []byte
	Body() io.Reader
	Respond(body io.Reader, headers ...*HeaderSegment) error
	ThrowError(code int, message string, headers ...*HeaderSegment) error
}

type UnknownContext interface {
	Context
	Node() Node
}

type contextSt struct {
	*Request
	stream NodeStream
	peerId PeerId
}

// PeerId ...
func (c *contextSt) PeerId() PeerId {
	return c.peerId
}

// Respond ...
func (c *contextSt) Respond(body io.Reader, headers ...*HeaderSegment) (err error) {
	err = c.respond(0, body, headers...)
	return
}

// ThrowError ...
func (c *contextSt) ThrowError(code int, message string, headers ...*HeaderSegment) (err error) {

	body := strings.NewReader(message)
	err = c.respond(code, body, headers...)
	return
}

// respond ...
func (c *contextSt) respond(code int, body io.Reader, headers ...*HeaderSegment) (err error) {
	res := NewReponse(code, body, headers...)
	reader, err := MarshalResponse(res)
	if err != nil {
		return
	}

	_, err = io.Copy(c.stream, reader)
	if err != nil {
		return
	}

	err = c.stream.Close()
	return
}

// NewContext ...
func NewContext(stream NodeStream, peerId PeerId) (ctx Context, err error) {

	req := new(Request)
	err = UnmarshalRequest(stream, req)
	if err == nil {
		ctxSt := new(contextSt)
		ctxSt.Request = req
		ctxSt.stream = stream
		ctxSt.peerId = peerId
		ctx = ctxSt
	}
	return
}
