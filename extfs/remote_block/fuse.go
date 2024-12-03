package remoteblock

import (
	"context"
	"errors"
	"io"
	"pan/app/peer"
	"sync"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type FUSERemoteFile interface {
	PeerID() peer.PeerID
	ItemID() int32
	ParentPath() string
	Name() string
}

type RemoteBlockServiceFUSEProvider interface {
	RemoteBlockService() *RemoteBlockService
}

type FUSERemoteBlockReader struct {
	FUSERemoteFile FUSERemoteFile
	Provider       RemoteBlockServiceFUSEProvider
	locker         sync.Mutex
	reader         io.ReadCloser
	offset         int64
	eof            bool
}

func (fuserfr *FUSERemoteBlockReader) closeReader(eof bool) error {
	reader := fuserfr.reader
	var err error
	if reader != nil {
		err = reader.Close()
		fuserfr.reader = nil
	}

	fuserfr.eof = eof
	return err
}

func (fuserfr *FUSERemoteBlockReader) initReader(off int64) error {
	var condition RemoteBlockSelectCondition
	condition.ItemID = fuserfr.FUSERemoteFile.ItemID()
	condition.ParentPath = fuserfr.FUSERemoteFile.ParentPath()
	condition.Name = fuserfr.FUSERemoteFile.Name()
	condition.Offset = off

	remoteBlockService := fuserfr.Provider.RemoteBlockService()
	reader, err := remoteBlockService.SelectWithCondition(fuserfr.FUSERemoteFile.PeerID(), &condition)
	if err == nil {
		fuserfr.reader = reader
		fuserfr.eof = false
	}
	return err
}

var _ = (fs.FileReader)((*FUSERemoteBlockReader)(nil))

func (fuserfr *FUSERemoteBlockReader) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	fuserfr.locker.Lock()
	defer fuserfr.locker.Unlock()
	if fuserfr.offset < 0 {
		return nil, syscall.ENOENT
	}

	if fuserfr.eof && off == fuserfr.offset {
		return fuse.ReadResultData(nil), fs.OK
	}

	if fuserfr.reader != nil && off != fuserfr.offset {
		fuserfr.closeReader(false)
	}

	if fuserfr.reader == nil {
		err := fuserfr.initReader(off)
		if err != nil {
			return nil, syscall.ENOENT
		}
	}

	n, err := fuserfr.reader.Read(dest)
	if err != nil {
		fuserfr.closeReader(true)
		if !errors.Is(err, io.EOF) {
			return nil, syscall.ENOENT
		}
	}

	limit := int64(len(dest))
	var buffer []byte
	if int64(n) >= limit {
		fuserfr.offset = off + limit
		buffer = dest
	} else {
		fuserfr.offset = off + int64(n)
		buffer = dest[:n]
	}
	return fuse.ReadResultData(buffer), fs.OK
}

var _ = (fs.FileReleaser)((*FUSERemoteBlockReader)(nil))

func (fuserfr *FUSERemoteBlockReader) Release(ctx context.Context) syscall.Errno {
	fuserfr.locker.Lock()
	defer fuserfr.locker.Unlock()
	fuserfr.offset = -1
	err := fuserfr.closeReader(true)
	if err != nil {
		return syscall.ENOENT
	}
	return 0
}
