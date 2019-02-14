// Ch03ex07 generates a Newton fractal for z^4 - 1 = 0 with supersampling.
package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

var rootCols = [4]color.RGBA{
	color.RGBA{255, 0, 0, 255},
	color.RGBA{0, 255, 0, 255},
	color.RGBA{0, 0, 255, 255},
	color.RGBA{255, 127, 0, 255},
}

func main() {
	const (
		xmin, ymin, xmax, ymax = -2.0, -2.0, 2.0, 2.0
		width, height          = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		ydiff := (ymax - ymin) / (4 * height)
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			xdiff := (xmax - xmin) / (4 * width)

			var colVec []color.Color
			colVec = append(colVec, newton(complex(x-xdiff, y-ydiff)))
			colVec = append(colVec, newton(complex(x-xdiff, y+ydiff)))
			colVec = append(colVec, newton(complex(x+xdiff, y-ydiff)))
			colVec = append(colVec, newton(complex(x+xdiff, y+ydiff)))

			// Image point (px, py) represents complex value z.
			img.Set(px, py, average(colVec))
		}
	}
	png.Encode(os.Stdout, img) // NOTE: ignoring errors
}

// newton solves z^4 - 1 = 0.
func newton(z complex128) color.Color {
	const (
		iterations = 200
		tol        = 1e-6
	)

	roots := []complex128{1, -1, 1i, -1i}

	for n := 0; n < iterations; n++ {
		z -= (cmplx.Pow(z, 4) - 1) / (4 * cmplx.Pow(z, 3))
		for i, v := range roots {
			if cmplx.Abs(z-v) < tol {
				return shade(rootCols[i], n, iterations)
			}
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
	rAvg, bAvg, gAvg := uint16(rSum/nel), uint16(gSum/nel), uint16(bSum/nel)
	return color.RGBA64{rAvg, bAvg, gAvg, 0xffff}
}

func shade(c color.Color, n, nmax int) color.Color {
	const contrast = 4
	r, g, b, _ := c.RGBA()
	rNew := float64(r) * (1 - contrast*float64(n)/float64(nmax))
	gNew := float64(g) * (1 - contrast*float64(n)/float64(nmax))
	bNew := float64(b) * (1 - contrast*float64(n)/float64(nmax))
	return color.RGBA64{uint16(max(0, rNew)), uint16(max(0, gNew)), uint16(max(0, bNew)), 0xffff}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
