package remotesearchfile

import (
	appnode "pan/app/app_node"
	nodesearchfile "pan/extfs/node_search_file"
	"time"
)

type RemoteSearchFileService struct {
	NodeSearchFileService  nodesearchfile.NodeSearchFileInternalService
	RemoteSearchFileBroker *RemoteSearchFileBroker
}

func (s *RemoteSearchFileService) Search(peerId string, condition RemoteSearchFileSearchCondition) (total int64, searchFiles []RemoteSearchFile, etag string, err error) {
	peerIdBytes, err := appnode.DecodePeerID(peerId)
	if err != nil {
		return
	}

	recordSearch := RemoteSearchFileRecordSearchCondition{
		Query:      condition.Query,
		Hash:       condition.Hash,
		RangeStart: int32(condition.RangeStart),
		RangeEnd:   int32(condition.RangeEnd),
	}

	recordList, err := s.RemoteSearchFileBroker.Search(peerIdBytes, &recordSearch)
	if err != nil {
		return
	}

	etag = recordList.Hash
	total = recordList.Total
	if total == 0 {
		return
	}

	for _, record := range recordList.Records {

		var searchFile RemoteSearchFile

		searchFile.ID = record.ID
		searchFile.ItemID = uint(record.ItemID)
		searchFile.Name = record.Name
		searchFile.FileType = record.FileType
		searchFile.FilePath = record.FilePath
		searchFile.MimeType = record.MimeType
		searchFile.Size = record.Size
		searchFile.CreatedAt = time.Unix(record.CreatedAt, 0)
		searchFile.UpdatedAt = time.Unix(record.UpdatedAt, 0)
		if len(record.Score) > 0 {
			searchFile.Score = int8(record.Score[0])
		}
		searchFile.Tokens = record.Tokens
		searchFile.Available = record.Available

		searchFiles = append(searchFiles, searchFile)
	}

	return
}

func (s *RemoteSearchFileService) SearchForTopic(condition *RemoteSearchFileRecordSearchCondition) (*RemoteSearchFileRecordList, error) {
	var condition_ nodesearchfile.NodeSearchFileSearchCondition

	condition_.Query = condition.Query
	condition_.Hash = condition.Hash
	condition_.RangeStart = int(condition.RangeStart)
	condition_.RangeEnd = int(condition.RangeEnd)

	total, searchFiles, _, err := s.NodeSearchFileService.Search(condition_)
	if err != nil {
		return nil, err
	}

	var list RemoteSearchFileRecordList
	list.Total = total
	if len(searchFiles) <= 0 {
		return &list, nil
	}

	for _, searchFile := range searchFiles {
		var record RemoteSearchFileRecord

		record.ID = searchFile.ID
		record.ItemID = uint32(searchFile.ItemID)
		record.Name = searchFile.Name
		record.FileType = searchFile.FileType
		record.FilePath = searchFile.FilePath
		record.MimeType = searchFile.MimeType
		record.Size = searchFile.Size
		record.CreatedAt = searchFile.CreatedAt.Unix()
		record.UpdatedAt = searchFile.UpdatedAt.Unix()
		record.Score = []byte{byte(searchFile.Score)}
		record.Tokens = searchFile.Tokens
		record.Available = searchFile.Available

		list.Records = append(list.Records, &record)
	}
	return &list, nil
}
