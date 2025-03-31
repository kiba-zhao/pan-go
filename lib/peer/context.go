// Define peer application execution context

package peer

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

// HashCode returns the key of the SessionItem as the hash code.
func (si *SessionItem) HashCode() SessionKey {
	return si.key
}

type Session = []*SessionItem

type Context struct {
	Response
	request *Request
	session Session
}

// Request returns the request associated with the context.
func (c *Context) Request() *Request {
	return c.request
}

// Name returns the request name associated with the context.
func (c *Context) Name() RequestName {
	return c.request.Name()
}

// RequestBody returns the request body associated with the context as an io.Reader.

func (c *Context) RequestBody() io.Reader {
	return c.request
}

// RequestHeader returns the request header associated with the context
// by the given key. The second return value is true if the header exists,
// false otherwise.
func (c *Context) RequestHeader(key []byte) ([]byte, bool) {
	return c.request.Header(key)
}

// Session retrieves the value associated with the given SessionKey from the context's session.
// It returns the value and a boolean indicating whether the key was found.

func (c *Context) Session(key SessionKey) (interface{}, bool) {
	idx, ok := slices.BinarySearchFunc(c.session, key, compareSessionItem)
	if !ok {
		return nil, ok
	}
	item := c.session[idx]
	return item.value, true
}

// Set associates the given value with the specified SessionKey in the context's session.
// If the key already exists, its value is updated. If the key does not exist, a new
// SessionItem is inserted in the session.

func (c *Context) Set(key SessionKey, value interface{}) {
	idx, ok := slices.BinarySearchFunc(c.session, key, compareSessionItem)
	if ok {
		c.session[idx].value = value
		return
	}

	c.session = slices.Insert(c.session, idx, &SessionItem{key, value})
}

// Del removes the SessionItem associated with the given SessionKey from the context's session.
// If the key is not found, Del does nothing.
func (c *Context) Del(key SessionKey) {
	idx, ok := slices.BinarySearchFunc(c.session, key, compareSessionItem)
	if ok {
		c.session = slices.Delete(c.session, idx, idx+1)
	}
}

// Respond sets the response body associated with the context.
//
// If the body also implements io.Closer, then the context will close the body when the context is closed.
func (c *Context) Respond(body io.Reader) {
	c.Reader = body
	c.code = CodeOK
	if closer, ok := body.(io.Closer); ok {
		c.Closer = closer
	}
}

// ThrowError sets the response code and body associated with the context based on the given error.
// The error is converted to a string and set as the response body.
func (c *Context) ThrowError(code int, err error) {
	c.code = code
	c.Reader = strings.NewReader(err.Error())
}

// InitContext initializes a Context with default values.
// It sets the response code to -1 (unknown), initializes the session as an empty slice,
// and sets the request and response bodies to empty values.
func InitContext(ctx *Context) {
	ctx.code = -1
	ctx.session = make([]*SessionItem, 0)

	ctx.request = &Request{}
	InitRequest(ctx.request)

	InitResponse(&ctx.Response)
}

// compareSessionItem compares the key of a SessionItem with a given SessionKey.
// Returns an integer indicating the result of the comparison:
// a negative number if item.key < key, zero if item.key == key, and a positive number if item.key > key.

func compareSessionItem(item *SessionItem, key SessionKey) int {
	return bytes.Compare(item.key, key)
}
