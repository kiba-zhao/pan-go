package peer

import (
	"bytes"
	"errors"
	"io"
)

type Request struct {
	method  []byte
	headers []*HeaderSegment
	body    io.Reader
}

// Method ...
func (r *Request) Method() []byte {
	return r.method
}

// Headers ...
func (r *Request) Headers() []*HeaderSegment {
	return r.headers
}

// Header ...
func (r *Request) Header(name []byte) (value []byte) {
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
func (r *Request) Body() io.Reader {
	return r.body
}

// NewRequest ...
func NewRequest(method []byte, body io.Reader, headers ...*HeaderSegment) *Request {
	req := new(Request)
	req.method = method
	req.body = body
	req.headers = headers
	return req
}

// UnmarshalRequest ...
func UnmarshalRequest(reader io.Reader, request *Request) (err error) {
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
		if bytes.Equal([]byte("Method"), header.Name()) {
			request.method = header.Value()
			continue
		}

		headers = append(headers, header)

	}

	if len(headers) > 0 {
		request.headers = headers
	}

	if ntype == BodySegmentType {
		request.body = reader
	}

	return
}

// MarshalRequest
func MarshalRequest(request *Request) (reader io.Reader, err error) {

	readers := make([]io.Reader, 0)
	method := request.Method()
	if method == nil {
		err = errors.New("Method Should not be nil")
		return
	}
	mstr := CreateSegmentType(HeaderSegmentType)
	msrs := CreateHeaderSegment(&HeaderSegment{name: []byte("Method"), value: method})
	readers = append(readers, mstr)
	readers = append(readers, msrs...)

	headers := request.Headers()
	if headers != nil {
		for _, header := range headers {
			if header == nil {
				continue
			}
			hstr := CreateSegmentType(HeaderSegmentType)
			hsrs := CreateHeaderSegment(header)
			readers = append(readers, hstr)
			readers = append(readers, hsrs...)
		}
	}

	body := request.Body()
	if body != nil {
		bstr := CreateSegmentType(BodySegmentType)
		readers = append(readers, bstr, body)
	}

	reader = io.MultiReader(readers...)
	return
}
