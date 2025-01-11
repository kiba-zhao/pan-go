package nodefile

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"os"
	"path"

	nodeitem "pan/extfs/node_item"
)

var ErrNodeFileUnavailable = errors.New("nodefile.NodeFileService Error: Unavailable")
var ErrNodeFileInvalidNodeItem = errors.New("nodefile.NodeFileService Error: Invalid Node Item")
var ErrNodeFileInvalidID = errors.New("nodefile.NodeFileService Error: Invalid Node File ID")

type NodeFileInternalService interface {
	TraverseWithCondition(func(item NodeFile) error, NodeFileSearchCondition) error
	SelectWithCondition(NodeFileSelectCondition) (NodeFile, error)
	IsNotExist(err error) bool
}

type NodeFileService struct {
	NodeItemService nodeitem.NodeItemInternalService
}

func (s *NodeFileService) IsNotExist(err error) bool {
	if os.IsNotExist(err) {
		return true
	}
	if errors.Is(err, ErrNodeFileUnavailable) {
		return true
	}
	return s.NodeItemService.IsNotExist(err)
}

func (s *NodeFileService) Select(id string) (NodeFile, error) {
	itemId, filePath, err := ParseNodeFileID(id)
	if err != nil {
		return NodeFile{}, err
	}
	nodeItem, err := s.NodeItemService.Select(itemId)
	if err != nil {
		return NodeFile{}, err
	}
	if nodeItem.FileType != nodeitem.FileTypeFolder {
		return NodeFile{}, ErrNodeFileInvalidNodeItem
	}
	if !nodeItem.Available {
		return NodeFile{}, ErrNodeFileUnavailable
	}

	realFilePath := path.Join(nodeItem.FilePath, filePath)

	fileStat, err := os.Stat(realFilePath)
	if err != nil {
		return NodeFile{}, err
	}

	var fileItem NodeFile
	fileItem.ID = id
	fileItem.ItemID = nodeItem.ID
	fileItem.Name = fileStat.Name()
	if fileStat.IsDir() {
		fileItem.FileType = nodeitem.FileTypeFolder
	} else {
		fileItem.FileType = nodeitem.FileTypeFile
		setMimeType(&fileItem, realFilePath)
	}
	fileItem.ParentPath = path.Dir(filePath)
	fileItem.FilePath = filePath
	fileItem.Size = fileStat.Size()
	fileItem.UpdatedAt = fileStat.ModTime()
	fileItem.CreatedAt = fileStat.ModTime()
	fileItem.Available = true

	return fileItem, nil
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
		setMimeType(&fileItem, filePath)
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
	fileItem.ID = GenerateNodeFileID(fileItem.ItemID, fileItem.FilePath)

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
			setMimeType(&item, path.Join(filePath, item.Name))
		}

		if conditions.ParentPath == nil {
			item.FilePath = item.Name
		} else {
			item.ParentPath = *conditions.ParentPath
			item.FilePath = path.Join(*conditions.ParentPath, item.Name)
		}
		item.ID = GenerateNodeFileID(item.ItemID, item.FilePath)

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

func GenerateNodeFileID(itemId uint, filePath string) string {
	filePathBytes := []byte(filePath)
	idBytes := make([]byte, 4+len(filePathBytes))
	binary.BigEndian.PutUint32(idBytes, uint32(itemId))
	copy(idBytes[4:], filePathBytes)
	return base64.StdEncoding.EncodeToString(idBytes)
}

func ParseNodeFileID(id string) (uint, string, error) {
	idBytes, err := base64.StdEncoding.DecodeString(id)
	if err == nil && len(idBytes) < 4 {
		err = ErrNodeFileInvalidID
	}
	if err != nil {
		return 0, "", err
	}
	itemId := binary.BigEndian.Uint32(idBytes)
	filePath := string(idBytes[4:])
	return uint(itemId), filePath, err
}

func setMimeType(fileItem *NodeFile, filePath string) error {
	mimeType, err := nodeitem.GenerateMimeType(filePath)
	if err == nil {
		fileItem.MimeType = mimeType
	}
	return err
}
