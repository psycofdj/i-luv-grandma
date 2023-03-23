// Copyright 2023 Xavier MARCELET. All rights reserved.
// Use of this source code is governed by Apache
// license that can be found in the LICENSE file.

// Package pbm provides a BPM format library which includes operations such
// as loading, writing and rotating images
package pbm

import (
	"io"
	"os"
	"strings"
)

// PBMMagicP1 is the special magic header of PBM black&white image format
// see https://en.wikipedia.org/wiki/Netpbm
const PBMMagicP1 string = "P1"

// Image - Represent a PBM image
type Image struct {
	width  int
	height int
	data   []bool
}

// Creates Image object from given string
func NewImageFromString(value string) (*Image, error) {
	reader := strings.NewReader(value)
	image := &Image{}
	if err := image.parse(reader); err != nil {
		return nil, err
	}
	return image, nil
}

// Creates Image object from given file path
func NewImageFromFile(path string) (*Image, error) {
	var reader io.Reader

	if path == "-" {
		reader = os.Stdin
	} else {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	}

	image := &Image{}
	if err := image.parse(reader); err != nil {
		return nil, err
	}
	return image, nil
}

// Returns image's width
func (i *Image) Width() int {
	return i.width
}

// Returns image's height
func (i *Image) Height() int {
	return i.height
}
