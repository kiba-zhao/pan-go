package quic

import (
	"bytes"
	"pan/lib/peer"
	"slices"
)

type stdQuicReception struct {
	peerId peer.PeerID
	ch     chan QuicConn
}

func compareReceptionItem(items []*stdQuicReception, peerId peer.PeerID) int {
	return bytes.Compare(items[0].peerId, peerId)
}

func searchReceptions(matrix [][]*stdQuicReception, peerId peer.PeerID) []*stdQuicReception {
	idx, ok := slices.BinarySearchFunc(matrix, peerId, compareReceptionItem)
	if !ok {
		return nil
	}
	return slices.Clone(matrix[idx])
}

func storeReception(matrix [][]*stdQuicReception, target *stdQuicReception) [][]*stdQuicReception {
	idx, ok := slices.BinarySearchFunc(matrix, target.peerId, compareReceptionItem)
	if !ok {
		matrix = slices.Insert(matrix, idx, []*stdQuicReception{target})
		return matrix
	}
	ridx := slices.Index(matrix[idx], target)
	if ridx < 0 {
		return matrix
	}
	matrix[idx] = append(matrix[idx], target)
	return matrix
}

func deleteReception(matrix [][]*stdQuicReception, target *stdQuicReception) [][]*stdQuicReception {
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

func takeOutReception(matrix [][]*stdQuicReception, peerId peer.PeerID) ([][]*stdQuicReception, []*stdQuicReception) {
	idx, ok := slices.BinarySearchFunc(matrix, peerId, compareReceptionItem)
	if !ok {
		return matrix, nil
	}
	receptions := matrix[idx]
	matrix = slices.Delete(matrix, idx, idx+1)
	return matrix, receptions
}
