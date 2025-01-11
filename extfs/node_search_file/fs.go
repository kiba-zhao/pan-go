package nodesearchfile

import (
	"io/fs"
	"iter"
	"os"
	"path"
)

type walkEntry struct {
	path       string
	isRoot     bool
	dirEntries []fs.DirEntry
	cursor     int
}

func (w *walkEntry) Next() fs.DirEntry {
	if w.Done() {
		return nil
	}
	entry := w.dirEntries[w.cursor]
	w.cursor++
	return entry
}

func (w *walkEntry) Done() bool {
	return w.cursor >= len(w.dirEntries)
}

func WalkRoot(root string) (iter.Seq[string], error) {

	stat, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	return func(yield func(string) bool) {

		rootDirEntry := fs.FileInfoToDirEntry(stat)
		var rootEntry walkEntry
		rootEntry.isRoot = true
		rootEntry.path = root
		rootEntry.cursor = 0
		rootEntry.dirEntries = []fs.DirEntry{rootDirEntry}

		entries := []*walkEntry{&rootEntry}
		cur := 0
		nextEntries := make([]*walkEntry, 0)

		for {
			if cur >= len(entries) {
				if len(nextEntries) <= 0 {
					return
				}
				entries = nextEntries
				nextEntries = make([]*walkEntry, 0)
				cur = 0
			}

			entry := entries[cur]
			dirEntry := entry.Next()
			if dirEntry == nil {
				return
			}
			if entry.Done() {
				cur++
			}

			var filePath string
			if entry.isRoot {
				filePath = entry.path
			} else {
				filePath = path.Join(entry.path, dirEntry.Name())
			}

			if !yield(filePath) {
				return
			}

			if dirEntry.IsDir() {
				subFiles, err := os.ReadDir(filePath)
				if err != nil {
					return
				}
				var subEntry walkEntry
				subEntry.path = filePath
				subEntry.dirEntries = subFiles
				subEntry.cursor = 0
				nextEntries = append(nextEntries, &subEntry)
			}
		}
	}, nil
}
