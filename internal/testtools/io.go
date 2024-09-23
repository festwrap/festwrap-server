package testtools

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func LoadTestDataOrError(t *testing.T, filePath string) []byte {
	t.Helper()

	data, err := os.ReadFile(filePath)

	if err != nil {
		t.Fatalf("Could not load data in %s: %v", filePath, err)
	}
	return data
}

func GetParentDir(t *testing.T) string {
	// Get the file of the caller, not the current one
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("Could not find the file directory of the caller")
	}
	return filepath.Dir(filepath.Dir(filename))
}

type ErrorReader struct{}

func NewErrorReader() ErrorReader {
	return ErrorReader{}
}

func (ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func (ErrorReader) Close() error {
	return nil
}
