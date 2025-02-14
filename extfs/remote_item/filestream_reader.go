package remoteitem

import (
	"io"
	"pan/app/peer"
	"sync"
	"time"
)

type RemoteFileStreamReader struct {
	service *RemoteFileStreamService
	record  *RemoteFileInfoRecord
	peerId  peer.PeerID
	reader  io.ReadCloser
	rw      sync.RWMutex
	closed  bool
	offset  int64
}

func (r *RemoteFileStreamReader) Read(p []byte) (int, error) {
	r.rw.Lock()
	defer r.rw.Lock()

	if r.closed {
		return 0, io.EOF
	}

	if r.reader == nil {
		err := initRemoteFileStreamReaderWithOffset(r, r.offset)
		if err != nil {
			return 0, err
		}
	}

	n, err := r.reader.Read(p)
	if err == nil {
		r.offset += int64(n)
	}
	return n, err
}

func (r *RemoteFileStreamReader) Close() error {
	r.rw.Lock()
	defer r.rw.Unlock()
	if r.closed {
		return nil
	}
	r.closed = true
	if r.reader == nil {
		return nil
	}
	return r.reader.Close()
}

func (r *RemoteFileStreamReader) Seek(offset int64, whence int) (int64, error) {
	r.rw.Lock()
	defer r.rw.Unlock()
	if r.closed {
		return 0, io.ErrClosedPipe
	}
	if r.reader != nil {
		r.reader.Close()
	}

	var offset_ int64
	if whence == 0 {
		offset_ = offset
	} else if whence == 1 {
		offset_ = r.offset + offset
	} else {
		offset_ = r.record.Size + offset
	}

	err := initRemoteFileStreamReaderWithOffset(r, offset_)
	if err != nil {
		return 0, err
	}
	return r.offset, nil
}

func (r *RemoteFileStreamReader) Name() string {
	return r.record.Name
}

func (r *RemoteFileStreamReader) ModTime() time.Time {
	return time.Unix(r.record.UpdatedAt, 0)
}

func (r *RemoteFileStreamReader) Offset() int64 {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.offset
}

func initRemoteFileStreamReaderWithOffset(r *RemoteFileStreamReader, offset int64) error {
	var condition RemoteFileStreamSelectCondition
	condition.ItemID = r.record.ItemID
	condition.FilePath = r.record.FilePath
	condition.Offset = offset

	reader, err := r.service.SelectWithCondition(r.peerId, &condition)
	if err != nil {
		return err
	}
	r.reader = reader
	r.offset = offset
	return nil
}
