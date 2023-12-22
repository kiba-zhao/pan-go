package broadcast

import (
	"pan/core"
)

type Context interface {
	core.Context
	Body() []byte
	Addr() []byte
	Net() Net
}

type contextStruct struct {
	method []byte
	body   []byte
	addr   []byte
	n      Net
}

// Method ...
func (c *contextStruct) Method() []byte {
	return c.method
}

// Body ...
func (c *contextStruct) Body() []byte {
	return c.body
}

// Addr ...
func (c *contextStruct) Addr() []byte {
	return c.addr
}

// Net ...
func (c *contextStruct) Net() Net {
	return c.n
}

// NewContext ...
func NewContext(method, body []byte, addr []byte, n Net) Context {

	ctx := new(contextStruct)
	ctx.method = method
	ctx.body = body
	ctx.addr = addr
	ctx.n = n

	return ctx
}
