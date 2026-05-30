// Define peer application response
package servlet

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Response struct {
	Message
	code int
	io.Closer
}

// Code returns the status code associated with the response.
func (r *Response) Code() int {
	return r.code
}

// Close closes the underlying io.Closer of the response if it exists.
// It returns an error if the Closer fails to close, otherwise it returns nil.

func (r *Response) Close() error {
	if r.Closer == nil {
		return nil
	}

	return r.Closer.Close()
}

// MarshalResponse serializes the `Response` into an `io.Reader`.
// It returns a reader containing the status code and the binary data of the response message.
// The status code is serialized by appending it as a uint32 to the beginning of the reader.
// Then the response message is serialized using `MarshalMessage`.

func MarshalResponse(response *Response) io.Reader {
	codeBuffer := make([]byte, 0)
	codeBuffer = binary.BigEndian.AppendUint32(codeBuffer, uint32(response.code))

	msgReader := MarshalMessage(&response.Message)
	return io.MultiReader(bytes.NewReader(codeBuffer), msgReader)
}

// UnmarshalResponse deserializes the response data from the given io.ReadCloser
// into the provided Response object. It reads the status code from the reader
// and assigns it to the response. The reader is also set as the Closer for the
// response. It then unmarshals the message part of the response using the
// UnmarshalMessage function. Returns an error if any part of the deserialization
// fails.

func UnmarshalResponse(reader io.ReadCloser, response *Response) error {

	code := uint32(0)
	err := binary.Read(reader, binary.BigEndian, &code)
	if err != nil {
		return err
	}
	response.code = int(code)
	response.Closer = reader

	return UnmarshalMessage(reader, &response.Message)
}

// InitResponse initializes the given Response by setting its header to an empty slice of HeaderItem pointers.
// This function should be called before using the Response to ensure that the items slice is properly initialized.
func InitResponse(response *Response) {
	response.header = &Header{}
	InitHeader(response.header)
}
