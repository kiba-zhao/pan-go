package services

import (
	"bytes"
	"io"
	"os"
	appConstant "pan/app/constant"
	appNode "pan/app/node"
	"pan/extfs/models"
	"path"
	"sync"

	"google.golang.org/protobuf/proto"
)

type RemoteFileReader struct {
	filePath string
	offset   int64
	limit    int64
	file     *os.File
	reader   io.Reader
	locker   sync.Mutex
	closed   bool
}

func (r *RemoteFileReader) openFile() error {
	if r.closed {
		return io.ErrClosedPipe
	}
	file, err := os.Open(r.filePath)
	if err != nil {
		return err
	}
	if r.offset > 0 {
		file.Seek(r.offset, 0)
	}
	if r.limit > 0 {
		r.reader = io.LimitReader(file, r.limit)
	} else {
		r.reader = file
	}
	r.file = file
	return nil
}

func (r *RemoteFileReader) closeFile() error {
	file := r.file
	if r.closed || file == nil {
		return nil
	}
	r.closed = true
	r.file = nil
	return file.Close()
}

func (r *RemoteFileReader) Read(p []byte) (n int, err error) {
	r.locker.Lock()
	defer r.locker.Unlock()
	if r.file == nil {
		err = r.openFile()
	}
	if err != nil {
		return
	}

	n, err = r.reader.Read(p)
	if err != nil {
		r.closeFile()
	}
	return
}

func (r *RemoteFileReader) Close() error {
	r.locker.Lock()
	defer r.locker.Unlock()
	return r.closeFile()
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

	filePath := nodeItem.FilePath
	if nodeItem.FileType == FileTypeFolder {
		if len(condition.ParentPath) > 0 {
			filePath = path.Join(filePath, condition.ParentPath)
		}
		if len(condition.Name) > 0 {
			filePath = path.Join(filePath, condition.Name)
		}
	}

	_, err = os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &RemoteFileReader{filePath: filePath, offset: condition.Offset, limit: condition.Limit}, nil
}

var RequestRemoteFileBlock = []byte("select_remote_file_block")

func (s *RemoteFileBlockService) SelectWithCondition(nodeId appNode.NodeID, condition *models.RemoteFileBlockSelectCondition) (io.ReadCloser, error) {
	requestBytes, err := proto.Marshal(condition)
	if err != nil {
		return nil, err
	}

	scope := s.NodeScopeModule.NodeScope()
	requestName := appNode.GenerateRouteName(scope, RequestRemoteFileBlock)
	request := appNode.NewRequest(requestName, bytes.NewReader(requestBytes))
	res, err := s.NodeModule.Do(nodeId, request)
	if err != nil {
		return nil, err
	}

	if res.Code() != appConstant.CodeOK {
		return nil, appConstant.ErrInternalError
	}

	return res, nil
}
