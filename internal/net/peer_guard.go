package net

type PeerGuard interface {
	AllowAccess(PeerID) bool
}
