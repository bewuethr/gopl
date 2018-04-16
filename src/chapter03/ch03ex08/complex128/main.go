// complex128 renders the Mandelbrot fractal using complex128.
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
	xmin, ymin, xmax, ymax := -2.0, -2.0, 2.0, 2.0
	if len(os.Args) == 4 {
		xmin, _ = strconv.ParseFloat(os.Args[1], 64)
		ymin, _ = strconv.ParseFloat(os.Args[2], 64)
		xmax, _ = strconv.ParseFloat(os.Args[3], 64)
		ymax = ymin + (xmax - xmin)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/float64(height)*(ymax-ymin) + ymin
		ydiff := (ymax - ymin) / float64(4*height)
		for px := 0; px < width; px++ {
			x := float64(px)/float64(width)*(xmax-xmin) + xmin
			xdiff := (xmax - xmin) / float64(4*width)

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

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 7

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
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
