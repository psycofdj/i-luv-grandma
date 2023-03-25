// Copyright 2023 Xavier MARCELET. All rights reserved.
// Use of this source code is governed by Apache
// license that can be found in the LICENSE file.

package pbm

import (
	"math"
)

// Rotator computes rotated coordinates for a given (x,y) point
// Stores constant values for a given rotation of a given image
type Rotator struct {
	x0      float64 // x coordinates of center of rotation
	y0      float64 // y coordinates of center of rotation
	sinθ    float64 // value of sin(angle)
	cosθ    float64 // value of cos(angle)
	offsetX float64 // shift of x coordinates for when center falls in inter-pixel
	offsetY float64 // shift of y coordinates for when center falls in inter-pixel
}

// NewRotator initializes rotator object
//
//  1. constant values for given angle and image
//  2. when center of rotation is on an inter-pixel for only one coordinates
//     translate point to half a pixel
func NewRotator(angle float64, image *Image) *Rotator {
	// 1.
	θ := angle * (math.Pi / float64(180))
	x0 := float64(image.Width()-1) / 2.0
	y0 := float64(image.Height()-1) / 2.0

	// 2.
	offsetX := 0.0
	offsetY := 0.0
	if image.Width()%2 == 1 && image.Height()%2 == 0 {
		offsetX = 0.5
	}
	if image.Width()%2 == 0 && image.Height()%2 == 1 {
		offsetY = 0.5
	}

	return &Rotator{
		x0:      x0,
		y0:      y0,
		sinθ:    math.Sin(θ),
		cosθ:    math.Cos(θ),
		offsetX: offsetX,
		offsetY: offsetY,
	}
}

// RotateCoord computes pixel new coordinates after rotation
func (r *Rotator) Compute(x int, y int) (int, int) {
	x1 := r.cosθ*(float64(x)-r.x0) - r.sinθ*(float64(y)-r.y0) + r.x0 + r.offsetX
	y1 := r.sinθ*(float64(x)-r.x0) + r.cosθ*(float64(y)-r.y0) + r.y0 + r.offsetY

	return int(math.Round(x1)), int(math.Round(y1))
}

// Rotate image to given angle
//
// Pixels projected outside image boundaries will be lost.
//
//  1. working buffer is already full of white pixels, we only need to rotate black ones
//  2. discard out-of-bout pixel coordinates
func (i *Image) Rotate(angle float64) {
	rotator := NewRotator(angle, i)

	result := make([]bool, i.width*i.height)
	for x := 0; x < i.width; x++ {
		for y := 0; y < i.height; y++ {
			// 1.
			if !i.data[x+y*i.width] {
				continue
			}

			pixelX, pixelY := rotator.Compute(x, y)

			// 2.
			if pixelX < 0 || pixelX >= i.width || pixelY < 0 || pixelY >= i.height {
				continue
			}

			index := pixelX + pixelY*i.width
			result[index] = true
		}
	}
	i.data = result
}
