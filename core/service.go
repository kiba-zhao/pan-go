package core

import (
	"errors"
)

// ParsePacket
func ParsePacket(payload []byte, offset int) (method, body []byte, err error) {
	mLen := int8(payload[offset])
	sOffset := offset + 1
	eOffset := sOffset + int(mLen)
	if eOffset > len(payload) {
		err = errors.New("Payload Not Enough")
		return
	}
	method = payload[sOffset:eOffset]
	body = payload[eOffset:]
	return
}

// MarshalPacket ...
func MarshalPacket(method, body []byte) ([]byte, []byte, []byte) {
	sizeBuf := []byte{byte(len(method))}
	return sizeBuf, method, body
}
