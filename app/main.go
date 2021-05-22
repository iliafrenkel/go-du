// Copyright 2021 Ilia Frenkel. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

// By default, the du utility shall write to standard output the
// size of the file space allocated to, and the size of the file
// space allocated to each subdirectory of, the file hierarchy
// rooted in each of the specified files. By default, when a
// symbolic link is encountered on the command line or in the file
// hierarchy, du shall count the size of the symbolic link (rather
// than the file referenced by the link), and shall not follow the
// link to another portion of the file hierarchy. The size of the
// file space allocated to a file of type directory shall be defined
// as the sum total of space allocated to all files in the file
// hierarchy rooted in the directory plus the space allocated to the
// directory itself.

// When du cannot stat() files or stat() or read directories, it
// shall report an error condition and the final exit status is
// affected. A file that occurs multiple times under one file
// operand and that has a link count greater than 1 shall be counted
// and written for only one entry. It is implementation-defined
// whether a file that has a link count no greater than 1 is counted
// and written just once, or is counted and written for each
// occurrence. It is implementation-defined whether a file that
// occurs under one file operand is counted for other file operands.
// The directory entry that is selected in the report is
// unspecified. By default, file sizes shall be written in 512-byte
// units, rounded up to the next 512-byte unit.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Format for printing out dir/file entry
const outFormat = "%d\t%s"

// Command-line flags
type options struct {
	BlockSize       bool `short:"k" default:"false" description:"Write the files sizes in units of 1024 bytes, rather than the default 512-byte units"`
	CountFiles      bool `short:"a" long:"all" default:"false" description:"write counts for all files, not just directories"`
	DereferenceAll  bool `short:"L" long:"dereference" default:"false" description:"dereference all symbolic links"`
	DereferenceArgs bool `short:"H" long:"dereference-args" default:"false" description:"dereference only symlinks that are listed on the command line"`
	OneFileSystem   bool `short:"x" long:"one-files-system" default:"false" description:"skip directories on different file systems"`
	Summarise       bool `short:"s" long:"summarise" default:"false" description:"display only a total for each argument"`
}

var opts options

// Holds the file/directory names from the command line arguments.
var argFiles []string

// Filesystem block size. The default is 4k but we will try to get the real
// size for each filesystem later.
var fsBlockSize int64 = 4096

// Unit size used for displaying. Posix standard says it should be 512 bytes
// but most modern implementations (such as GNU coreutils) use 1024. We will
// stick with Posix. The `-k` flag allows to switch to 1024 instead.
var unitSize int64 = 512

// A logger that outputs to stderr without the timestamp.
var errLog = log.New(os.Stderr, "", 0)

// A directory tree with accumulated sizes for each directory
type dirTree struct {
	path    string
	size    int64
	files   []fileInfo
	subdirs []dirTree
}

// A simple structs that represents a file in a directory
type fileInfo struct {
	path string
	size int64
}

// calcSize receives size in bytes and returns size in units.
//
// Filesystem allocates space in blocks and not in bytes. That is why the
// actual size of the file is usually smaller than the space allocated for
// it by the system. Since we want to report the acual space in use and
// not the file size we need calculate the number of filesystem blocks
// allocated to the file.
func calcSize(size int64) int64 {
	// Set the unit size to 1024 if "-k" is specified
	if opts.BlockSize {
		unitSize = 1024
	} else {
		unitSize = 512
	}
	allocSize := (1 + (size-1)/fsBlockSize) * fsBlockSize
	return 1 + (allocSize-1)/unitSize
}

