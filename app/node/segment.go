package node

import (
	"encoding/binary"
	"io"
)

func ParseSegment(reader io.Reader) (segment []byte, err error) {

	unsignedSize := uint32(0)
	err = binary.Read(reader, binary.BigEndian, &unsignedSize)
	size := int(unsignedSize)

	segment = make([]byte, size)
	total := 0
	buffer := segment
	for {
		num, err := reader.Read(buffer)
		if err != nil && num != size {
			break
		}
		total += num
		if total >= size {
			break
		}
		buffer = segment[:total]
	}

	return
}
