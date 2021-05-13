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

func Test_Main(t *testing.T) {
	err := createTestData()
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	got := dirSize("./testdata")
	if got != 4096 {
		t.Errorf("Expecting ./testdata size to be 4096 and not %v", got)
	}
}
