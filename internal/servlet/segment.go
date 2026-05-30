package servlet

import (
	"encoding/binary"
	"io"
)

// ParseSegment reads a segment from the provided io.Reader.
// It first reads a uint32 value to determine the size of the segment,
// then reads the bytes of the segment based on the determined size.
// It returns the segment as a byte slice and any error encountered
// during the read operations.

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
