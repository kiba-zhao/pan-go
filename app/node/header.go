package node

import (
	"bytes"
	"encoding/binary"
	"io"
	"slices"
)

type HeaderItem struct {
	key   []byte
	value []byte
}

type Header struct {
	items []*HeaderItem
}

func (h *Header) Set(key, value []byte) {
	idx, ok := slices.BinarySearchFunc(h.items, key, compareHeaderItem)
	if ok {
		h.items[idx].value = value
		return
	}

	h.items = slices.Insert(h.items, idx, &HeaderItem{key, value})
}

func (h *Header) Get(key []byte) ([]byte, bool) {
	idx, ok := slices.BinarySearchFunc(h.items, key, compareHeaderItem)
	if !ok {
		return nil, ok
	}

	item := h.items[idx]
	return item.value, ok
}

func (h *Header) Del(key []byte) {
	idx, ok := slices.BinarySearchFunc(h.items, key, compareHeaderItem)
	if ok {
		h.items = slices.Delete(h.items, idx, idx+1)
	}
}

func MarshalHeader(header *Header) (io.Reader, int) {

	items := header.items
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
	header.items = make([]*HeaderItem, 0)
}

func compareHeaderItem(item *HeaderItem, key []byte) int {
	return bytes.Compare(item.key, key)
}
