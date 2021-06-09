// Package dirtree provides objects and methods to work with a filesystem
// hierarchy.
package dirtree

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

// A logger that outputs to stderr without the timestamp.
var errLog = log.New(os.Stderr, "", 0)

// A simple structs that represents a file in a directory
type FileInfo struct {
	path string
	size int64
}

// A directory tree with accumulated sizes for each directory
type DirTree struct {
	// Root of the tree
	path string
	// Cumulative size of the tree
	size int64
	// List of files on the root level
	files []FileInfo
	// List of sub-directories
	subdirs []*DirTree
	// Unit size used for displaying. Posix standard says it should be 512 bytes
	// but most modern implementations (such as GNU coreutils) use 1024. We will
	// stick with Posix. The `-k` flag allows to switch to 1024 instead.
	unitSize int64
	// Filesystem block size. The default is 4k but we will try to get the real
	// size for each filesystem later.
	blockSize int64
}

// New creates a new directory tree rooted at `path`.
func New(path string, unitSize int64) *DirTree {
	dt := new(DirTree)
	dt.path = path
	dt.unitSize = unitSize
	bs, err := getFSBlockSize(path)
	if err != nil {
		dt.blockSize = 4096
	} else {
		dt.blockSize = bs
	}
	dt.buildDirTree()

	return dt
}

// buildDirTree builds a hierarchy of directories represented by a dirTree
// structure starting from the directory given by `dt`.
//
// Any errors encountered during the traversal will be printed to stderr and
// will not cause the function to fail.
func (dt *DirTree) buildDirTree() {
	dtInfo, err := os.Stat(dt.path)
	if err != nil {
		errLog.Println(err)
	}
	dt.size = dt.calcSize(dtInfo.Size())
	if !dtInfo.IsDir() {
		return
	}

	files, err := os.ReadDir(dt.path)
	if err != nil {
		errLog.Println(err)
	}
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			errLog.Println(err)
			continue
		}
		if f.IsDir() {
			sdt := New(filepath.Join(dt.path, f.Name()), dt.unitSize)
			dt.size = dt.size + sdt.size
			dt.subdirs = append(dt.subdirs, sdt)
		} else {
			dt.size = dt.size + dt.calcSize(info.Size())
			fi := FileInfo{
				path: filepath.Join(dt.path, info.Name()),
				size: dt.calcSize(info.Size()),
			}
			dt.files = append(dt.files, fi)
		}
	}
}

// PrintDirTree walks over `dt` recursively and returns a slice of strings.

// Each line in the slice is either a file or a directory and it's size
// formatted according to `outFormat` string.
// If `countFiles` is true files are printed first. If `summarise` is true
// sub-directories are not printed out.
func (dt *DirTree) PrintDirTree(outFormat string, countFiles bool, summarise bool) []string {
	var out []string
	// If "-a" is provided output files first
	if countFiles {
		for _, f := range dt.files {
			out = append(out, fmt.Sprintf(outFormat, f.size, fixPath(f.path)))
		}
	}
	if !summarise {
		for _, d := range dt.subdirs {
			out = append(out, d.PrintDirTree(outFormat, countFiles, summarise)...)
		}
	}
	out = append(out, fmt.Sprintf(outFormat, dt.size, fixPath(filepath.Clean(dt.path))))

	return out
}

// calcSize receives size in bytes and returns size in units.
//
// Filesystem allocates space in blocks and not in bytes. That is why the
// actual size of the file is usually smaller than the space allocated for
// it by the system. Since we want to report the acual space in use and
// not the file size we need calculate the number of filesystem blocks
// allocated to the file.
func (dt *DirTree) calcSize(size int64) int64 {
	allocSize := (1 + (size-1)/dt.blockSize) * dt.blockSize
	return 1 + (allocSize-1)/dt.unitSize
}

// Get filesystem block size.
//
// Returns block size in bytes or error.
// https://man7.org/linux/man-pages/man2/statfs.2.html
func getFSBlockSize(path string) (int64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	} else {
		return stat.Bsize, nil
	}
}

// fixPath adds './' to the begining of the relative paths.
func fixPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if path[0] == '.' {
		return path
	}
	return "./" + path
}
