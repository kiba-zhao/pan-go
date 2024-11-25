package peer

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Message struct {
	header *Header
	io.Reader
}

func (r *Message) Header(key []byte) ([]byte, bool) {
	return r.header.Get(key)
}

func (r *Message) SetHeader(key, value []byte) {
	if value == nil {
		r.header.Del(key)
		return
	}
	r.header.Set(key, value)
}

func MarshalMessage(message *Message) io.Reader {
	headerReader, headerSize := MarshalHeader(message.header)

	headerSizeBuffer := make([]byte, 0)
	headerSizeBuffer = binary.BigEndian.AppendUint32(headerSizeBuffer, uint32(headerSize))
	headerSizeReader := bytes.NewReader(headerSizeBuffer)

	var readers []io.Reader
	readers = append(readers, headerSizeReader)
	if headerSize > 0 {
		readers = append(readers, headerReader)
	}
	if message.Reader != nil {
		readers = append(readers, message.Reader)
	}

	if len(readers) > 1 {
		return io.MultiReader(readers...)
	}
	return readers[0]
}

func UnmarshalMessage(reader io.Reader, message *Message) error {

	headerSizeBuffer := make([]byte, 4)
	headerSizeNum, err := reader.Read(headerSizeBuffer)
	if err != nil && headerSizeNum != 4 {
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
	message.Reader = reader
	return nil
}
