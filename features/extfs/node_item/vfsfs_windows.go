//go:build windows

package nodeitem

import (
	"math"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/winfsp/cgofuse/fuse"
)

type VFSWinFSRuntime interface {
	VFSFUSENodeItemRuntime
	Parse(path string) (name string, subPath string)
}

func NewVFSWinFS(runtime VFSWinFSRuntime) *VFSWinFS {
	var winfs VFSWinFS
	winfs.runtime = runtime
	winfs.fhMap = make(map[uint64]*os.File)
	winfs.fh = 0
	return &winfs
}

type VFSWinFS struct {
	fuse.FileSystemBase
	runtime VFSWinFSRuntime
	fhMap   map[uint64]*os.File
	fhRW    sync.RWMutex
	fh      uint64
}

func (winfs *VFSWinFS) SetStatWithFileInfo(stat *fuse.Stat_t, info os.FileInfo) {
	if info.IsDir() {
		stat.Mode = fuse.S_IFDIR | 0555
	} else {
		stat.Mode = fuse.S_IFREG | 0444
	}

	stat.Size = int64(info.Size())

	winStat := info.Sys().(*syscall.Win32FileAttributeData)
	stat.Atim = fuse.Timespec{Sec: 0, Nsec: winStat.LastAccessTime.Nanoseconds()}
	stat.Mtim = fuse.Timespec{Sec: 0, Nsec: winStat.LastWriteTime.Nanoseconds()}
	stat.Ctim = fuse.Timespec{Sec: 0, Nsec: winStat.CreationTime.Nanoseconds()}

}

func (winfs *VFSWinFS) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {

	nodeItemService := winfs.runtime.NodeItemService()
	if nodeItemService == nil {
		return fuse.EFAULT
	}

	// list all node items
	if len(path) <= 1 {
		nodeItems, err := nodeItemService.SelectAllWithEnabled(true)
		if err != nil {
			return fuse.EFAULT
		}
		for _, nodeItem := range nodeItems {
			info, err := os.Stat(nodeItem.FilePath)
			if err != nil {
				continue
			}
			var stat fuse.Stat_t
			winfs.SetStatWithFileInfo(&stat, info)
			fill(nodeItem.Name, &stat, 0)
		}
		return 0
	}
	//

	name, subPath := winfs.runtime.Parse(path)
	nodeItem, err := nodeItemService.SelectByName(name)
	if nodeItemService.IsNotExist(err) {
		return fuse.ENOTDIR
	}
	if err != nil {
		return fuse.EFAULT
	}
	fullPath := nodeItem.FilePath
	if len(subPath) >= 1 {
		fullPath = filepath.Join(fullPath, filepath.FromSlash(subPath))
	}

	dirEntries, err := os.ReadDir(fullPath)
	if err != nil {
		return fuse.ENOTDIR
	}

	if len(dirEntries) <= 0 {
		return 0
	}

	for _, dirEntry := range dirEntries {
		info, err := dirEntry.Info()
		if err != nil {
			continue
		}
		var stat fuse.Stat_t
		winfs.SetStatWithFileInfo(&stat, info)
		fill(dirEntry.Name(), &stat, 0)
	}

	return 0
}

func (winfs *VFSWinFS) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	nodeItemService := winfs.runtime.NodeItemService()
	if nodeItemService == nil {
		return fuse.EFAULT
	}

	name, subPath := winfs.runtime.Parse(path)
	nodeItem, err := nodeItemService.SelectByName(name)
	if nodeItemService.IsNotExist(err) {
		return fuse.ENOATTR
	}
	if err != nil {
		return fuse.EFAULT
	}
	fullPath := nodeItem.FilePath
	if len(subPath) >= 1 {
		fullPath = filepath.Join(fullPath, filepath.FromSlash(subPath))
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return fuse.ENOATTR
	}
	winfs.SetStatWithFileInfo(stat, info)
	return 0
}

func (winfs *VFSWinFS) Open(path string, flags int) (errc int, fh uint64) {
	nodeItemService := winfs.runtime.NodeItemService()
	if nodeItemService == nil {
		return fuse.EFAULT, 0
	}

	name, subPath := winfs.runtime.Parse(path)
	nodeItem, err := nodeItemService.SelectByName(name)
	if nodeItemService.IsNotExist(err) {
		return fuse.ENFILE, 0
	}
	if err != nil {
		return fuse.EFAULT, 0
	}
	fullPath := nodeItem.FilePath
	if len(subPath) >= 1 {
		fullPath = filepath.Join(fullPath, filepath.FromSlash(subPath))
	}

	winfs.fhRW.Lock()
	defer winfs.fhRW.Unlock()

	file, err := os.OpenFile(fullPath, flags, 0)
	if err != nil {
		return fuse.EFAULT, 0
	}

	count := 0
	for ; count < math.MaxUint8; count++ {
		winfs.fh++
		if _, ok := winfs.fhMap[winfs.fh]; !ok {
			winfs.fhMap[winfs.fh] = file
			break
		}
	}

	if count >= math.MaxUint8 {
		return fuse.ENFILE, 0
	}

	return 0, winfs.fh
}

func (winfs *VFSWinFS) Read(path string, buff []byte, ofst int64, fh uint64) int {
	winfs.fhRW.RLock()
	file, ok := winfs.fhMap[fh]
	winfs.fhRW.RUnlock()
	if !ok {
		return 0
	}

	// TODO: log error
	n, _ := file.ReadAt(buff, ofst)
	return n
}

func (winfs *VFSWinFS) Release(path string, fh uint64) int {
	ret := 0
	winfs.fhRW.Lock()
	defer winfs.fhRW.Unlock()

	file, ok := winfs.fhMap[fh]
	if ok {
		delete(winfs.fhMap, fh)
		// TODO: log error
		file.Close()
	} else {
		ret = int(fuse.EFAULT)
	}

	return ret
}
