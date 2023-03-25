package pbm

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewFromFile_doesNotExist(t *testing.T) {
	_, err := NewImageFromFile("/dos/not/exist")
	if err == nil {
		t.Fatalf("should have fail, path does not exist")
	}
}

func TestNewFromFile_valid(t *testing.T) {
	_, srcPath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("could not determine current source file path")
	}

	datasetPath := filepath.Join(filepath.Dir(srcPath), "..", "dataset")
	testPath := filepath.Join(datasetPath, "valid_480p.pbm")
	if _, err := NewImageFromFile(testPath); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

// TestNewFromFile_stdin
//
//  1. restore stdin after test
func TestNewFromFile_stdin(t *testing.T) {
	_, srcPath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("could not determine current source file path")
	}
	datasetPath := filepath.Join(filepath.Dir(srcPath), "..", "dataset")
	testPath := filepath.Join(datasetPath, "valid_480p.pbm")
	input, err := os.Open(testPath)
	if err != nil {
		t.Fatalf("could not read test file '%s': %s", testPath, err)
	}

	// 1.
	defer func(file *os.File) {
		os.Stdin = file
	}(os.Stdin)
	os.Stdin = input

	if _, err := NewImageFromFile("-"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
