// Define peer application message
package app

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Message struct {
	header *Header
	io.Reader
}

// Header returns the value associated with the given key from the message header.
// It returns the value and a boolean indicating whether the key was found.
func (r *Message) Header(key []byte) ([]byte, bool) {
	return r.header.Get(key)
}

// SetHeader sets the value of a header item with the specified key.
// If the value is nil, the header item is removed. Otherwise, the header item is set or updated.

func (r *Message) SetHeader(key, value []byte) {
	if value == nil {
		r.header.Del(key)
		return
	}
	r.header.Set(key, value)
}

// MarshalMessage serializes a Message into an io.Reader.
// It marshals the message header and its size into a byte stream,
// followed by the message body if it exists. The resulting io.Reader
// can be used to read the complete serialized message.

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

// UnmarshalMessage deserializes a Message from the given io.Reader.
// It reads the header size from the reader, initializes a new Header,
// and unmarshals the header if its size is greater than zero. The function
// then assigns the deserialized header and remaining reader to the Message.
// Returns an error if any part of the deserialization fails.

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
