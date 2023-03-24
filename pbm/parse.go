// Copyright 2023 Xavier MARCELET. All rights reserved.
// Use of this source code is governed by Apache
// license that can be found in the LICENSE file.

package pbm

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// extract next valid token from buffer ignoring comment, whitespaces and newlines
//  1. atEOF tells us if that no more data is available in stream
//  2. increment inside loop, advance must have incremented value in switch returns
//  3. eat up anything until end of comment or end of input
func nextToken(data []byte, atEOF bool) (int, []byte, error) {
	var (
		err       error
		available = len(data)
		advance   = 0
		token     = []byte{}
	)

	// 1.
	if atEOF {
		err = bufio.ErrFinalToken
	}

	// 2.
	for advance < available {
		char := data[advance]
		advance++
		switch {
		case char == ' ', char == '\n':
			if len(token) != 0 {
				return advance, token, err
			}
		case char == '#':
			// 3.
			for advance < available && char != '\n' {
				char = data[advance]
				advance++
			}
		case char >= '0' || char <= '9':
			token = append(token, char)
		default:
			return advance, token, fmt.Errorf("unexpected input '%c'", char)
		}
	}

	return advance, token, err
}

func (i *Image) parse(stream io.Reader) error {
	if err := i.parseMagic(stream); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stream)
	scanner.Split(nextToken)

	if err := i.parseHeader(scanner); err != nil {
		return err
	}
	return i.parseData(scanner)
}

func (i *Image) parseMagic(stream io.Reader) error {
	buffer := make([]byte, 2)
	count, err := stream.Read(buffer)
	if err != nil || count != 2 {
		return fmt.Errorf("invalid format, expected magic number")
	}
	if string(buffer) != PBMMagicP1 {
		return fmt.Errorf("invalid magic number '%s', expecting %s", string(buffer), PBMMagicP1)
	}
	return nil
}

func (i *Image) parseHeader(scanner *bufio.Scanner) error {
	var err error

	if !scanner.Scan() || scanner.Err() != nil {
		return fmt.Errorf("invalid format, expected image width: %s", scanner.Err())
	}
	i.width, err = strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("invalid width '%s', expecting number", scanner.Text())
	}

	if !scanner.Scan() || scanner.Err() != nil {
		return fmt.Errorf("invalid format, expected image height: %s", scanner.Err())
	}
	i.height, err = strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("invalid height '%s', expecting number", scanner.Text())
	}

	return nil
}

func (i *Image) parseData(scanner *bufio.Scanner) error {
	i.data = make([]bool, i.width*i.height)
	index := 0
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("invalid input: %s", err)
		}
		for _, digit := range scanner.Text() {
			if digit != '0' && digit != '1' {
				return fmt.Errorf("invalid pixel value '%d', expecting 0 or 1", digit)
			}
			if index >= (i.width * i.height) {
				return fmt.Errorf("invalid data, expecting no more than %d pixels", i.width*i.height)
			}
			i.data[index] = (digit == '1')
			index++
		}
	}

	if index != (i.width * i.height) {
		return fmt.Errorf("invalid data, got '%d' out of '%d' expected pixels", index, (i.width * i.height))
	}

	return nil
}
