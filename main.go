package main

import (
	"flag"
	"fmt"
	"os"
)

var opts struct {
	BlockSize       bool `short:"k" default:"false" description:"Write the files sizes in units of 1024 bytes, rather than the default 512-byte units"`
	CountFiles      bool `short:"a" long:"all" default:"false" description:"write counts for all files, not just directories"`
	DereferenceAll  bool `short:"L" long:"dereference" default:"false" description:"dereference all symbolic links"`
	DereferenceArgs bool `short:"H" long:"dereference-args" default:"false" description:"dereference only symlinks that are listed on the command line"`
	OneFileSystem   bool `short:"x" long:"one-files-system" default:"false" description:"skip directories on different file systems"`
	Summarise       bool `short:"s" long:"summarise" default:"false" description:"display only a total for each argument"`
}

var files []string

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
	files = flag.Args()
	if len(files) == 0 {
		files = append(files, ".")
	}

	fmt.Println(opts)
	fmt.Println(files)
}
