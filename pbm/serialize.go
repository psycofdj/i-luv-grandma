// Copyright 2023 Xavier MARCELET. All rights reserved.
// Use of this source code is governed by Apache
// license that can be found in the LICENSE file.

package pbm

import (
	"fmt"
	"io"
	"os"
)

// Serialize image into file in ascii/plain representation
func (i *Image) EncodeASCIIToFile(path string) error {
	var writer io.Writer

	if path == "-" {
		writer = os.Stdout
	} else {
		flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		file, err := os.OpenFile(path, flags, 0644)
		if err != nil {
			return err
		}
		writer = file
	}
	return i.EncodeASCII(writer)
}

// Serialize image into stream in ascii/plain representation
func (i *Image) EncodeASCII(stream io.Writer) error {
	if err := i.encodeASCIIHeader(stream); err != nil {
		return err
	}
	return i.encodeASCIIData(stream)
}

func (i *Image) encodeASCIIHeader(stream io.Writer) error {
	_, err := fmt.Fprintf(stream, "%s\n", PBMMagicP1)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(stream, "%d %d\n", i.Width(), i.Height())
	if err != nil {
		return err
	}
	return nil
}

// serialize data section
//
// This function serialize image into memory then writes to output stream
// for performance considerations
//
//  1. allocates image-size plus enough room for newlines written every "width" bytes
//  2. cOffset tracks how many separator were written so far which affects byte index
//     to write in result variable
func (i *Image) encodeASCIIData(stream io.Writer) error {
	var (
		// 1.
		result    = make([]byte, (i.Width()+1)*i.Height())
		separator = byte(10)
		white     = byte(48)
		black     = byte(49)
	)
	for cOffset, cIdx := 0, 0; cIdx < i.Width()*i.Height(); cIdx++ {
		if i.data[cIdx] {
			result[cIdx+cOffset] = black
		} else {
			result[cIdx+cOffset] = white
		}
		// 2.
		if (cIdx+1)%i.Width() == 0 {
			cOffset++
			result[cIdx+cOffset] = separator
		}
	}
	if _, err := stream.Write(result); err != nil {
		return err
	}
	return nil
}
