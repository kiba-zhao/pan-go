//go:build windows

package vfs

import (
	"errors"
	nodeitem "pan/features/extfs/node_item"
	remoteitem "pan/features/extfs/remote_item"
	remotenode "pan/features/extfs/remote_node"
	"strings"
	"sync"

	"github.com/winfsp/cgofuse/fuse"
)

var ErrVFSWinFSMountFailed = errors.New("vfs.VFSWinFS Error: Mount Failed")

const VFSWinFSPathSeparator = "/"

type VFSWinFSNodeSystem = fuse.FileSystemInterface

func NewVFSFS(server *stdVFSServer) *VFSWinFS {
	var vfsfs VFSWinFS
	vfsfs.stdVFSFSRuntime = server.runtime
	vfsfs.vfsServer = server
	vfsfs.nodeItemSys = nodeitem.NewVFSWinFS(&vfsfs)
	vfsfs.remoteItemSys = remoteitem.NewVFSWinFS(&vfsfs)
	return &vfsfs
}

type VFSWinFS struct {
	fuse.FileSystemBase
	*stdVFSFSRuntime
	vfsServer     *stdVFSServer
	locker        sync.Mutex
	host          *fuse.FileSystemHost
	nodeItemSys   VFSWinFSNodeSystem
	remoteItemSys VFSWinFSNodeSystem
}

var _ = (VFSFS)((*VFSWinFS)(nil))

func (winfs *VFSWinFS) Mount() error {
	winfs.locker.Lock()
	defer winfs.locker.Unlock()

	vfsServer := winfs.vfsServer
	mountPath := vfsServer.MountPath()
	if len(mountPath) <= 0 {
		return ErrVFSWinFSMountFailed
	}

	opts := make([]string, 0)
	opts = append(opts, "-o", "uid=-1")
	opts = append(opts, "-o", "gid=-1")
	opts = append(opts, "--FileSystemName=extfs")

	winfs.host = fuse.NewFileSystemHost(winfs)
	ok := winfs.host.Mount(mountPath, opts)
	if !ok {
		return ErrVFSWinFSMountFailed
	}
	return nil
}

func (winfs *VFSWinFS) Unmount() error {
	winfs.locker.Lock()
	defer winfs.locker.Unlock()
	if winfs.host != nil {
		winfs.host.Unmount()
	}
	winfs.host = nil
	return nil
}

func (winfs *VFSWinFS) Parse(path string) (name string, subPath string) {
	if len(path) <= 1 {
		return "", VFSWinFSPathSeparator
	}

	idx := strings.Index(path[1:], VFSWinFSPathSeparator)
	if idx < 0 {
		return path[1:], VFSWinFSPathSeparator
	}
	return path[1 : idx+1], path[idx+1:]
}

func (winfs *VFSWinFS) lookupNodeSys(name string) (VFSWinFSNodeSystem, bool) {
	if len(name) <= 0 {
		return nil, false
	}

	vfsServer := winfs.vfsServer
	hostname := vfsServer.HostName()
	if len(hostname) > 0 && name == hostname {
		return winfs.nodeItemSys, false
	}
	return winfs.remoteItemSys, true
}

func (winfs *VFSWinFS) lookupNodeStat(name string) *fuse.Stat_t {

	var stat fuse.Stat_t
	stat.Mode = fuse.S_IFDIR | 0555
	return &stat
}

func (winfs *VFSWinFS) Init() {
	winfs.nodeItemSys.Init()
	winfs.remoteItemSys.Init()
}

func (winfs *VFSWinFS) Destroy() {
	winfs.nodeItemSys.Destroy()
	winfs.remoteItemSys.Destroy()
}

func (winfs *VFSWinFS) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {

	name, subPath := winfs.Parse(path)
	if len(name) > 0 {
		subSys, fullPass := winfs.lookupNodeSys(name)
		if fullPass {
			return subSys.Readdir(path, fill, ofst, fh)
		}
		return subSys.Readdir(subPath, fill, ofst, fh)
	}

	// fill(".", winfs.lookupNodeStat("."), 0)
	// fill("..", winfs.lookupNodeStat(".."), 0)
	vfsServer := winfs.vfsServer
	hostname := vfsServer.HostName()
	if len(hostname) > 0 {
		fill(hostname, winfs.lookupNodeStat(hostname), 0)
	}

	var remoteNodeArr []remotenode.RemoteNode
	var err error
	if remoteNodeService := winfs.RemoteNodeService(); remoteNodeService != nil {
		_, remoteNodeArr, err = remoteNodeService.SelectAll()
	}

	if err != nil || len(remoteNodeArr) <= 0 {
		return 0
	}

	for _, remoteNode := range remoteNodeArr {
		fill(remoteNode.Name, winfs.lookupNodeStat(remoteNode.Name), 0)
	}
	return 0
}

func (winfs *VFSWinFS) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	name, subPath := winfs.Parse(path)
	if len(name) > 0 && len(subPath) > 1 {
		subSys, fullPass := winfs.lookupNodeSys(name)
		if fullPass {
			return subSys.Getattr(path, stat, fh)
		}
		return subSys.Getattr(subPath, stat, fh)
	}

	if len(name) > 0 {
		vfsServer := winfs.vfsServer
		hostname := vfsServer.HostName()
		if len(hostname) > 0 && hostname != name {
			remoteNodeService := winfs.RemoteNodeService()
			if remoteNodeService == nil {
				return fuse.ENOATTR
			}
			_, err := remoteNodeService.SelectByName(name)
			if err != nil {
				return fuse.ENOATTR
			}
		}
	}

	stat.Mode = fuse.S_IFDIR | 0555
	return 0
}

func (winfs *VFSWinFS) Open(path string, flags int) (errc int, fh uint64) {
	name, subPath := winfs.Parse(path)
	if len(name) > 0 && len(subPath) > 1 {
		subSys, fullPass := winfs.lookupNodeSys(name)
		if fullPass {
			return subSys.Open(path, flags)
		}
		return subSys.Open(subPath, flags)
	}

	return fuse.ENFILE, 0
}

func (winfs *VFSWinFS) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	name, subPath := winfs.Parse(path)
	if len(name) > 0 && len(subPath) > 1 {
		subSys, fullPass := winfs.lookupNodeSys(name)
		if fullPass {
			return subSys.Read(path, buff, ofst, fh)
		}
		return subSys.Read(subPath, buff, ofst, fh)
	}
	return fuse.ENFILE
}

func (winfs *VFSWinFS) Release(path string, fh uint64) int {
	name, subPath := winfs.Parse(path)
	if len(name) > 0 && len(subPath) > 1 {
		subSys, fullPass := winfs.lookupNodeSys(name)
		if fullPass {
			return subSys.Release(path, fh)
		}
		return subSys.Release(subPath, fh)
	}
	return fuse.ENFILE
}
