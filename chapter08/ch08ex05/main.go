// Ch08ex05 is a concurrent version of the Mandelbrot program from Chapter 3.3.
// The channel is buffered to 10 elements per goroutine.
package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

type pixel struct {
	x, y int
	c    color.Color
}

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 5000, 5000
	)

	var nProcs int
	flag.IntVar(&nProcs, "P", 1, "number of goroutines")
	flag.Parse()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	pixels := make(chan pixel, 10*nProcs)
	procHeight := height / nProcs

	for i := 0; i < nProcs; i++ {
		minY := i * procHeight
		var maxY int
		if i == nProcs-1 {
			maxY = height
		} else {
			maxY = (i + 1) * procHeight
		}

		go func(minY, maxY int) {
			for py := minY; py < maxY; py++ {
				y := float64(py)/height*(ymax-ymin) + ymin
				for px := 0; px < width; px++ {
					x := float64(px)/width*(xmax-xmin) + xmin
					z := complex(x, y)
					pixels <- pixel{
						x: px,
						y: py,
						c: mandelbrot(z),
					}
				}
			}
		}(minY, maxY)
	}

	for i := 0; i < width*height; i++ {
		pixel := <-pixels
		img.Set(pixel.x, pixel.y, pixel.c)
	}

	png.Encode(os.Stdout, img) // NOTE: ignoring errors
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
