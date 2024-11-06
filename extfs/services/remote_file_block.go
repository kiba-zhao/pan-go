package services

import (
	"bytes"
	"io"
	"os"
	appConstant "pan/app/constant"
	appNode "pan/app/node"
	"pan/extfs/models"
	"path"

	"google.golang.org/protobuf/proto"
)

type RemoteFileReader struct {
	file   *os.File
	reader io.Reader
}

func (r *RemoteFileReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if err == io.EOF || n < len(p) {
		r.file.Close()
	}
	return
}

type RemoteFileBlockService struct {
	NodeItemService NodeItemInternalService
	NodeModule      appNode.NodeModule
	NodeScopeModule appNode.NodeScopeModule
}

func (s *RemoteFileBlockService) SelectForNode(condition *models.RemoteFileBlockSelectCondition) (io.Reader, error) {
	nodeItem, err := s.NodeItemService.Select(uint(condition.ItemID))
	if err != nil {
		return nil, err
	}
	if nodeItem.FileType != FileTypeFile {
		return nil, appConstant.ErrRefused
	}

	filePath := nodeItem.FilePath
	if len(condition.ParentPath) > 0 {
		filePath = path.Join(filePath, condition.ParentPath)
	}
	if len(condition.Name) > 0 {
		filePath = path.Join(filePath, condition.Name)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	if condition.Offset > 0 {
		file.Seek(condition.Offset, 0)
	}

	var reader io.Reader
	if condition.Limit > 0 {
		reader = io.LimitReader(file, condition.Limit)
	} else {
		reader = file
	}

	return &RemoteFileReader{file: file, reader: reader}, nil
}

var RequestRemoteFileBlock = []byte("select_remote_file_block")

func (s *RemoteFileBlockService) SelectWithCondition(nodeId appNode.NodeID, condition *models.RemoteFileBlockSelectCondition) (io.Reader, error) {
	requestBytes, err := proto.Marshal(condition)
	if err != nil {
		return nil, err
	}

	scope := s.NodeScopeModule.NodeScope()
	requestName := appNode.GenerateRouteName(scope, RequestRemoteFileBlock)
	request := appNode.NewRequest(requestName, bytes.NewReader(requestBytes))
	response, err := s.NodeModule.Do(nodeId, request)
	if err != nil {
		return nil, err
	}

	if response.Code() != appConstant.CodeOK {
		return nil, appConstant.ErrInternalError
	}

	return response.Body(), nil
}
