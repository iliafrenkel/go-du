package dirtree

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

type options struct {
	BlockSize       bool
	CountFiles      bool
	DereferenceAll  bool
	DereferenceArgs bool
	OneFileSystem   bool
	Summarise       bool
	Version         bool
}

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
	output   []string   // expected output
}

var testCases = []testCase{
	{
		name:  "File does not exist",
		files: []testFile{},
		path:  "./ak5i8fg74",
		opts: options{
			BlockSize:       false,
			CountFiles:      false,
			DereferenceAll:  false,
			DereferenceArgs: false,
			OneFileSystem:   false,
			Summarise:       false,
		},
		expected: 0,
		output: []string{
			"0\t./ak5i8fg74",
		},
	},
	{
		name:  "Access denied",
		files: []testFile{},
		path:  "/root",
		opts: options{
			BlockSize:       false,
			CountFiles:      false,
			DereferenceAll:  false,
			DereferenceArgs: false,
			OneFileSystem:   false,
			Summarise:       false,
		},
		expected: 8,
		output: []string{
			"8\t/root",
		},
	},
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
		expected: 3456 + 4096,
		output:   []string{"16\t" + testFilesRoot},
	},
	{
		name: "Single file argument",
		files: []testFile{
			{filepath.Join(testFilesRoot, "under_4k.txt"), 3456},
		},
		path: testFilesRoot + "/under_4k.txt",
		opts: options{
			BlockSize:       false,
			CountFiles:      false,
			DereferenceAll:  false,
			DereferenceArgs: false,
			OneFileSystem:   false,
			Summarise:       false,
		},
		expected: 4096,
		output:   []string{"8\t" + testFilesRoot + "/under_4k.txt"},
	},
	{
		name: "Single file and -k",
		files: []testFile{
			{filepath.Join(testFilesRoot, "under_4k.txt"), 3456},
		},
		path: testFilesRoot,
		opts: options{
			BlockSize:       true, // <--
			CountFiles:      false,
			DereferenceAll:  false,
			DereferenceArgs: false,
			OneFileSystem:   false,
			Summarise:       false,
		},
		expected: 3456 + 4096,
		output:   []string{"8\t" + testFilesRoot},
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
		expected: 4096 + 4096 + 8192 + 4096,
		output:   []string{"40\t" + testFilesRoot},
	},
	{
		name: "Miltiple directories",
		files: []testFile{
			{filepath.Join(testFilesRoot, "under_4k.txt"), 3456},
			{filepath.Join(testFilesRoot, "exactly_4k.txt"), 4096},
			{filepath.Join(testFilesRoot, "over_4k.txt"), 5678},
			{filepath.Join(testFilesRoot, "subdir", "over_4m.txt"), 5678 * 1024},
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
		expected: 4096 + 4096 + 8192 + 5818367 + 4096,
		output: []string{
			"11368\t" + fixPath(filepath.Join(testFilesRoot, "subdir")),
			"11408\t" + testFilesRoot,
		},
	},
	{
		name: "Summarise",
		files: []testFile{
			{filepath.Join(testFilesRoot, "under_4k.txt"), 3456},
			{filepath.Join(testFilesRoot, "exactly_4k.txt"), 4096},
			{filepath.Join(testFilesRoot, "over_4k.txt"), 5678},
			{filepath.Join(testFilesRoot, "subdir", "over_4m.txt"), 5678 * 1024},
		},
		path: testFilesRoot,
		opts: options{
			BlockSize:       false,
			CountFiles:      false,
			DereferenceAll:  false,
			DereferenceArgs: false,
			OneFileSystem:   false,
			Summarise:       true,
		},
		expected: 4096 + 4096 + 8192 + 5818367 + 4096,
		output: []string{
			"11408\t" + testFilesRoot,
		},
	},
	{
		name: "All files",
		files: []testFile{
			{filepath.Join(testFilesRoot, "under_4k.txt"), 3456},
			{filepath.Join(testFilesRoot, "exactly_4k.txt"), 4096},
			{filepath.Join(testFilesRoot, "over_4k.txt"), 5678},
			{filepath.Join(testFilesRoot, "subdir", "over_4m.txt"), 5678 * 1024},
		},
		path: testFilesRoot,
		opts: options{
			BlockSize:       false,
			CountFiles:      true,
			DereferenceAll:  false,
			DereferenceArgs: false,
			OneFileSystem:   false,
			Summarise:       false,
		},
		expected: 4096 + 4096 + 8192 + 5818367 + 4096,
		output: []string{
			"8\t" + testFilesRoot + "/exactly_4k.txt",
			"16\t" + testFilesRoot + "/over_4k.txt",
			"8\t" + testFilesRoot + "/under_4k.txt",
			"11360\t" + testFilesRoot + "/subdir/over_4m.txt",
			"11368\t" + testFilesRoot + "/subdir",
			"11408\t" + testFilesRoot,
		},
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
func resetTestData() error {
	if err := os.RemoveAll(testFilesRoot); err != nil {
		return err
	}

	return nil
}

func Test_BuildTree(t *testing.T) {
	for _, tc := range testCases {
		err := createTestData(tc.files)
		if err != nil {
			t.Fatalf("Failed to create test data: %v", err)
		}
		t.Run(tc.name, func(t *testing.T) {
			defer resetTestData()
			var unitSize int64 = 512
			if tc.opts.BlockSize {
				unitSize = 1024
			}
			dt := New(tc.path, unitSize)
			got := dt.size
			want := dt.calcSize(tc.expected)
			if got != want {
				t.Errorf("Expecting size to be %v and not %v", want, got)
			}
			out := dt.PrintDirTree("%d\t%s", tc.opts.CountFiles, tc.opts.Summarise)
			for i, v := range out {
				if v != tc.output[i] {
					t.Errorf("Output string #%v '%s' is not equal to the expected one '%s'", i, v, tc.output[i])
				}
			}
		})
	}
}
