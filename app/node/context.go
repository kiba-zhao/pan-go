package node

import (
	"bytes"
	"io"
	"pan/app/cache"
	"strings"
)

type SessionKey = []byte
type SessionItem struct {
	key   SessionKey
	value interface{}
}

func (si *SessionItem) HashCode() SessionKey {
	return si.key
}

type Session = cache.Bucket[SessionKey, *SessionItem]

type Context struct {
	Response
	request *Request
	session Session
}

func (c *Context) Request() *Request {
	return c.request
}

func (c *Context) Name() RequestName {
	return c.request.Name()
}

func (c *Context) RequestBody() io.Reader {
	return c.request.Body()
}

func (c *Context) RequestHeader(key []byte) ([]byte, bool) {
	return c.request.Header(key)
}

func (c *Context) SetHeader(key, value []byte) {
	if value == nil {
		c.header.Del(key)
		return
	}
	c.header.Set(key, value)
}

func (c *Context) Session(key SessionKey) (interface{}, bool) {
	return c.session.Search(key)
}

func (c *Context) Set(key SessionKey, value interface{}) {
	c.session.Swap(&SessionItem{key, value})
}

func (c *Context) Del(key SessionKey) {
	item, ok := c.session.Search(key)
	if ok {
		c.session.Delete(item)
	}
}

func (c *Context) Respond(body io.Reader) {
	c.body = body
	c.code = 0
}

func (c *Context) ThrowError(code int, err error) {
	c.code = code
	c.body = strings.NewReader(err.Error())
}

func InitContext(ctx *Context) {
	ctx.code = -1
	ctx.session = cache.NewBucket[SessionKey, *SessionItem](bytes.Compare)

	ctx.request = &Request{}
	InitRequest(ctx.request)

	InitResponse(&ctx.Response)
}
