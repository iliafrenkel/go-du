package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// The below is needed because "packages that call flag.Parse during package
// initialization may cause tests to fail" as explained here under "testing":
// https://golang.org/doc/go1.13#testing
// Workaround suggested here:
// https://github.com/onsi/ginkgo/issues/602#issuecomment-555502868
var _ = func() bool {
	testing.Init()
	return true
}()

const testFilesRoot = "./testdata"

// Represents a single file of a given size
type testFile struct {
	path string // relative or absolute path
	size int64  // size in bytes
}

// Represents a test case - one or multiple files in one multiple directories
// with the expected total size.
type testCase struct {
	name     string     // friendly name for the test
	files    []testFile // a set of test files
	path     string     // path for the test
	opts     options    // simulated command line options
	expected int64      // expected size in bytes
}

var testCases = []testCase{
	{
		name: "Single file",
		files: []testFile{
			{filepath.Join(testFilesRoot, "under_4k.txt"), 3456},
		},
		path: testFilesRoot,
		opts: options{
			BlockSize:       false,
			CountFiles:      false,
			DereferenceAll:  false,
			DereferenceArgs: false,
			OneFileSystem:   false,
			Summarise:       false,
		},
		expected: 3456,
	},
	{
		name: "Single directory",
		files: []testFile{
			{filepath.Join(testFilesRoot, "under_4k.txt"), 3456},
			{filepath.Join(testFilesRoot, "exactly_4k.txt"), 4096},
			{filepath.Join(testFilesRoot, "over_4k.txt"), 5678},
		},
		path: testFilesRoot,
		opts: options{
			BlockSize:       false,
			CountFiles:      false,
			DereferenceAll:  false,
			DereferenceArgs: false,
			OneFileSystem:   false,
			Summarise:       false,
		},
		expected: 4096 + 4096 + 8192,
	},
}

// Creates and empty file of a given size with a given name.
//
// The `file` must be a full path, can be either relative or absolute.
// If the path doesn't exist it will be created. If the file already
// exists it will be trunctated to a given size.
// Permissions for created folders are 0755 and permissions for the
// file are 0644.
func createDummyFile(file string, size int64) error {
	// check if file exists
	path := filepath.Dir(file)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	// create a temp buffer of a given size
	tmpBuf := make([]byte, size)
	// write the file, truncate to a given size if the file already exists
	return ioutil.WriteFile(file, tmpBuf, 0644)
}

// Creates test directories and files.
func createTestData(testFiles []testFile) error {
	for _, f := range testFiles {
		if err := createDummyFile(f.path, f.size); err != nil {
			return err
		}
	}

	return nil
}

// Delete test directories and files.
func deleteTestData(testFiles []testFile) error {
	for _, f := range testFiles {
		if err := os.RemoveAll(f.path); err != nil {
			return err
		}
	}

	return nil
}

// Test default values for flags.
func Test_Flags(t *testing.T) {
	var tests = []struct {
		name string
		flag bool
		want bool
	}{
		{"-k", opts.BlockSize, false},
		{"-a", opts.CountFiles, false},
		{"-L", opts.DereferenceAll, false},
		{"-H", opts.DereferenceArgs, false},
		{"-x", opts.OneFileSystem, false},
		{"-s", opts.Summarise, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.flag
			if got != tt.want {
				t.Errorf("Expecting %v to be %v and not %v if %s is not provided.", tt.flag, tt.want, got, tt.name)
			}
		})
	}
}

func Test_BuildTree(t *testing.T) {
	for _, tc := range testCases {
		err := createTestData(tc.files)
		if err != nil {
			t.Fatalf("Failed to create test data: %v", err)
		}
		t.Run(tc.name, func(t *testing.T) {
			dt := dirTree{
				path:    tc.path,
				size:    0,
				files:   []fileInfo{},
				subdirs: []dirTree{},
			}
			opts = tc.opts
			buildDirTree(&dt)
			got := dt.size
			want := calcSize(tc.expected)
			if got != want {
				t.Errorf("Expecting size to be %v and not %v", want, got)
			}
		})
		err = deleteTestData(tc.files)
		if err != nil {
			t.Fatalf("Failed to delete test data: %v", err)
		}
	}
}
