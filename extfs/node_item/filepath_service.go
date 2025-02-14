package nodeitem

import (
	"errors"
	"path/filepath"
)

var ErrNodeFilePathUnavailable = errors.New("nodeitem.NodeFilePathService Error: Unavailable")

type NodeFilePathInternalService interface {
	IsNotExist(err error) bool
	Select(id uint, filePath string) (string, error)
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
