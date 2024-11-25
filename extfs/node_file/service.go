package nodefile

import (
	"encoding/base64"
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	nodeitem "pan/extfs/node_item"
)

var ErrNodeFileUnavailable = errors.New("nodefile.NodeFileService Error: Unavailable")
var ErrNodeFileInvalidNodeItem = errors.New("nodefile.NodeFileService Error: Invalid Node Item")

type NodeFileInternalService interface {
	TraverseWithCondition(func(item NodeFile) error, NodeFileSearchCondition) error
	SelectWithCondition(NodeFileSelectCondition) (NodeFile, error)
}

type NodeFileService struct {
	NodeItemService nodeitem.NodeItemInternalService
}

func (s *NodeFileService) Search(conditions NodeFileSearchCondition) (int64, []NodeFile, error) {

	items := make([]NodeFile, 0)
	err := s.TraverseWithCondition(func(item NodeFile) error {
		items = append(items, item)
		return nil
	}, conditions)
	return int64(len(items)), items, err
}

func (s *NodeFileService) SelectWithCondition(condition NodeFileSelectCondition) (NodeFile, error) {
	nodeItem, err := s.NodeItemService.Select(condition.ItemID)
	if err != nil {
		return NodeFile{}, err
	}

	if nodeItem.FileType != nodeitem.FileTypeFolder {
		return NodeFile{}, ErrNodeFileInvalidNodeItem
	}

	if !nodeItem.Available {
		return NodeFile{}, ErrNodeFileUnavailable
	}

	filePath := nodeItem.FilePath
	if len(condition.ParentPath) > 0 {
		filePath = path.Join(filePath, condition.ParentPath)
	}
	filePath = path.Join(filePath, condition.Name)

	fileStat, err := os.Stat(filePath)
	if err != nil {
		return NodeFile{}, err
	}

	var fileItem NodeFile

	fileItem.ItemID = nodeItem.ID
	fileItem.Name = fileStat.Name()
	if fileStat.IsDir() {
		fileItem.FileType = nodeitem.FileTypeFolder
	} else {
		fileItem.FileType = nodeitem.FileTypeFile
	}
	fileItem.ParentPath = condition.ParentPath
	if len(fileItem.ParentPath) > 0 {
		fileItem.FilePath = path.Join(fileItem.ParentPath, fileStat.Name())
	} else {
		fileItem.FilePath = fileStat.Name()
	}
	fileItem.Size = fileStat.Size()
	fileItem.Available = true
	fileItem.CreatedAt = fileStat.ModTime()
	fileItem.UpdatedAt = fileStat.ModTime()
	fileItem.ID = generateNodeFileID(fileItem.ItemID, fileItem.FilePath)

	return fileItem, err

}

func (s *NodeFileService) TraverseWithCondition(traverseFn func(item NodeFile) error, conditions NodeFileSearchCondition) error {
	nodeItem, err := s.NodeItemService.Select(conditions.ItemID)
	if err != nil {
		return err
	}

	if nodeItem.FileType != nodeitem.FileTypeFolder {
		return ErrNodeFileInvalidNodeItem
	}

	if !nodeItem.Available {
		return ErrNodeFileUnavailable
	}

	filePath := nodeItem.FilePath
	if conditions.ParentPath != nil {
		filePath = path.Join(filePath, *conditions.ParentPath)
	}

	files, err := os.ReadDir(filePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		var item NodeFile

		item.ItemID = nodeItem.ID
		item.Name = file.Name()

		if file.IsDir() {
			item.FileType = nodeitem.FileTypeFolder
		} else {
			item.FileType = nodeitem.FileTypeFile
		}

		if conditions.ParentPath == nil {
			item.FilePath = item.Name
		} else {
			item.ParentPath = *conditions.ParentPath
			item.FilePath = path.Join(*conditions.ParentPath, item.Name)
		}
		item.ID = generateNodeFileID(item.ItemID, item.FilePath)

		item.Available = true
		info, infoErr := file.Info()
		if infoErr == nil {
			item.UpdatedAt = info.ModTime()
			item.Size = info.Size()
		} else {
			item.Available = false
		}
		err = traverseFn(item)
		if err != nil {
			return err
		}
	}
	return err
}

const NodeFileSep = "_"

func generateNodeFileID(itemId uint, filePath string) string {
	idStr := strings.Join([]string{strconv.FormatUint(uint64(itemId), 10), filePath}, NodeFileSep)
	return base64.StdEncoding.EncodeToString([]byte(idStr))
}
