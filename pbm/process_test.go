// Copyright 2023 Xavier MARCELET. All rights reserved.
// Use of this source code is governed by Apache
// license that can be found in the LICENSE file.

package pbm

import (
	"bytes"
	"testing"
)

func checkRotation(t *testing.T, angle float64, in string, expect string) {
	t.Helper()

	img, err := NewImageFromString(in)
	if err != nil {
		t.Fatalf("expected parse error: %s", err)
	}

	img.Rotate(angle)

	writer := bytes.Buffer{}
	if err := img.EncodeASCII(&writer); err != nil {
		t.Fatalf("expected encode error: %s", err)
	}

	if writer.String() != expect {
		t.Fatalf("unexpected output: %s", writer.String())
	}
}

func TestRotate_invariants(t *testing.T) {
	in := `P1
3 3
100
000
000
`
	// check invariant rotations
	checkRotation(t, 0, in, in)
	checkRotation(t, 360, in, in)
	checkRotation(t, -360, in, in)
	checkRotation(t, 720, in, in)
	checkRotation(t, -720, in, in)
}

func TestRotate_evenInvariants(t *testing.T) {
	in := `P1
2 2
10
00
`
	// check invariant rotations
	checkRotation(t, 0, in, in)
	checkRotation(t, 360, in, in)
	checkRotation(t, -360, in, in)
	checkRotation(t, 720, in, in)
	checkRotation(t, -720, in, in)
}

func TestRotate_evenSquare(t *testing.T) {
	in := `P1
2 2
10
00
`
	out := `P1
2 2
01
00
`
	checkRotation(t, 90, in, out)

	out = `P1
2 2
00
01
`
	checkRotation(t, 180, in, out)

	out = `P1
2 2
00
10
`
	checkRotation(t, -90, in, out)
}

func TestRotate_backslash(t *testing.T) {
	// backslash input
	const in = `P1
3 3
100
010
001
`

	// -> rotate backlash 45 degrees
	out := `P1
3 3
010
010
010
`
	checkRotation(t, 45, in, out)

	// -> rotate backlash 90 degrees
	out = `P1
3 3
001
010
100
`
	checkRotation(t, 90, in, out)

	// -> rotate backlash 135 degrees
	out = `P1
3 3
000
111
000
`
	checkRotation(t, 135, in, out)

	// -> rotate backlash 180 degrees
	out = `P1
3 3
100
010
001
`
	checkRotation(t, 180, in, out)

	// -> rotate backlash 180 degrees
	out = `P1
3 3
100
010
001
`
	checkRotation(t, 180, in, out)
}

func TestRotate_oddSquare(t *testing.T) {
	// odd square input
	in := `P1
5 5
00100
00100
00100
00100
00100
`
	out := `P1
5 5
00000
00000
11111
00000
00000
`
	checkRotation(t, 90, in, out)
}

func TestRotate_squareLine(t *testing.T) {
	// even square input
	in := `P1
4 4
0010
0010
0010
0010
`
	out := `P1
4 4
0000
0000
1111
0000
`
	checkRotation(t, 90, in, out)
}

func TestRotate_evenWidthRect(t *testing.T) {
	// even-width rectangle
	in := `P1
6 3
000000
011110
000000
`
	out := `P1
6 3
000100
000100
000100
`
	checkRotation(t, 90, in, out)
}

func TestRotate_evenHeight(t *testing.T) {
	// even-height rectangle
	in := `P1
3 6
010
010
010
010
010
010
`
	out := `P1
3 6
000
000
000
111
000
000
`
	checkRotation(t, 90, in, out)
}

func TestRotate_evenHeight2(t *testing.T) {
	// even-height rectangle
	in := `P1
3 6
010
010
111
010
010
010
`
	out := `P1
3 6
000
000
001
111
001
000
`
	checkRotation(t, 90, in, out)
}

func TestRotate_clock(t *testing.T) {
	// even-height rectangle
	in := `P1
11 11
00000100000
00000100000
00000100000
00000100000
00000100000
00000100000
00000000000
00000000000
00000000000
00000000000
00000000000
`
	out := `P1
11 11
00000000000
00000000010
00000000100
00000001000
00000010000
00000100000
00000000000
00000000000
00000000000
00000000000
00000000000
`
	checkRotation(t, 45, in, out)

	out = `P1
11 11
00000000000
00000000000
00000000000
00000000011
00000001100
00000110000
00000000000
00000000000
00000000000
00000000000
00000000000
`
	checkRotation(t, 67.5, in, out)
}
