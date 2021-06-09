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

	"github.com/iliafrenkel/go-du/app/dirtree"
)

// Format for printing out dir/file entry
const outFormat = "%d\t%s"

// Command-line flags
type options struct {
	BlockSize       bool `short:"k" default:"false" description:"Write the files sizes in units of 1024 bytes, rather than the default 512-byte units"`
	CountFiles      bool `short:"a" long:"all" default:"false" description:"write counts for all files, not just directories"`
	DereferenceAll  bool `short:"L" long:"dereference" default:"false" description:"dereference all symbolic links"`
	DereferenceArgs bool `short:"H" long:"dereference-args" default:"false" description:"dereference only symlinks that are listed on the command line"`
	OneFileSystem   bool `short:"x" long:"one-file-system" default:"false" description:"skip directories on different file systems"`
	Summarise       bool `short:"s" long:"summarise" default:"false" description:"display only a total for each argument"`
	Version         bool `short:"v" long:"version" default:"false" description:"show version info and exit"`
}

var opts options

// Version information, comes from the build flags (see Makefile)
var (
	revision = "unknown"
	version  = "unknown"
	branch   = "unknown"
)

// Holds the file/directory names from the command line arguments.
var argFiles []string

// A logger that outputs to stderr without the timestamp.
var errLog = log.New(os.Stderr, "", 0)

// conflictingFlags checks that there are no conflicts between various command line
// flags. Prints out an error message and returns true if there are some
// conflicting flags, returns false otherwise.
func conflictingFlags() bool {
	if opts.CountFiles && opts.Summarise {
		errLog.Println("Cannot both summarise and show all entries.")
		return true
	}

	return false
}

// printVersion prints out version, license and contact information.
func printVersion() {
	fmt.Println("go-du", version)
	fmt.Println("Copyright (c) 2021 Ilia Frenkel")
	fmt.Println("MIT License <https://opensource.org/licenses/MIT>")
	fmt.Println("Source code <https://github.com/iliafrenkel/go-du/>")
	fmt.Println("\nWritten by Ilia Frenkel<frenkel.ilia@gmail.com>")
	fmt.Println()
}

// Declare and parse command line flags.
func init() {
	// Define command-line flags
	flag.Usage = func() {
		fmt.Println("Usage: go-du [-a|-s] [-kx] [-H|-L] [FILE...]")
		fmt.Println("Summarise disk usage of the set of FILEs, recursively for directories.")
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println("\nThis is POSIX compatible implementation of the du utility. For exended")
		fmt.Println("documentation see https://man7.org/linux/man-pages/man1/du.1p.html")
		fmt.Println("\nDisplay values are in 512-byte units, rounded up to the next 512-byte unit")
		fmt.Println("unless -k flag is specified.")
		fmt.Println("\nCreated by Ilia Frenkel<frenkel.ilia@gmail.com>")
		fmt.Println("Report bugs at https://github.com/iliafrenkel/go-du")
		fmt.Printf("Revision: %s\n", revision)
	}
	flag.BoolVar(&opts.BlockSize, "k", false, "\tWrite the files sizes in units of 1024 bytes, rather than the\n\tdefault 512-byte units")
	flag.BoolVar(&opts.CountFiles, "a", false, "\twrite counts for all files, not just directories")
	flag.BoolVar(&opts.DereferenceAll, "L", false, "\tdereference all symbolic links")
	flag.BoolVar(&opts.DereferenceArgs, "H", false, "\tdereference only symlinks that are listed on the command line")
	flag.BoolVar(&opts.OneFileSystem, "x", false, "\tskip directories on different file systems")
	flag.BoolVar(&opts.Summarise, "s", false, "\tdisplay only a total for each argument")
	flag.BoolVar(&opts.Version, "version", false, "\t")
	flag.BoolVar(&opts.Version, "v", false, "\tshow version info and exit")
	flag.Parse()

	if conflictingFlags() {
		os.Exit(1)
	}

	// If version is requested print out the info and ignore all other flags
	if opts.Version {
		printVersion()
		os.Exit(0)
	}
}

func main() {
	// If there are no arguments provided, default to the current directory
	argFiles = flag.Args()
	if len(argFiles) == 0 {
		argFiles = append(argFiles, ".")
	}

	// If -k is provided set block size to 1024
	var bs int64 = 512
	if opts.BlockSize {
		bs = 1024
	}

	for _, file := range argFiles {
		dt := dirtree.New(file, bs)
		for _, s := range dt.PrintDirTree(outFormat, opts.CountFiles, opts.Summarise) {
			fmt.Println(s)
		}
	}
}
