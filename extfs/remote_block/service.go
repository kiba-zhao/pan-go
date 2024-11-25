package remoteblock

import (
	"io"
	"os"
	"pan/app/peer"
	"path"
	sync "sync"

	nodeitem "pan/extfs/node_item"
)

type RemoteBlockReader struct {
	filePath string
	offset   int64
	limit    int64
	file     *os.File
	reader   io.Reader
	locker   sync.Mutex
	closed   bool
}

func (r *RemoteBlockReader) openFile() error {
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

func (r *RemoteBlockReader) closeFile() error {
	file := r.file
	if r.closed || file == nil {
		return nil
	}
	r.closed = true
	r.file = nil
	return file.Close()
}

func (r *RemoteBlockReader) Read(p []byte) (n int, err error) {
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

func (r *RemoteBlockReader) Close() error {
	r.locker.Lock()
	defer r.locker.Unlock()
	return r.closeFile()
}

type RemoteBlockService struct {
	NodeItemService   nodeitem.NodeItemInternalService
	RemoteBlockBroker *RemoteBlockBroker
}

func (s *RemoteBlockService) SelectForTopic(condition *RemoteBlockSelectCondition) (io.Reader, error) {
	nodeItem, err := s.NodeItemService.Select(uint(condition.ItemID))
	if err != nil {
		return nil, err
	}

	filePath := nodeItem.FilePath
	if nodeItem.FileType == nodeitem.FileTypeFolder {
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

	return &RemoteBlockReader{filePath: filePath, offset: condition.Offset, limit: condition.Limit}, nil
}

func (s *RemoteBlockService) SelectWithCondition(peerId peer.PeerID, condition *RemoteBlockSelectCondition) (io.ReadCloser, error) {
	return s.RemoteBlockBroker.Select(peerId, condition)
}
