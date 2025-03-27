//go:build linux || (darwin && amd64)

package vfs

import (
	"bytes"
	"context"
	"errors"
	"os"
	appnode "pan/app/app_node"
	nodeitem "pan/extfs/node_item"
	remoteitem "pan/extfs/remote_item"
	remotenode "pan/extfs/remote_node"
	"slices"
	"syscall"
	"time"

	"sync"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var ErrVFSFUSEFSUnavailable = errors.New("vfs.VFSFUSEFS Error: FileSystem Unavailable")
var ErrVFSFUSEFSNameConflict = errors.New("vfs.VFSFUSEFS Error: Name Conflict")

type VFSFUSEFS struct {
	fs.Inode
	runtime  *VFSFSRuntime
	server   *fuse.Server
	settings *VFSSettings
	rw       sync.RWMutex
}

func NewVFSFS(runtime *VFSFSRuntime, settings VFSSettings) *VFSFUSEFS {
	var vfsfs VFSFUSEFS
	vfsfs.settings = &settings
	vfsfs.runtime = runtime
	return &vfsfs
}

// Mount mounts the FUSE file system.
//
// It takes a VFSSettings object as parameter and
// returns an error if the settings are invalid or if the mount fails.
// If the FUSE file system has already been mounted, it returns ErrFUSEUnavailable.
//
// The mount path is created if it does not exist.
func (fusefs *VFSFUSEFS) Mount() error {
	fusefs.rw.Lock()
	defer fusefs.rw.Unlock()
	if fusefs.server != nil {
		return ErrVFSFUSEFSUnavailable
	}
	settings := fusefs.settings

	opts := &fs.Options{}
	opts.DisableXAttrs = true
	opts.DisableReadDirPlus = true

	// auto create mount path
	mountPath := settings.MountPath
	if err := os.MkdirAll(mountPath, 0755); err != nil {
		return err
	}

	rawFS := fs.NewNodeFS(fusefs, opts)
	server, err := fuse.NewServer(rawFS, mountPath, &opts.MountOptions)
	if err != nil {
		return err
	}

	fusefs.server = server

	go server.Serve()
	err = server.WaitMount()
	return err
}

// Unmount unmounts the FUSE file system.
//
// It returns an error if the unmount fails. If the FUSE file system has already been unmounted,
// it returns nil.
func (fusefs *VFSFUSEFS) Unmount() error {

	fusefs.rw.Lock()
	defer fusefs.rw.Unlock()
	if fusefs.server == nil {
		return nil
	}

	err := fusefs.server.Unmount()
	if err == nil {
		fusefs.server.Wait()
		fusefs.server = nil
	}
	return err
}

// Getattr returns the attributes of the root node of the remote node.
// It always returns a directory with a size of 4096 bytes, and the
// current time as the last modified time. The number of hard links to
// the file is always 1.
// It implements the Getattr method of the fs.Inode interface.
func (fusefs *VFSFUSEFS) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = fuse.S_IFDIR
	out.Size = 4096
	now := time.Now()
	out.SetTimes(nil, &now, nil)
	out.Nlink = 1
	return fs.OK
}

// Readdir returns a directory stream containing the local node and all
// remote nodes as subdirectories. If any remote node has the same name
// as the local node, an error is returned. The returned directory stream
// is released by calling RmChild on the FUSERemoteNode instance with
// the names of the remote nodes that are not in the returned stream.
// It implements the Readdir method of the fs.Inode interface.
func (fusefs *VFSFUSEFS) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	dirs := make([]fuse.DirEntry, 0)
	names := make([]string, 0)

	fusefs.rw.RLock()
	settings := fusefs.settings
	runtime := fusefs.runtime
	var remoteNodeService *remotenode.RemoteNodeService
	if runtime != nil {
		remoteNodeService = runtime.RemoteNodeService()
	}
	defer fusefs.rw.RUnlock()

	if settings != nil {
		localEntry := fuse.DirEntry{Name: settings.LocalName, Mode: fuse.S_IFDIR}
		dirs = append(dirs, localEntry)
		names = append(names, localEntry.Name)
	}

	if remoteNodeService == nil {
		return fs.NewListDirStream(dirs), 0
	}

	_, remotes, err := remoteNodeService.SelectAll()
	if err != nil {
		return fs.NewListDirStream(dirs), 0
	}

	for _, remote := range remotes {
		idx, ok := slices.BinarySearch(names, remote.Name)
		if ok {
			err = ErrVFSFUSEFSNameConflict
			break
		}
		names = slices.Insert(names, idx, remote.Name)
		dirEntry := fuse.DirEntry{Name: remote.Name, Mode: fuse.S_IFDIR}
		dirs = append(dirs, dirEntry)
	}

	releaseNames := make([]string, 0)
	for n := range fusefs.Children() {
		if _, ok := slices.BinarySearch(names, n); ok {
			continue
		}
		releaseNames = append(releaseNames, n)
	}
	if len(releaseNames) > 0 {
		fusefs.RmChild(releaseNames...)
	}
	if err != nil {
		return nil, syscall.ENOENT
	}

	return fs.NewListDirStream(dirs), 0
}

// Lookup searches for a child inode with the given name in the current remote node.
// If the name matches the local node, it ensures the inode represents a local node
// and returns it. If the name corresponds to a remote node, it retrieves the
// corresponding remote node information and creates a new inode for it if necessary.
// Returns the found inode and fs.OK on success, or nil and ENOENT if the name
// does not correspond to any known node.
// It implements the Lookup method of the fs.Inode interface.
func (fusefs *VFSFUSEFS) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {

	fusefs.rw.RLock()
	settings := fusefs.settings
	runtime := fusefs.runtime
	var remoteNodeService *remotenode.RemoteNodeService
	if runtime != nil {
		remoteNodeService = runtime.RemoteNodeService()
	}
	defer fusefs.rw.RUnlock()

	inode := fusefs.GetChild(name)
	if settings != nil && name == settings.LocalName {
		if inode != nil {
			_, ok := inode.Operations().(*nodeitem.VFSFUSENodeItem)
			if ok {
				return inode, fs.OK
			}
			fusefs.RmChild(name)
		}
		node := nodeitem.NewVFSFUSENodeItem(runtime)
		inode = fusefs.NewInode(ctx, node, fs.StableAttr{Mode: fuse.S_IFDIR})
		return inode, fs.OK
	}

	remote, err := remoteNodeService.SelectByName(name)
	if err != nil {
		return nil, syscall.ENOENT
	}

	peerId, err := appnode.DecodePeerID(remote.PeerID)
	if err != nil {
		return nil, syscall.ENOENT
	}

	if inode != nil {
		fuseRemoteItemList, ok := inode.Operations().(*remoteitem.VFSFUSERemoteItemList)
		if ok && bytes.Equal(fuseRemoteItemList.PeerID(), peerId) {
			return inode, fs.OK
		}
		fusefs.RmChild(name)
	}

	node := remoteitem.NewVFSFUSERemoteItem(peerId, runtime)
	inode = fusefs.NewInode(ctx, node, fs.StableAttr{Mode: fuse.S_IFDIR})
	return inode, fs.OK
}
