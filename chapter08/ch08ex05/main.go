// Ch08ex05 is a concurrent version of the Mandelbrot program from Chapter 3.3.
// Usage:
//
//     ch08ex05 [-P NPROC]
//
// where NPROC is the number of goroutines to be used.
package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 5000, 5000
	)

	var nProcs int
	flag.IntVar(&nProcs, "P", 1, "number of goroutines")
	flag.Parse()

	var (
		img  = image.NewRGBA(image.Rect(0, 0, width, height))
		rows = make(chan int)
		wg   = &sync.WaitGroup{}
	)

	for i := 0; i < nProcs; i++ {
		wg.Add(1)
		go func(rows <-chan int) {
			for py := range rows {
				y := float64(py)/height*(ymax-ymin) + ymin
				for px := 0; px < width; px++ {
					x := float64(px)/width*(xmax-xmin) + xmin
					c := mandelbrot(complex(x, y))
					img.Set(px, py, c)
				}
			}
			wg.Done()
		}(rows)
	}

	go func() {
		for py := 0; py < height; py++ {
			rows <- py
		}
		close(rows)
	}()

	wg.Wait()
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
