// complex64 renders the Mandelbrot fractal using complex64.
package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"strconv"
)

var colors [256]color.RGBA
var rainbow = []color.RGBA{
	color.RGBA{255, 0, 0, 255},
	color.RGBA{148, 0, 211, 255},
	color.RGBA{75, 0, 130, 255},
	color.RGBA{0, 0, 255, 255},
	color.RGBA{0, 255, 0, 255},
	color.RGBA{255, 255, 0, 255},
	color.RGBA{255, 127, 0, 255},
	color.RGBA{255, 0, 0, 255},
}

func init() {
	for i := range colors {
		colors[i] = interp(i, len(colors), len(rainbow)-1)
	}
}

// interp maps i/length to 0..mapTo, then interpolates the rainbow colours at
// the indices above and below.
func interp(i, length, mapTo int) color.RGBA {
	lower := i * mapTo / length
	weight := float64(i*mapTo)/float64(length) - float64(lower)
	r := rainbow[lower].R + uint8(weight*float64(rainbow[lower+1].R-rainbow[lower].R))
	g := rainbow[lower].G + uint8(weight*float64(rainbow[lower+1].G-rainbow[lower].G))
	b := rainbow[lower].B + uint8(weight*float64(rainbow[lower+1].B-rainbow[lower].B))
	return color.RGBA{r, g, b, 255}
}

func main() {
	const width, height = 1024, 1024
	xmin64, ymin64, xmax64, ymax64 := -2.0, -2.0, 2.0, 2.0
	if len(os.Args) == 4 {
		xmin64, _ = strconv.ParseFloat(os.Args[1], 32)
		ymin64, _ = strconv.ParseFloat(os.Args[2], 32)
		xmax64, _ = strconv.ParseFloat(os.Args[3], 32)
		ymax64 = ymin64 + (xmax64 - xmin64)
	}
	xmin, ymin, xmax, ymax := float32(xmin64), float32(ymin64), float32(xmax64), float32(ymax64)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float32(py)/float32(height)*(ymax-ymin) + ymin
		ydiff := (ymax - ymin) / float32(4*height)
		for px := 0; px < width; px++ {
			x := float32(px)/float32(width)*(xmax-xmin) + xmin
			xdiff := (xmax - xmin) / float32(4*width)

			var colVec []color.Color
			colVec = append(colVec, mandelbrot(complex(x-xdiff, y-ydiff)))
			colVec = append(colVec, mandelbrot(complex(x-xdiff, y+ydiff)))
			colVec = append(colVec, mandelbrot(complex(x+xdiff, y-ydiff)))
			colVec = append(colVec, mandelbrot(complex(x+xdiff, y+ydiff)))

			// Image point (px, py) represents complex value z.
			img.Set(px, py, average(colVec))
		}
	}
	png.Encode(os.Stdout, img) // NOTE: ignoring errors
}

func mandelbrot(z complex64) color.Color {
	const iterations = 200
	const contrast = 7

	var v complex64
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(complex128(v)) > 2 {
			return colors[255-contrast*n]
		}
	}
	return color.Black
}

func average(colVec []color.Color) color.Color {
	var rSum, gSum, bSum uint32
	for _, v := range colVec {
		r, g, b, _ := v.RGBA()
		rSum += r
		gSum += g
		bSum += b
	}
	nel := uint32(len(colVec))
	rAvg := uint16(rSum / nel)
	bAvg := uint16(gSum / nel)
	gAvg := uint16(bSum / nel)
	return color.RGBA64{rAvg, bAvg, gAvg, 0xffff}
}
