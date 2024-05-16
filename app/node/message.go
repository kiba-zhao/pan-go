package node

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Message struct {
	header *Header
	body   io.Reader
}

func (r *Message) Header(key []byte) ([]byte, bool) {
	return r.header.Get(key)
}

func (r *Message) Body() io.Reader {
	return r.body
}

func MarshalMessage(message *Message) io.Reader {
	headerReader, headerSize := MarshalHeader(message.header)

	headerSizeBuffer := make([]byte, 0)
	headerSizeBuffer = binary.BigEndian.AppendUint32(headerSizeBuffer, uint32(headerSize))
	headerSizeReader := bytes.NewReader(headerSizeBuffer)
	return io.MultiReader(headerSizeReader, headerReader, message.body)
}

func UnmarshalMessage(reader io.Reader, message *Message) error {

	headerSizeBuffer := make([]byte, 4)
	_, err := reader.Read(headerSizeBuffer)
	if err != nil {
		return err
	}

	header := &Header{}
	InitHeader(header)

	headerSize := binary.BigEndian.Uint32(headerSizeBuffer)
	if headerSize > 0 {
		headerReader := io.LimitReader(reader, int64(headerSize))
		err = UnmarshalHeader(headerReader, header)
		if err != nil {
			return err
		}
	}

	message.header = header
	message.body = reader
	return nil
}
