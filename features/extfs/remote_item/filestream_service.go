package remoteitem

import (
	"io"
	"os"
	nodeitem "pan/features/extfs/node_item"
	"pan/lib/peer"
)

type RemoteFileStreamService struct {
	RemoteFileInfoService  RemoteFileInfoInternalService
	FilePathService        nodeitem.NodeFilePathInternalService
	RemoteFileStreamBroker *RemoteFileStreamBroker
}

func (s *RemoteFileStreamService) IsNotExist(err error) bool {
	if s.FilePathService.IsNotExist(err) {
		return true
	}
	return s.RemoteFileInfoService.IsNotExist(err)
}

func (s *RemoteFileStreamService) Read(peerId string, itemId uint, filePath string) (*RemoteFileStreamReader, error) {

	peerIdBytes, err := peer.DecodePeerID(peerId)
	if err != nil {
		return nil, err
	}

	return s.ReadWithPeerID(peerIdBytes, itemId, filePath)
}

func (s *RemoteFileStreamService) ReadWithPeerID(peerId []byte, itemId uint, filePath string) (*RemoteFileStreamReader, error) {

	var condition RemoteFileInfoRecordSelectCondition
	condition.ItemID = uint32(itemId)
	condition.FilePath = filePath

	record, err := s.RemoteFileInfoService.SelectWithCondition(peerId, &condition)
	if err != nil {
		return nil, err
	}

	var reader RemoteFileStreamReader
	reader.peerId = peerId
	reader.service = s
	reader.record = record
	return &reader, nil
}

func (s *RemoteFileStreamService) SelectForTopic(condition *RemoteFileStreamSelectCondition) (io.ReadCloser, error) {
	filePath, err := s.FilePathService.SelectWithoutFolder(uint(condition.ItemID), condition.FilePath)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	if condition.Offset > 0 {
		_, err = file.Seek(condition.Offset, 0)
	}
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *RemoteFileStreamService) SelectWithCondition(peerId peer.PeerID, condition *RemoteFileStreamSelectCondition) (io.ReadCloser, error) {
	return s.RemoteFileStreamBroker.Select(peerId, condition)
}
