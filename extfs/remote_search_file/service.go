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

// Search sends a request to search for remote search files based on the given peer ID and search condition.
// It utilizes the RemoteSearchFileBroker to perform the request and returns a list of remote search file records
// or an error if the request fails.
//
// Parameters:
//   - peerId: The ID of the peer to which the search request is sent.
//   - condition: The search condition used to filter the remote search files.
//
// Returns:
//   - The total number of search results.
//   - A list of remote search file records.
//   - An ETag of the search result.
//   - An error if the search request fails.
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
		if record.Score > 0 {
			searchFile.Score = uint(record.Score)
		}
		searchFile.Tokens = record.Tokens
		searchFile.Available = record.Available

		searchFiles = append(searchFiles, searchFile)
	}

	return
}

// SearchForTopic searches for remote search files based on the given condition.
//
// Parameters:
//   - condition: The search condition used to filter the remote search files.
//
// Returns:
//   - A pointer to a RemoteSearchFileRecordList containing the search results.
//   - An error if the search request fails.
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
		record.Score = uint32(searchFile.Score)
		record.Tokens = searchFile.Tokens
		record.Available = searchFile.Available

		list.Records = append(list.Records, &record)
	}
	return &list, nil
}
