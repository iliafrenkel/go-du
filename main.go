package main

import (
	"flag"
	"fmt"
	"os"
)

var all bool

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-a|-s] [-kx] [-H|-L] [file...]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVar(&all, "a", false,
		`In addition to the default output, report the size of
each file not of type directory in the file hierarchy
rooted in the specified file.  The -a option shall not
affect whether non-directories given as file operands
are listed.`)

	flag.Parse()

	if all {
		fmt.Println("All flag is enabled.")
	} else {
		fmt.Println("All flag is disabled.")
	}

	args := flag.Args()
	if len(args) == 0 {
		os.Exit(1)
	}
	fmt.Println(args)
}
