package node

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

func (r *Request) Name() RequestName {
	return r.name
}

func MarshalRequest(request *Request) io.Reader {
	name := request.Name()
	nameBuffer := make([]byte, 0)
	nameBuffer = binary.BigEndian.AppendUint32(nameBuffer, uint32(len(name)))
	nameBuffer = slices.Concat(nameBuffer, name)

	msgReader := MarshalMessage(&request.Message)
	return io.MultiReader(bytes.NewReader(nameBuffer), msgReader)
}

func UnmarshalRequest(reader io.Reader, request *Request) error {

	name, err := ParseSegment(reader)
	if err != nil {
		return err
	}
	request.name = name

	return UnmarshalMessage(reader, &request.Message)
}

func InitRequest(request *Request) {
	request.header = &Header{}
	InitHeader(request.header)
}
