package peer

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	HeaderSegmentType = uint8(iota)
	BodySegmentType
)

type HeaderSegment struct {
	name  []byte
	value []byte
}

// Name ...
func (hs *HeaderSegment) Name() []byte {
	return hs.name
}

// Value ...
func (hs *HeaderSegment) Value() []byte {
	return hs.value
}

// NewHeaderSegment ...
func NewHeaderSegment(name, value []byte) *HeaderSegment {
	segment := new(HeaderSegment)
	segment.name = name
	segment.value = value
	return segment
}

// ParseSegmentType ...
func ParseSegmentType(reader io.Reader) (stype uint8, err error) {
	buf := make([]byte, 1)
	_, err = reader.Read(buf)
	stype = buf[0]
	return
}

// ParseHeaderSegment ...
func ParseHeaderSegment(reader io.Reader) (header *HeaderSegment, err error) {

	name, err := ParseHeaderSegmentField(reader)
	if err != nil {
		return
	}
	value, err := ParseHeaderSegmentField(reader)
	if err == nil {
		header = new(HeaderSegment)
		header.name = name
		header.value = value
	}
	return
}

// ParseHeaderSegmentField ...
func ParseHeaderSegmentField(reader io.Reader) (data []byte, err error) {

	usize := uint16(0)
	err = binary.Read(reader, binary.BigEndian, &usize)
	size := int(usize)
	if size <= 0 {
		return
	}

	data = make([]byte, size)
	total := 0
	buffer := data
	for {
		num, err := reader.Read(buffer)
		if err != nil {
			break
		}
		total += num
		if total >= size {
			break
		}
		buffer = data[:total]
	}

	return
}

// CreateSegmentType ...
func CreateSegmentType(segmentType uint8) io.Reader {
	reader := bytes.NewReader([]byte{segmentType})
	return reader
}

// CreateHeaderSegment ...
func CreateHeaderSegment(header *HeaderSegment) []io.Reader {

	readers := make([]io.Reader, 0)
	nsr, ndr := CreateHeaderSegmentField(header.Name())
	vsr, vdr := CreateHeaderSegmentField(header.Value())
	readers = append(readers, nsr, ndr, vsr, vdr)
	return readers
}

// CreateHeaderSegmentField ...
func CreateHeaderSegmentField(data []byte) (sr io.Reader, dr io.Reader) {

	size := len(data)
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(size))
	sr = bytes.NewReader(b)
	dr = bytes.NewReader(data)

	return
}
