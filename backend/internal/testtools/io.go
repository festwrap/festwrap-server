package testtools

import (
	"path/filepath"
	"runtime"
	"testing"
)

func GetParentDir(t *testing.T) string {
	// Get the file of the caller, not the current one
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("Could not find the file directory of the caller")
	}
	return filepath.Dir(filepath.Dir(filename))
}
