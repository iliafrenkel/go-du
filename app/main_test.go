package main

import (
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
		{"-v", opts.Version, false},
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

func Test_FlagConflicts(t *testing.T) {
	opts = options{
		BlockSize:       false,
		CountFiles:      true,
		DereferenceAll:  false,
		DereferenceArgs: false,
		OneFileSystem:   false,
		Summarise:       true,
	}
	if !conflictingFlags() {
		t.Errorf("Expecting conflict between -a and -s flags.")
	}
}
