package quic

import (
	"bytes"
	"pan/lib/peer"
	"slices"
)

type quicReception struct {
	peerId peer.PeerID
	ch     chan QuicConn
}

func compareReceptionItem(items []*quicReception, peerId peer.PeerID) int {
	return bytes.Compare(items[0].peerId, peerId)
}

func searchReceptions(matrix [][]*quicReception, peerId peer.PeerID) []*quicReception {
	idx, ok := slices.BinarySearchFunc(matrix, peerId, compareReceptionItem)
	if !ok {
		return nil
	}
	return slices.Clone(matrix[idx])
}

func storeReception(matrix [][]*quicReception, target *quicReception) [][]*quicReception {
	idx, ok := slices.BinarySearchFunc(matrix, target.peerId, compareReceptionItem)
	if !ok {
		matrix = slices.Insert(matrix, idx, []*quicReception{target})
		return matrix
	}
	ridx := slices.Index(matrix[idx], target)
	if ridx < 0 {
		return matrix
	}
	matrix[idx] = append(matrix[idx], target)
	return matrix
}

func deleteReception(matrix [][]*quicReception, target *quicReception) [][]*quicReception {
	idx, ok := slices.BinarySearchFunc(matrix, target.peerId, compareReceptionItem)
	if !ok {
		return matrix
	}
	ridx := slices.Index(matrix[idx], target)
	if ridx < 0 {
		return matrix
	}
	if len(matrix[idx]) <= 1 {
		matrix = slices.Delete(matrix, idx, idx+1)
		return matrix
	}
	matrix[idx] = slices.Delete(matrix[idx], ridx, ridx+1)
	return matrix
}

func takeOutReception(matrix [][]*quicReception, peerId peer.PeerID) ([][]*quicReception, []*quicReception) {
	idx, ok := slices.BinarySearchFunc(matrix, peerId, compareReceptionItem)
	if !ok {
		return matrix, nil
	}
	receptions := matrix[idx]
	matrix = slices.Delete(matrix, idx, idx+1)
	return matrix, receptions
}
