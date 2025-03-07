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

// Next returns the next directory entry in the walk sequence. If there are no more entries,
// it returns nil. This method advances the internal cursor, so subsequent calls will return
// the next entry in the list. If the walk is completed, indicated by the Done method, it
// returns nil without advancing the cursor.

func (w *walkEntry) Next() fs.DirEntry {
	if w.Done() {
		return nil
	}
	entry := w.dirEntries[w.cursor]
	w.cursor++
	return entry
}

// Done checks if the walk sequence has been completed by determining if the
// internal cursor has reached or exceeded the number of directory entries.
// It returns true if there are no more entries to iterate over, otherwise false.

func (w *walkEntry) Done() bool {
	return w.cursor >= len(w.dirEntries)
}

// WalkRoot returns a sequence of strings representing the paths of all files
// and directories below the given root directory. The sequence is lazily
// generated and will stop iterating if the callback function returns false.
//
// The sequence first yields the root directory, then all of its subdirectories
// and files in depth-first order. If a directory is encountered, its path is
// yielded first, then all of its subdirectories and files are yielded, and so
// on.
//
// If an error occurs while walking the directory tree, the sequence will stop
// iterating and the error will be returned.
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
