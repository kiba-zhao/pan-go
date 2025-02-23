package nodeitem

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrNodeFilePathUnavailable = errors.New("nodeitem.NodeFilePathService Error: Unavailable")
var ErrNodeFilePathWithoutFolder = errors.New("nodeitem.NodeFilePathService Error: Without Folder")

type NodeFilePathInternalService interface {
	IsNotExist(err error) bool
	Select(id uint, filePath string) (string, error)
	SelectWithoutFolder(id uint, filePath string) (string, error)
}

type NodeFilePathService struct {
	NodeItemService NodeItemInternalService
}

func (s *NodeFilePathService) IsNotExist(err error) bool {
	return s.NodeItemService.IsNotExist(err)
}

func (s *NodeFilePathService) Select(id uint, filePath string) (string, error) {
	nodeItem, err := s.NodeItemService.Select(id)
	if err != nil {
		return "", err
	}
	if !nodeItem.Available {
		return "", ErrNodeFilePathUnavailable
	}

	if len(filePath) <= 0 {
		return nodeItem.FilePath, nil
	}

	filePath_ := filepath.FromSlash(filePath)
	realFilePath := filepath.Join(nodeItem.FilePath, filePath_)
	return realFilePath, err

}

func (s *NodeFilePathService) SelectWithoutFolder(id uint, filePath string) (string, error) {
	realFilePath, err := s.Select(id, filePath)
	if err != nil {
		return "", err
	}
	fileStat, err := os.Stat(realFilePath)
	if err != nil {
		return realFilePath, err
	}
	if fileStat.IsDir() {
		return "", ErrNodeFilePathWithoutFolder
	}
	return realFilePath, err
}