// buildDirTree builds a hierarchy of directories represented by a dirTree
// structure starting from the directory given by `dt`.
//
// Any errors encountered during the traversal will be printed to stderr and
// will not cause the function to fail.
func buildDirTree(dt *dirTree) {
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
			sdt := dirTree{
				path:    filepath.Join(dt.path, f.Name()),
				size:    calcSize(fsBlockSize),
				files:   []fileInfo{},
				subdirs: []dirTree{},
			}
			buildDirTree(&sdt)
			dt.size = dt.size + sdt.size
			dt.subdirs = append(dt.subdirs, sdt)
		} else {
			dt.size = dt.size + calcSize(info.Size())
			fi := fileInfo{
				path: filepath.Join(dt.path, info.Name()),
				size: calcSize(info.Size()),
			}
			dt.files = append(dt.files, fi)
		}
	}
}

func fixPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if path[0] == '.' {
		return path
	}
	return "./" + path
}

// printDirTree walks over `dt` recursively and returns a slice of strings that
// represents the output of the `go-du` command taking into account various
// command line flags.
//
// Files are printed first but only if `-a` flag was specified.
// Format is defined by the `outFormat` constant.
//
// We return a slice of strings instead of printing it out directly to allow
// testing of the output.
func printDirTree(dt dirTree) []string {
	var out []string
	// If "-a" is provided output files first
	if opts.CountFiles {
		for _, f := range dt.files {
			out = append(out, fmt.Sprintf(outFormat, f.size, fixPath(f.path)))
		}
	}
	if !opts.Summarise {
		for _, d := range dt.subdirs {
			out = append(out, printDirTree(d)...)
		}
	}
	out = append(out, fmt.Sprintf(outFormat, dt.size, fixPath(filepath.Clean(dt.path))))

	return out
}

// Declare and parse command line flags.
func init() {
	// Define command-line flags
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-a|-s] [-kx] [-H|-L] [FILE...]\n", os.Args[0])
		fmt.Printf("Summarise disk usage of the set of FILEs, recursively for directories.\n\n")
		flag.PrintDefaults()
		fmt.Println("\nThis is POSIX compatible implementation of the du utility. For exended")
		fmt.Println("documentation see https://man7.org/linux/man-pages/man1/du.1p.html")
		fmt.Println("\nDisplay values are in 512-byte units, rounded up to the next 512-byte unit")
		fmt.Println("unless -k flag is specified.")
		fmt.Println("\nCreated by Ilia Frenkel<frenkel.ilia@gmail.com>")
		fmt.Println("Report bugs at https://github.com/iliafrenkel/go-du")
	}
	flag.BoolVar(&opts.BlockSize, "k", false, "\tWrite the files sizes in units of 1024 bytes, rather than the\n\tdefault 512-byte units")
	flag.BoolVar(&opts.CountFiles, "a", false, "\twrite counts for all files, not just directories")
	flag.BoolVar(&opts.DereferenceAll, "L", false, "\tdereference all symbolic links")
	flag.BoolVar(&opts.DereferenceArgs, "H", false, "\tdereference only symlinks that are listed on the command line")
	flag.BoolVar(&opts.OneFileSystem, "x", false, "\tskip directories on different file systems")
	flag.BoolVar(&opts.Summarise, "s", false, "\tdisplay only a total for each argument")
	flag.Parse()
	// Check that there are no conflicts between flags
	if opts.CountFiles && opts.Summarise {
		errLog.Fatal("Cannot both summarise and show all entries.")
	}
	fmt.Printf("os.Args: %v\n", os.Args)
}

func main() {
	// If there are no arguments provided, default to the current directory
	argFiles = flag.Args()
	if len(argFiles) == 0 {
		argFiles = append(argFiles, ".")
	}

	for _, file := range argFiles {
		f, err := os.Stat(file)
		if err != nil {
			errLog.Println(err)
			continue
		}
		if f.Mode().IsRegular() { // it's a file, print out its size
			fmt.Printf("%v\t%s\n", calcSize(f.Size()), f.Name())
		} else { // it's a dir, count all the file sizes
			dt := dirTree{
				path:    file,
				size:    0,
				files:   []fileInfo{},
				subdirs: []dirTree{},
			}
			buildDirTree(&dt)
			for _, s := range printDirTree(dt) {
				fmt.Println(s)
			}
		}
	}
}
