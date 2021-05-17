package main

import (
	"io/ioutil"
	"os"
	"testing"
)

// Creates directories and file of particular size.
func createTestData() error {
	// create sub-directories
	err := os.MkdirAll("./testdata/subdir/subsubdir", 0755)
	if err != nil {
		return err
	}

	// create files
	tmpBuf := make([]byte, 3456)
	ioutil.WriteFile("./testdata/under_4k.txt", tmpBuf, 0644)
	tmpBuf = make([]byte, 4096)
	ioutil.WriteFile("./testdata/exactly_4k.txt", tmpBuf, 0644)
	tmpBuf = make([]byte, 5678)
	ioutil.WriteFile("./testdata/over_4k.txt", tmpBuf, 0644)

	tmpBuf = make([]byte, 4096*1024)
	ioutil.WriteFile("./testdata/subdir/exactly_4m.txt", tmpBuf, 0644)

	tmpBuf = make([]byte, 4096*1024+1)
	ioutil.WriteFile("./testdata/subdir/subsubdir/over_4m.txt", tmpBuf, 0644)

	return nil
}

func Test_BuildTree(t *testing.T) {
	err := createTestData()
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	dt := dirTree{
		path:    "./testdata",
		size:    calcSize(fsBlockSize),
		files:   []fileInfo{},
		subdirs: []dirTree{},
	}
	buildDirTree(&dt)

	got := dt.size
	if got != 16432 {
		t.Errorf("Expecting ./testdata size to be 16432 and not %v", got)
	}

	got = dt.subdirs[0].size
	if got != 16392 {
		t.Errorf("Expecting %s size to be 16392 and not %v", dt.subdirs[0].path, got)
	}

	got = dt.subdirs[0].subdirs[0].size
	if got != 8200 {
		t.Errorf("Expecting %s size to be 8200 and not %v", dt.subdirs[0].subdirs[0].path, got)
	}
}
