package models

import "pan/peer"

type PeerStat struct {
	ID   string         `json:"id"`
	Stat peer.PeerState `json:"stat"`
}
