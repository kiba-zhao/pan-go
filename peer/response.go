package peer

import (
	"bytes"
	"encoding/binary"

	"io"
)

type Response struct {
	code    int
	headers []*HeaderSegment
	body    io.Reader
}

// IsError ...
func (r *Response) IsError() bool {
	return r.code != 0
}

// Code ...
func (r *Response) Code() int {
	return r.code
}

// Headers ...
func (r *Response) Headers() []*HeaderSegment {
	return r.headers
}

// Header ...
func (r *Response) Header(name []byte) (value []byte) {
	if r.headers != nil {
		for _, header := range r.headers {
			if bytes.Equal(header.Name(), name) {
				value = header.Value()
				break
			}
		}
	}
	return
}

// Body ...
func (r *Response) Body() io.Reader {
	return r.body
}

// NewReponse ...
func NewReponse(code int, body io.Reader, headers ...*HeaderSegment) *Response {
	response := new(Response)
	response.code = code
	response.body = body
	response.headers = headers
	return response
}

// UnmarshalResponse ...
func UnmarshalResponse(reader io.Reader, response *Response) (err error) {
	headers := make([]*HeaderSegment, 0)
	var ntype uint8
	for {
		headerType, err := ParseSegmentType(reader)
		if err != nil || headerType != HeaderSegmentType {
			ntype = headerType
			break
		}
		header, err := ParseHeaderSegment(reader)
		if err != nil {
			break
		}

		if bytes.Equal([]byte("PeerCode"), header.Name()) {
			code := binary.BigEndian.Uint32(header.Value())
			response.code = int(code)
			continue
		}

		headers = append(headers, header)

	}

	if len(headers) > 0 {
		response.headers = headers
	}

	if ntype == BodySegmentType {
		response.body = reader
	}

	return
}

// MarshalResponse
func MarshalResponse(response *Response) (reader io.Reader, err error) {

	readers := make([]io.Reader, 0)
	code := make([]byte, 4)
	binary.BigEndian.PutUint32(code, uint32(response.Code()))
	cstr := CreateSegmentType(HeaderSegmentType)
	csrs := CreateHeaderSegment(&HeaderSegment{name: []byte("PeerCode"), value: code})
	readers = append(readers, cstr)
	readers = append(readers, csrs...)

	headers := response.Headers()
	if headers != nil {
		for _, header := range headers {
			hstr := CreateSegmentType(HeaderSegmentType)
			hsrs := CreateHeaderSegment(header)
			readers = append(readers, hstr)
			readers = append(readers, hsrs...)
		}
	}

	body := response.Body()
	if body != nil {
		bstr := CreateSegmentType(BodySegmentType)
		readers = append(readers, bstr, body)
	}

	reader = io.MultiReader(readers...)
	return
}
