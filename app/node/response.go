package node

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Response struct {
	Message
	code int
}

func (r *Response) Code() int {
	return r.code
}

func MarshalResponse(response *Response) io.Reader {
	codeBuffer := make([]byte, 0)
	codeBuffer = binary.BigEndian.AppendUint32(codeBuffer, uint32(response.code))

	msgReader := MarshalMessage(&response.Message)
	return io.MultiReader(bytes.NewReader(codeBuffer), msgReader)
}

func UnmarshalResponse(reader io.Reader, response *Response) error {

	code := uint32(0)
	err := binary.Read(reader, binary.BigEndian, &code)
	if err != nil {
		return err
	}
	response.code = int(code)

	return UnmarshalMessage(reader, &response.Message)
}

func InitResponse(response *Response) {
	response.header = &Header{}
	InitHeader(response.header)
}
