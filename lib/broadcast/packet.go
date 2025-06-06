// Define packet buffer for broadcast
//
// It is used to cache broadcast messages.
package broadcast

import (
	"bytes"
	"cmp"
	"context"
	"encoding/binary"
	"sync"
)

const CHECKSUM_THRESHOLD = uint16(65535)

type stdPacketBuffer struct {
	size    int
	content []byte
	addr    string
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// HashCode returns the address of the packet buffer as the hash code.
func (bpb *stdPacketBuffer) HashCode() string {
	return bpb.addr
}

func comparePacketBuffer(item *stdPacketBuffer, key string) int {
	return cmp.Compare(item.addr, key)
}

func parsePacketBuffer(block []byte) ([]byte, int) {
	if block[0] != 0 {
		return block, 0
	}
	offset := 1
	checksum := binary.BigEndian.Uint16(block[offset : offset+2])
	offset += 2
	size16 := binary.BigEndian.Uint16(block[offset : offset+2])
	offset += 2
	if checksum^size16 != CHECKSUM_THRESHOLD {
		return block, 0
	}

	size := int(size16)
	if size < len(block)-5 {
		return block, 0
	}
	return block[offset:], size
}

func packBuffer(buffer []byte) []byte {
	size := len(buffer)
	checksum := CHECKSUM_THRESHOLD ^ uint16(size)

	return bytes.Join([][]byte{
		binary.BigEndian.AppendUint16([]byte{0}, checksum),
		binary.BigEndian.AppendUint16(nil, uint16(size)),
		buffer,
	}, nil)
}
