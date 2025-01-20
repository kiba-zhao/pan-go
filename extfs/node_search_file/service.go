package nodesearchfile

import (
	"cmp"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"os"
	nodeitem "pan/extfs/node_item"
	"path"
	"slices"
	"time"
)

var ErrNodeSearchFileUnavailable = errors.New("nodesearchfile.NodeSearchFileService Error: Unavailable")

type NodeSearchFileInternalService interface {
	Search(condition NodeSearchFileSearchCondition) (total int64, searchFiles []NodeSearchFile, etag string, err error)
}

type NodeSearchFileService struct {
	NodeItemService      nodeitem.NodeItemInternalService
	NodeSearchFileRepo   NodeSearchFileRepository
	NodeSearchTaskRepo   NodeSearchTaskRepository
	NodeSearchTaskWorker NodeSearchTaskWorker
}

func (s *NodeSearchFileService) IsNotExist(err error) bool {
	return errors.Is(err, ErrNodeSearchTaskNotFound)
}

func (s *NodeSearchFileService) InitWithTaskID(taskId uint64) error {
	return s.NodeSearchFileRepo.Init(taskId)
}

func (s *NodeSearchFileService) DestroyWithTaskID(taskId uint64) error {
	return s.NodeSearchFileRepo.Destroy(taskId)
}

func (s *NodeSearchFileService) SaveWithTaskID(taskId uint64, searchFile NodeSearchFile) (NodeSearchFile, error) {
	return s.NodeSearchFileRepo.Save(taskId, searchFile)
}

func (s *NodeSearchFileService) Search(condition NodeSearchFileSearchCondition) (total int64, searchFiles []NodeSearchFile, etag string, err error) {

	var task NodeSearchTask
	var ok bool

	if len(condition.Hash) > 0 {
		task, err = s.NodeSearchTaskRepo.SelectWithQueryAndHash(condition.Query, condition.Hash)
	} else {
		task.Query = condition.Query
		task.Hash = genesearchFileTaskHash(task)
		task.Status = NodeSearchTaskStatusPending
		task, ok, err = s.NodeSearchTaskRepo.SelectOrCreate(task)
	}

	if err != nil {
		return
	}

	etag = task.Hash
	if err == nil && ok {
		s.NodeSearchTaskWorker.Reload()
		total = -1
		return
	}

	total, searchFiles, err = s.NodeSearchFileRepo.Search(task.ID, condition.RangeCondition)
	if err != nil {
		return
	}
	if task.Status == NodeSearchTaskStatusPending {
		total = -1
	}
	if len(searchFiles) <= 0 || len(condition.Hash) <= 0 {
		return
	}

	searchFiles_ := make([]NodeSearchFile, 0)
	items_ := make([]nodeitem.NodeItem, 0)
	for _, searchFile := range searchFiles {
		var searchFileErr error
		var item nodeitem.NodeItem
		itemIdx, ok := slices.BinarySearchFunc(items_, searchFile.ItemID, binarySearchFuncWithNodeItem)
		if !ok {
			item, searchFileErr = s.NodeItemService.Select(searchFile.ItemID)
			if searchFileErr == nil {
				items_ = slices.Insert(items_, itemIdx, item)
			}
		} else {
			item = items_[itemIdx]
		}

		if searchFileErr == nil {
			filepath := path.Join(item.FilePath, searchFile.FilePath)
			stat, searchFileErr := os.Stat(filepath)
			if searchFileErr == nil {
				if stat.IsDir() {
					searchFile.Available = searchFile.FileType == nodeitem.FileTypeFolder
				} else {
					searchFile.Available = searchFile.FileType == nodeitem.FileTypeFile
				}
			}
		}

		if searchFileErr != nil {
			searchFile.Available = false
		}
		searchFiles_ = append(searchFiles_, searchFile)
	}

	searchFiles = searchFiles_
	return
}

func binarySearchFuncWithNodeItem(item nodeitem.NodeItem, target uint) int {
	return cmp.Compare(item.ID, target)
}

func genesearchFileTaskHash(task NodeSearchTask) string {
	hashBytes := make([]byte, 16)

	binary.BigEndian.PutUint64(hashBytes, uint64(time.Now().Unix()))
	rand.Read(hashBytes[8:])

	return base64.StdEncoding.EncodeToString(hashBytes)
}
