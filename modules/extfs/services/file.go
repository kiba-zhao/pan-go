package services

import (
	extfs "pan/modules/extfs/peer"
	"pan/peer"
)

type FileService struct {
	API extfs.API
}

// SyncRemoteFile ...
func (s *FileService) SyncRemoteFiles(peerId peer.PeerId) {

}
