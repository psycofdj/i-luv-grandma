// Copyright 2023 Xavier MARCELET. All rights reserved.
// Use of this source code is governed by Apache
// license that can be found in the LICENSE file.

package pbm

import (
	"math"
)

// Rotate image to given angle
//
// Pixels projected outside image boundaries will be lost.
//
//  1. working buffer is already full of black pixels, we only need to rotate white ones
//  2. rotation equation where: θ angle of rotation, x,y: pixel coordinates, x0,y0: center of rotation
//  3. pathological case where pixel coordinates is 0.5, which is rounded to 1 instead of 0
//  4. approximate to closest integer coordinates
//  5. discard out-of-bout pixel coordinates
func (i *Image) Rotate(angle float64) {
	θ := angle * (math.Pi / float64(180))
	x0 := float64(i.width-1) / 2.0
	y0 := float64(i.height-1) / 2.0
	sinθ := math.Sin(θ)
	cosθ := math.Cos(θ)

	result := make([]bool, i.width*i.height)
	for x := 0; x < i.width; x++ {
		for y := 0; y < i.height; y++ {
			// 1.
			if !i.data[x+y*i.width] {
				continue
			}

			// 2.
			newX := cosθ*(float64(x)-x0) - sinθ*(float64(y)-y0) + x0
			newY := sinθ*(float64(x)-x0) + cosθ*(float64(y)-y0) + y0

			// 3.
			if (newY + 0.5) < math.SmallestNonzeroFloat64 {
				newY = 0.0
			}
			if (newX + 0.5) < math.SmallestNonzeroFloat64 {
				newX = 0.0
			}

			// 4.
			pixelX := int(math.Round(newX))
			pixelY := int(math.Round(newY))

			// 5.
			if pixelX < 0 || pixelX >= i.width {
				continue
			}
			if pixelY < 0 || pixelY >= i.height {
				continue
			}

			index := pixelX + pixelY*i.width
			result[index] = true
		}
	}
	i.data = result
}
