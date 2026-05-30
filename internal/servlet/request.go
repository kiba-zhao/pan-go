// Define peer application request
package servlet

import (
	"bytes"
	"encoding/binary"
	"io"
	"slices"
)

type RequestName = []byte

type Request struct {
	Message
	name RequestName
}

// NewRequest creates a new request with the given name and body.
// It returns the created request.
func NewRequest(name RequestName, body io.Reader) *Request {
	request := &Request{}
	request.name = name
	request.Reader = body
	InitRequest(request)
	return request
}

// Name returns the name associated with the request.
func (r *Request) Name() RequestName {
	return r.name
}

// MarshalRequest serializes the `Request` into an `io.Reader`.
// It returns a reader containing the name of the request and the binary data of the request message.
// The name is serialized by first appending the length of the name as a uint32, followed by the name bytes.
// Then the request message is serialized using `MarshalMessage`.
func MarshalRequest(request *Request) io.Reader {
	name := request.Name()
	nameBuffer := make([]byte, 0)
	nameBuffer = binary.BigEndian.AppendUint32(nameBuffer, uint32(len(name)))
	nameBuffer = slices.Concat(nameBuffer, name)

	msgReader := MarshalMessage(&request.Message)
	return io.MultiReader(bytes.NewReader(nameBuffer), msgReader)
}

// UnmarshalRequest deserializes the request items from the given io.Reader
// and sets each item on the given Request.
// It reads the request items until an io.EOF is encountered.
// If any error occurs during the deserialization, it returns the error.
func UnmarshalRequest(reader io.Reader, request *Request) error {

	name, err := ParseSegment(reader)
	if err != nil {
		return err
	}
	request.name = name

	return UnmarshalMessage(reader, &request.Message)
}

// InitRequest initializes the given Request by setting its header to an empty slice of HeaderItem pointers.
// This function should be called before using the Request to ensure that the items slice is properly initialized.
func InitRequest(request *Request) {
	request.header = &Header{}
	InitHeader(request.header)
}

// SetRequestScope sets the name of the given request to the given scope.
// The scope and the current request name are combined using the GenerateRouteName function.
func SetRequestScope(request *Request, scope RequestName) {
	request.name = GenerateRouteName(scope, request.Name())
}
