// Copyright 2023 Xavier MARCELET. All rights reserved.
// Use of this source code is governed by Apache
// license that can be found in the LICENSE file.

package pbm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSerialize_simple(t *testing.T) {
	// simple test
	image := Image{2, 2, []bool{true, true, false, false}}
	expect := `P1
2 2
11
00
`
	content := new(strings.Builder)
	err := image.EncodeASCII(content)
	if err != nil {
		t.Fatalf("unexpected serialization error: %s", err)
	}
	if expect != content.String() {
		t.Fatalf("unexpected serialization output: %v, want %v", []byte(content.String()), []byte(expect))
	}
}

func TestSerialize_closeWriter(t *testing.T) {
	// output to closed writer
	image := Image{2, 2, []bool{true, true, false, false}}
	file, _ := os.CreateTemp("dir", "prefix")
	file.Close()
	err := image.EncodeASCII(file)
	if err == nil {
		t.Fatalf("should have fail: cannot write to closed file")
	}
}

func TestSerialize_writeFile(t *testing.T) {
	// actually write to file
	image := Image{3, 3, []bool{
		true, true, true,
		false, false, false,
		true, false, true,
	}}
	expect := `P1
3 3
111
000
101
`
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("could not create temp directory")
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "output")
	err = image.EncodeASCIIToFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	output, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read generated file path '%s': %s", path, err)
	}
	if string(output) != expect {
		t.Fatalf("unexpected serialization output: %v, want %v", output, []byte(expect))
	}
}
