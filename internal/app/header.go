// Define context header with p2p handles
package app

import (
	"bytes"
	"encoding/binary"
	"io"
	"slices"
)

type HeaderItem struct {
	Key   []byte
	Value []byte
}

type Header struct {
	items []*HeaderItem
}

// Set sets the header item with the given key to the specified value.
// If an item with the same key exists, its value is updated.
// Otherwise, a new header item is inserted in the correct order.

func (h *Header) Set(key, value []byte) {
	idx, ok := slices.BinarySearchFunc(h.items, key, compareHeaderItem)
	if ok {
		h.items[idx].Value = value
		return
	}

	h.items = slices.Insert(h.items, idx, &HeaderItem{key, value})
}

// Get retrieves the value associated with the given key from the header.
// It returns the value and a boolean indicating whether the key was found.
func (h *Header) Get(key []byte) ([]byte, bool) {
	idx, ok := slices.BinarySearchFunc(h.items, key, compareHeaderItem)
	if !ok {
		return nil, ok
	}

	item := h.items[idx]
	return item.Value, ok
}

// Del removes the header item associated with the given key from the header.
// If the key is not found, Del does nothing.
func (h *Header) Del(key []byte) {
	idx, ok := slices.BinarySearchFunc(h.items, key, compareHeaderItem)
	if ok {
		h.items = slices.Delete(h.items, idx, idx+1)
	}
}

// MarshalHeader serializes the `Header` into an `io.Reader`.
// It returns a reader containing the binary data of the header items
// and the total size of the serialized header.
// Each header item is serialized by first appending the key size as a
// uint32, followed by the key bytes, then the value size as a uint32,
// and finally the value bytes. If the header contains no items, it
// returns a nil reader and a size of 0.

func MarshalHeader(header *Header) (io.Reader, int) {

	items := header.items
	if len(items) <= 0 {
		return nil, 0
	}

	buffer := make([]byte, 0)
	for _, item := range items {
		keySize := len(item.Key)
		valueSize := len(item.Value)
		buffer = binary.BigEndian.AppendUint32(buffer, uint32(keySize))
		buffer = append(buffer, item.Key...)
		buffer = binary.BigEndian.AppendUint32(buffer, uint32(valueSize))
		buffer = append(buffer, item.Value...)
	}
	return bytes.NewReader(buffer), len(buffer)
}

// UnmarshalHeader deserializes the header items from the given io.Reader
// and sets each item on the given Header.
// It reads the header items until an io.EOF is encountered.
// If any error occurs during the deserialization, it returns the error.
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

// InitHeader initializes the given Header by creating an empty slice of HeaderItem pointers.
// This function should be called before using the Header to ensure that the items slice is properly initialized.

func InitHeader(header *Header) {
	header.items = make([]*HeaderItem, 0)
}

// compareHeaderItem compares the key of a HeaderItem with a given byte slice.
// It returns an integer indicating the result of the comparison:
// a negative number if item.Key < key, zero if item.Key == key, and a positive number if item.Key > key.
func compareHeaderItem(item *HeaderItem, key []byte) int {
	return bytes.Compare(item.Key, key)
}
