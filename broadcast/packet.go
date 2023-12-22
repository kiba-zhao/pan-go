package broadcast

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// ParsePacket ...
func ParsePacket(payload []byte) (packet []byte, size int, err error) {
	maxLen := len(payload)
	size = -1
	offset := -1
	if maxLen <= 5 {
		err = errors.New("Payload Not Enough")
		return
	}

	version, offset := parseVersion(payload, 0)
	if version != 1 {
		return
	}

	size, offset, err = parseSize(payload, offset)
	if err == nil {
		packet = payload[offset:size]
	}

	return
}

func parseVersion(payload []byte, offset int) (version uint8, eOffset int) {
	version = uint8(payload[offset])
	eOffset = offset + 1
	return
}

func parseSize(payload []byte, offset int) (size int, eOffset int, err error) {
	eOffset = offset + 4
	maxLen := len(payload)
	if eOffset > maxLen {
		err = errors.New("Payload Not Enough")
		return
	}

	var usize uint32 = 0
	err = binary.Read(bytes.NewReader(payload[offset:eOffset]), binary.BigEndian, &usize)
	if err != nil {
		return
	}

	size = int(usize)
	if size > maxLen {
		err = errors.New("Payload Not Enough")
	}
	return
}

// MarshalPacket ...
func MarshalPacket(parts ...[]byte) (payload []byte) {

	_parts := make([][]byte, 0)
	head := []byte{1, 0, 0, 0, 0}
	_parts = append(_parts, head)
	_parts = append(_parts, parts...)

	size := len(head)
	for _, part := range parts {
		size += len(part)
	}

	binary.BigEndian.PutUint32(head[1:], uint32(size))
	payload = bytes.Join(_parts, nil)

	return
}
