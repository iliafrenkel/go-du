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
	"syscall"
)

// Command-line flags
var opts struct {
	BlockSize       bool `short:"k" default:"false" description:"Write the files sizes in units of 1024 bytes, rather than the default 512-byte units"`
	CountFiles      bool `short:"a" long:"all" default:"false" description:"write counts for all files, not just directories"`
	DereferenceAll  bool `short:"L" long:"dereference" default:"false" description:"dereference all symbolic links"`
	DereferenceArgs bool `short:"H" long:"dereference-args" default:"false" description:"dereference only symlinks that are listed on the command line"`
	OneFileSystem   bool `short:"x" long:"one-files-system" default:"false" description:"skip directories on different file systems"`
	Summarise       bool `short:"s" long:"summarise" default:"false" description:"display only a total for each argument"`
}

// Holds the file/directory names from the command line arguments
var argFiles []string

// Filesystem block size
var fsBlockSize int64 = 4096

// Unit size
var unitSize int64 = 512

// A logger that outputs to stderr without the timestamp
var errLog = log.New(os.Stderr, "", 0)

// Receives size in bytes and returns size in units.
//
// Filesystem allocates space in blocks and not in bytes. That is why the
// actual size of the file is usually smaller than the space allocated for
// it by the system. Since we want to report the acual space in use and
// not the file size we need calculate the number of filesystem blocks
// allocated to the file.
func calcSize(size int64) int64 {
	allocSize := (1 + (size-1)/fsBlockSize) * fsBlockSize
	return 1 + (allocSize-1)/unitSize
}

// Prints out total size of all files in a directory recursively.
// Any errors encountered during traversal will be printed to stderr and will
// not cause the function to fail.
func dirSize(path string) int64 {
	files, err := os.ReadDir(path)
	if err != nil {
		errLog.Println(err)
	}
	var size int64
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			errLog.Println(err)
			continue
		}
		if f.IsDir() {
			size = size + dirSize(filepath.Join(path, f.Name()))
		} else {
			size = size + info.Size()
			if opts.CountFiles {
				fmt.Printf("%v\t%s\n", calcSize(info.Size()), path+string(filepath.Separator)+f.Name())
			}
		}
	}

	fmt.Printf("%v\t%s\n", calcSize(size), path)
	return size
}

func main() {
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

	// If there are no arguments provided, default to the current directory
	argFiles = flag.Args()
	if len(argFiles) == 0 {
		argFiles = append(argFiles, ".")
	}

	// https://man7.org/linux/man-pages/man2/statfs.2.html
	// We will need the following to identify which filesystem
	// the file is on and what is the block size. For now we
	// assume "/" is the only one.
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		errLog.Println(err)
	} else {
		fsBlockSize = stat.Bsize
	}
	// fmt.Printf("Fs type [FS ID]: 0x%xd[%v]\n", stat.Type, stat.Fsid)

	// Set the unit size to 1024 if "-k" is specified
	if opts.BlockSize {
		unitSize = 1024
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
			dirSize(file)
		}
	}
}
