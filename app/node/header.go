package node

import (
	"bytes"
	"encoding/binary"
	"io"
	"pan/app/cache"
)

type HeaderItem struct {
	key   []byte
	value []byte
}

func (item *HeaderItem) HashCode() []byte {
	return item.key
}

type Header struct {
	bucket cache.Bucket[[]byte, *HeaderItem]
}

func (h *Header) Set(key, value []byte) {
	h.bucket.Swap(&HeaderItem{key, value})
}

func (h *Header) Get(key []byte) ([]byte, bool) {
	item, ok := h.bucket.Search(key)
	if !ok {
		return nil, ok
	}
	return item.value, ok
}

func (h *Header) Del(key []byte) {
	item, ok := h.bucket.Search(key)
	if ok {
		h.bucket.Delete(item)
	}
}

func MarshalHeader(header *Header) (io.Reader, int) {

	items := header.bucket.Items()
	if len(items) <= 0 {
		return nil, 0
	}

	buffer := make([]byte, 0)
	for _, item := range items {
		keySize := len(item.key)
		valueSize := len(item.value)
		buffer = binary.BigEndian.AppendUint32(buffer, uint32(keySize))
		buffer = append(buffer, item.key...)
		buffer = binary.BigEndian.AppendUint32(buffer, uint32(valueSize))
		buffer = append(buffer, item.value...)
	}
	return bytes.NewReader(buffer), len(buffer)
}

func UnmarshalHeader(reader io.Reader, header *Header) error {
	for {
		var value []byte
		key, err := ParseSegment(reader)
		if err == nil {
			value, err = ParseSegment(reader)
		}
		if err != nil {
			return err
		}
		header.Set(key, value)
	}
}

func InitHeader(header *Header) {
	header.bucket = cache.NewBucket[[]byte, *HeaderItem](bytes.Compare)
}
