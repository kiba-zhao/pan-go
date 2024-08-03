package node

import (
	"bytes"
	"io"

	"slices"

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

type Session = []*SessionItem

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
	idx, ok := slices.BinarySearchFunc(c.session, key, compareSessionItem)
	if !ok {
		return nil, ok
	}
	item := c.session[idx]
	return item.value, true
}

func (c *Context) Set(key SessionKey, value interface{}) {
	idx, ok := slices.BinarySearchFunc(c.session, key, compareSessionItem)
	if ok {
		c.session[idx].value = value
		return
	}

	c.session = slices.Insert(c.session, idx, &SessionItem{key, value})
}

func (c *Context) Del(key SessionKey) {
	idx, ok := slices.BinarySearchFunc(c.session, key, compareSessionItem)
	if ok {
		c.session = slices.Delete(c.session, idx, idx+1)
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
	ctx.session = make([]*SessionItem, 0)

	ctx.request = &Request{}
	InitRequest(ctx.request)

	InitResponse(&ctx.Response)
}

func compareSessionItem(item *SessionItem, key SessionKey) int {
	return bytes.Compare(item.key, key)
}
