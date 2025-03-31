package remoteitem

import (
	"math"
	nodeitem "pan/features/extfs/node_item"
	"sync"

	"github.com/winfsp/cgofuse/fuse"
)

type VFSWinFSRuntime interface {
	VFSFUSERemoteItemRuntime
	Parse(path string) (name string, subPath string)
}

func NewVFSWinFS(runtime VFSWinFSRuntime) *VFSWinFS {
	var winfs VFSWinFS
	winfs.runtime = runtime
	winfs.fhMap = make(map[uint64]*RemoteFileStreamReader)
	winfs.fh = 0
	return &winfs
}

type VFSWinFS struct {
	fuse.FileSystemBase
	runtime VFSWinFSRuntime
	fhMap   map[uint64]*RemoteFileStreamReader
	fhRW    sync.RWMutex
	fh      uint64
}

func (winfs *VFSWinFS) SetStatWithRemoteFile(stat *fuse.Stat_t, remoteStat RemoteFileStat) {
	if remoteStat.FileType == nodeitem.FileTypeFolder {
		stat.Mode = fuse.S_IFDIR | 0555
	} else {
		stat.Mode = fuse.S_IFREG | 0444
	}

	stat.Size = remoteStat.Size
	stat.Ctim = fuse.Timespec{Sec: 0, Nsec: remoteStat.CreatedAt.UnixNano()}
	stat.Mtim = fuse.Timespec{Sec: 0, Nsec: remoteStat.UpdatedAt.UnixNano()}
	stat.Atim = stat.Mtim
}

func (winfs *VFSWinFS) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	remoteNodeService := winfs.runtime.RemoteNodeService()
	remoteItemService := winfs.runtime.RemoteItemService()
	remoteFileInfoService := winfs.runtime.RemoteFileInfoService()
	if remoteItemService == nil || remoteNodeService == nil || remoteFileInfoService == nil {
		return fuse.EFAULT
	}

	nodeName, nodePath := winfs.runtime.Parse(path)
	remoteNode, err := remoteNodeService.SelectByName(nodeName)
	if err != nil {
		return fuse.EFAULT
	}

	// list all remote items
	if len(nodePath) <= 1 {

		_, remoteItems, err := remoteItemService.Search(remoteNode.PeerID)
		if err != nil {
			return fuse.EFAULT
		}
		for _, remoteItem := range remoteItems {
			var stat fuse.Stat_t
			winfs.SetStatWithRemoteFile(&stat, remoteItem.RemoteFileStat)
			fill(remoteItem.Name, &stat, 0)
		}
		return 0
	}
	//

	name, subPath := winfs.runtime.Parse(nodePath)
	remoteItem, err := remoteItemService.SelectByName(remoteNode.PeerID, name)
	if remoteItemService.IsNotExist(err) {
		return fuse.EBADF
	}
	if err != nil {
		return fuse.EFAULT
	}
	_, remoteFileInfos, err := remoteFileInfoService.Search(remoteNode.PeerID, remoteItem.ItemID, subPath[1:])
	if err != nil {
		return fuse.EFAULT
	}

	for _, remoteFileInfo := range remoteFileInfos {
		var stat fuse.Stat_t
		winfs.SetStatWithRemoteFile(&stat, remoteFileInfo.RemoteFileStat)
		fill(remoteFileInfo.Name, &stat, 0)
	}

	return 0
}

func (winfs *VFSWinFS) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	remoteNodeService := winfs.runtime.RemoteNodeService()
	remoteItemService := winfs.runtime.RemoteItemService()
	remoteFileInfoService := winfs.runtime.RemoteFileInfoService()
	if remoteItemService == nil || remoteNodeService == nil || remoteFileInfoService == nil {
		return fuse.EFAULT
	}

	// Select remote node with remote name
	nodeName, nodePath := winfs.runtime.Parse(path)
	remoteNode, err := remoteNodeService.SelectByName(nodeName)
	if err != nil {
		return fuse.EFAULT
	}
	//

	// Select remote item with remote item name
	name, subPath := winfs.runtime.Parse(nodePath)
	remoteItem, err := remoteItemService.SelectByName(remoteNode.PeerID, name)
	if remoteItemService.IsNotExist(err) {
		return fuse.EBADF
	}
	if err != nil {
		return fuse.EFAULT
	}

	if len(subPath) <= 1 {
		winfs.SetStatWithRemoteFile(stat, remoteItem.RemoteFileStat)
		return 0
	}
	//

	// Select remote file with remote file path
	remoteFileInfo, err := remoteFileInfoService.Select(remoteNode.PeerID, remoteItem.ItemID, subPath[1:])
	if remoteFileInfoService.IsNotExist(err) {
		return fuse.ENFILE
	}
	if err != nil {
		return fuse.EFAULT
	}

	winfs.SetStatWithRemoteFile(stat, remoteFileInfo.RemoteFileStat)
	return 0
}

func (winfs *VFSWinFS) Open(path string, flags int) (errc int, fh uint64) {
	remoteNodeService := winfs.runtime.RemoteNodeService()
	remoteItemService := winfs.runtime.RemoteItemService()
	remoteFileStreamService := winfs.runtime.RemoteFileStreamService()
	if remoteItemService == nil || remoteNodeService == nil || remoteFileStreamService == nil {
		return fuse.EFAULT, 0
	}

	nodeName, nodePath := winfs.runtime.Parse(path)
	remoteNode, err := remoteNodeService.SelectByName(nodeName)
	if err != nil {
		return fuse.ENFILE, 0
	}

	name, subPath := winfs.runtime.Parse(nodePath)
	remoteItem, err := remoteItemService.SelectByName(remoteNode.PeerID, name)
	if remoteItemService.IsNotExist(err) {
		return fuse.ENFILE, 0
	}
	if err != nil {
		return fuse.EFAULT, 0
	}

	reader, err := remoteFileStreamService.Read(remoteNode.PeerID, remoteItem.ItemID, subPath[1:])
	if remoteFileStreamService.IsNotExist(err) {
		return fuse.ENFILE, 0
	}
	if err != nil {
		return fuse.EFAULT, 0
	}

	winfs.fhRW.Lock()
	defer winfs.fhRW.Unlock()

	count := 0
	for ; count < math.MaxUint8; count++ {
		winfs.fh++
		if _, ok := winfs.fhMap[winfs.fh]; !ok {
			winfs.fhMap[winfs.fh] = reader
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
	reader, ok := winfs.fhMap[fh]
	winfs.fhRW.RUnlock()
	if !ok {
		return 0
	}

	// TODO: log error
	_, err := reader.Seek(ofst, 0)
	if err != nil {
		return 0
	}

	n, _ := reader.Read(buff)
	return n
}

func (winfs *VFSWinFS) Release(path string, fh uint64) int {
	ret := 0
	winfs.fhRW.Lock()
	defer winfs.fhRW.Unlock()

	reader, ok := winfs.fhMap[fh]
	if ok {
		delete(winfs.fhMap, fh)
		// TODO: log error
		reader.Close()
	} else {
		ret = int(fuse.EFAULT)
	}

	return ret
}
