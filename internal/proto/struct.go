package proto

import (
	"bytes"
	"io"

	"google.golang.org/protobuf/proto"
)

func MarshalWithReader(message proto.Message) (io.Reader, error) {
	contentBytes, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(contentBytes), nil
}

func UnmarshalWithReader(reader io.Reader, message proto.Message) error {
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return proto.Unmarshal(content, message)
}
