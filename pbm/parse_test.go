// Copyright 2023 Xavier MARCELET. All rights reserved.
// Use of this source code is governed by Apache
// license that can be found in the LICENSE file.

package pbm

import (
	"testing"
)

func TestParse_empty(t *testing.T) {
	input := ""
	_, err := NewImageFromString(input)
	if err == nil {
		t.Fatalf("should have fail: missing magic number")
	}
}

func TestParse_missingMagic(t *testing.T) {
	input := "2 2 0000"
	_, err := NewImageFromString(input)
	if err == nil {
		t.Fatalf("should have fail: missing magic number")
	}
}

func TestParse_invalidMagic(t *testing.T) {
	input := " \nP2 2 2 0000"
	_, err := NewImageFromString(input)
	if err == nil {
		t.Fatalf("should have fail: missing magic number (shoud be first bytes)")
	}
}

func TestParse_incorrectMagic(t *testing.T) {
	input := "P4 2 2 0000"
	_, err := NewImageFromString(input)
	if err == nil {
		t.Fatalf("should have fail: incorrect magic number")
	}
}

func TestParse_incorrectWidth(t *testing.T) {
	input := "P1 x 2 0000"
	_, err := NewImageFromString(input)
	if err == nil {
		t.Fatalf("should have fail: invalid width")
	}
}

func TestParse_incorrectHeight(t *testing.T) {
	input := "P1 2 x 0000"
	_, err := NewImageFromString(input)
	if err == nil {
		t.Fatalf("should have fail: invalid height")
	}
}

func expect(t *testing.T, image *Image, err error, width, height int, data string) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected parse error: %s", err)
	}
	if width != image.width {
		t.Fatalf("expected width %d, got '%d'", width, image.width)
	}
	if width != image.height {
		t.Fatalf("expected height %d, got '%d'", height, image.height)
	}
	for cIdx, cPixel := range image.data {
		want := "0"
		if cPixel {
			want = "1"
		}
		if string(data[cIdx]) != want {
			t.Fatalf("unexpected data result, what '%s', got '%v'", data, image.data)
		}
	}
}

func TestParse_simple(t *testing.T) {
	// simple input
	input := `P1 2 2 1 0 0 1`
	image, err := NewImageFromString(input)
	expect(t, image, err, 2, 2, "1001")
}

func TestParse_newlines(t *testing.T) {
	// newlines
	input := `P1
2 2
1 0 0 1`
	image, err := NewImageFromString(input)
	expect(t, image, err, 2, 2, "1001")
}

func TestParse_endLineFeed(t *testing.T) {
	// with end line-feed
	input := `P1
2 2
1 0 0 1
`
	image, err := NewImageFromString(input)
	expect(t, image, err, 2, 2, "1001")
}

func TestParse_leadingSpaces(t *testing.T) {
	// with starting spaces
	input := `P1
           2 2
           1 0 0 1`
	image, err := NewImageFromString(input)
	expect(t, image, err, 2, 2, "1001")
}

func TestParse_multipleComments(t *testing.T) {
	// newlines with multiple comments
	input := `P1 # this is a magic header
2 2 # my grandma's dearest memory
1 0 0 1
`
	image, err := NewImageFromString(input)
	expect(t, image, err, 2, 2, "1001")
}

func TestParse_emptyComment(t *testing.T) {
	// newlines with empty
	input := `P1 #
3 3 #
1 0 1
0 1 1
1 0 1
`
	image, err := NewImageFromString(input)
	expect(t, image, err, 3, 3, "101011101")
}

func TestParse_dataComment(t *testing.T) {
	// newlines with comment before data
	input := `P1
2 2
# my grandma's dearest memory
1 0 0 1
`
	image, err := NewImageFromString(input)
	expect(t, image, err, 2, 2, "1001")
}

func TestParse_compactedData(t *testing.T) {
	// compacted data
	input := `P1 2 2 1001`
	image, err := NewImageFromString(input)
	expect(t, image, err, 2, 2, "1001")
}

func TestParse_tooMuchData(t *testing.T) {
	input := "P1 2 2 00000"
	_, err := NewImageFromString(input)
	if err == nil {
		t.Fatalf("should have fail: too much data")
	}
}

func TestParse_invalidData(t *testing.T) {
	input := "P1 2 2 020"
	_, err := NewImageFromString(input)
	if err == nil {
		t.Fatalf("should have fail: invalid digit 2")
	}
}

func TestParse_missingData(t *testing.T) {
	input := "P1 2 2 01"
	_, err := NewImageFromString(input)
	if err == nil {
		t.Fatalf("should have fail: not enough data")
	}
}
