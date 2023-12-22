package broadcast

import (
	"bytes"
	"errors"
	"net"

	"pan/core"
)

// Accept ...
func Accept(app core.App[Context], network Net) {
	var addr []byte
	var payload []byte
	bufferSize := -1
	for {
		msg, srcAddr, err := network.Read(bufferSize)
		if errors.Is(err, net.ErrClosed) {
			break
		}
		if err != nil {
			panic(err)
		}

		if payload == nil || !bytes.Equal(addr, srcAddr) {
			payload = msg
		} else {
			payload = bytes.Join([][]byte{payload, msg}, nil)
		}

		packet, size, err := ParsePacket(payload)
		if size > 0 && err != nil {
			addr = srcAddr
			bufferSize = size - len(payload)
			continue
		}

		if size <= 0 {
			payload = nil
			addr = nil
			bufferSize = -1
			continue
		}

		method, body, err := core.ParsePacket(packet, 0)
		if err == nil {
			ctx := NewContext(method, body, srcAddr, network)
			go app.Run(ctx)
			if len(payload) > size {
				payload = payload[size:]
				continue
			}
		}
		payload = nil
		addr = nil
		bufferSize = -1

	}

}

// Dispatch ...
func Dispatch(method, body []byte, network Net) (err error) {
	s, m, b := core.MarshalPacket(method, body)
	payload := MarshalPacket(s, m, b)
	err = network.Write(payload)
	return
}
