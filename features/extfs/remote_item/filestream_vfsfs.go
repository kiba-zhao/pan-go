//go:build linux || (darwin && amd64)

package remoteitem

import (
	"context"
	"errors"
	"io"
	"sync"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type VFSFUSERemoteFileStreamReader struct {
	reader *RemoteFileStreamReader
	locker sync.Mutex
}

var _ = (fs.FileReader)((*VFSFUSERemoteFileStreamReader)(nil))

func (fuserfr *VFSFUSERemoteFileStreamReader) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {

	fuserfr.locker.Lock()
	defer fuserfr.locker.Unlock()

	var err error
	if fuserfr.reader.Offset() != off {
		_, err = fuserfr.reader.Seek(off, 0)
	}

	if err != nil {
		return nil, syscall.ENOENT
	}

	n, err := fuserfr.reader.Read(dest)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, syscall.ENOENT
	}

	if n == 0 {
		return fuse.ReadResultData(nil), fs.OK
	}

	if n < len(dest) {
		dataBytes := dest[:n]
		return fuse.ReadResultData(dataBytes), fs.OK
	}

	return fuse.ReadResultData(dest), fs.OK
}

var _ = (fs.FileReleaser)((*VFSFUSERemoteFileStreamReader)(nil))

func (fuserfr *VFSFUSERemoteFileStreamReader) Release(ctx context.Context) syscall.Errno {
	fuserfr.locker.Lock()
	defer fuserfr.locker.Unlock()

	err := fuserfr.reader.Close()
	if err != nil {
		return syscall.ENOENT
	}
	return 0
}
