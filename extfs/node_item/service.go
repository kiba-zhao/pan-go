package nodeitem

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

const (
	FileTypeFolder = "D"
	FileTypeFile   = "F"
)

type NodeItemInternalService interface {
	TraverseAll(func(NodeItem) error) error
	Select(uint) (NodeItem, error)
	SelectByName(string) (NodeItem, error)
	IsNotExist(error) bool
	SelectAllWithEnabled(bool) ([]NodeItem, error)
}

type NodeItemService struct {
	NodeItemRepo NodeItemRepository
}

func (s *NodeItemService) IsNotExist(err error) bool {
	return errors.Is(err, ErrNodeItemNotFound)
}

func (s *NodeItemService) SelectAll() (int64, []NodeItem, error) {
	items := make([]NodeItem, 0)
	err := s.TraverseAll(func(nodeItem NodeItem) error {
		items = append(items, nodeItem)
		return nil
	})

	if err != nil {
		return 0, nil, err
	}
	return int64(len(items)), items, nil
}

func (s *NodeItemService) Create(fields NodeItemFields) (NodeItem, error) {
	stat, err := os.Stat(fields.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return NodeItem{}, err
	}
	available := err == nil

	var nodeItem NodeItem
	nodeItem.Name = fields.Name
	nodeItem.FilePath = fields.FilePath
	nodeItem.Enabled = fields.Enabled
	nodeItem.Available = available
	if stat.IsDir() {
		nodeItem.FileType = FileTypeFolder
	} else {
		nodeItem.FileType = FileTypeFile
		mimeType, err := GenerateMimeType(nodeItem.FilePath)
		if err == nil {
			nodeItem.MimeType = mimeType
		}
	}

	nodeItem_, err := s.NodeItemRepo.Save(nodeItem)
	if err == nil {
		setNodeItemAvailableWithFileStat(&nodeItem_)
	}
	return nodeItem_, err
}

func (s *NodeItemService) Update(fields NodeItemFields, id uint) (NodeItem, error) {
	nodeItem, err := s.NodeItemRepo.Select(id)
	if err != nil {
		return nodeItem, err
	}

	stat, err := os.Stat(fields.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return nodeItem, err
	}
	err = nil

	nodeItem.Name = fields.Name
	nodeItem.FilePath = fields.FilePath
	nodeItem.Enabled = fields.Enabled
	if stat.IsDir() {
		nodeItem.FileType = FileTypeFolder
	} else {
		nodeItem.FileType = FileTypeFile
		mimeType, err := GenerateMimeType(nodeItem.FilePath)
		if err == nil {
			nodeItem.MimeType = mimeType
		}
	}

	nodeItem, err = s.NodeItemRepo.Save(nodeItem)
	if err == nil {
		setNodeItemAvailableWithFileStat(&nodeItem)
	}
	return nodeItem, err
}

func (s *NodeItemService) Select(id uint) (NodeItem, error) {
	nodeItem, err := s.NodeItemRepo.Select(id)
	if err != nil {
		return nodeItem, err
	}
	setNodeItemAvailableWithFileStat(&nodeItem)
	setNodeItemAvailableWithMimeType(&nodeItem)
	return nodeItem, nil
}

func (s *NodeItemService) SelectByName(name string) (NodeItem, error) {
	nodeItem, err := s.NodeItemRepo.SelectByName(name)
	if err != nil {
		return nodeItem, err
	}
	setNodeItemAvailableWithFileStat(&nodeItem)
	setNodeItemAvailableWithMimeType(&nodeItem)
	return nodeItem, nil
}

func (s *NodeItemService) TraverseAll(traverseFn func(NodeItem) error) error {
	return s.NodeItemRepo.TraverseAll(func(nodeItem NodeItem) error {
		setNodeItemAvailableWithFileStat(&nodeItem)
		setNodeItemAvailableWithMimeType(&nodeItem)
		return traverseFn(nodeItem)
	})
}

func (s *NodeItemService) SelectAllWithEnabled(enabled bool) ([]NodeItem, error) {
	return s.NodeItemRepo.SelectAllWithEnabled(enabled)
}

func (s *NodeItemService) Delete(id uint) error {
	nodeItem, err := s.NodeItemRepo.Select(id)
	if err != nil {
		return err
	}
	err = s.NodeItemRepo.Delete(nodeItem)
	return err
}

func setNodeItemAvailableWithFileStat(nodeItem *NodeItem) {
	nodeItem.Available = *nodeItem.Enabled
	if !nodeItem.Available {
		return
	}

	stat, err := os.Stat(nodeItem.FilePath)
	nodeItem.Available = err == nil
	if !nodeItem.Available {
		return
	}

	nodeItem.Size = stat.Size()
	if stat.IsDir() {
		nodeItem.Available = nodeItem.FileType == FileTypeFolder
	} else {
		nodeItem.Available = nodeItem.FileType == FileTypeFile
	}
}

func setNodeItemAvailableWithMimeType(nodeItem *NodeItem) {
	if strings.Compare(nodeItem.FileType, FileTypeFile) != 0 {
		return
	}
	if !nodeItem.Available {
		return
	}
	mimeType, _ := GenerateMimeType(nodeItem.FilePath)
	nodeItem.Available = strings.Compare(nodeItem.MimeType, mimeType) == 0
}

func GenerateMimeType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	fileType := http.DetectContentType(buffer)
	return fileType, nil
}
