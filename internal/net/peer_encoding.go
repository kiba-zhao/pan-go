package net

import "encoding/base64"

func EncodePeerID(peerId PeerID) string {
	return base64.RawURLEncoding.EncodeToString(peerId)
}

func DecodePeerID(peerId string) (PeerID, error) {
	return base64.RawURLEncoding.DecodeString(peerId)
}
