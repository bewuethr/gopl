// BigRat renders the Mandelbrot fractal using math/big.Rat.
package main

import (
	"image"
	"image/color"
	"image/png"
	"math/big"
	"os"
	"strconv"
)

const prec = 100

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
	const width, height = 200, 200
	var (
		xmin = big.NewRat(-2, 1)
		ymin = big.NewRat(-2, 1)
		xmax = big.NewRat(2, 1)
		ymax = big.NewRat(2, 1)
	)
	if len(os.Args) == 4 {
		xmin64, _ := strconv.ParseFloat(os.Args[1], 64)
		ymin64, _ := strconv.ParseFloat(os.Args[2], 64)
		xmax64, _ := strconv.ParseFloat(os.Args[3], 64)

		xmin.SetFloat64(xmin64)
		ymin.SetFloat64(ymin64)
		xmax.SetFloat64(xmax64)
		t := new(big.Rat)
		ymax.Add(ymin, t.Sub(xmax, xmin))
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y, ydiff := new(big.Rat), new(big.Rat)
		t1, t2, t3 := new(big.Rat), new(big.Rat), new(big.Rat)

		// y := float64(py)/float64(height)*(ymax-ymin) + ymin
		y.Add(t1.Mul(t2.Quo(new(big.Rat).SetInt64(int64(py)), new(big.Rat).SetInt64(height)), t3.Sub(ymax, ymin)), ymin)

		// ydiff := (ymax - ymin) / float64(4*height)
		ydiff.Quo(t1.Sub(ymax, ymin), t2.Mul(big.NewRat(4, 1), new(big.Rat).SetInt64(height)))

		for px := 0; px < width; px++ {
			x, xdiff := new(big.Rat), new(big.Rat)

			// x := float64(px)/float64(width)*(xmax-xmin) + xmin
			x.Add(t1.Mul(t2.Quo(new(big.Rat).SetInt64(int64(px)), new(big.Rat).SetInt64(width)), t3.Sub(xmax, xmin)), xmin)

			// xdiff := (xmax - xmin) / float64(4*width)
			xdiff.Quo(t1.Sub(xmax, xmin), t2.Mul(new(big.Rat).SetInt64(4), new(big.Rat).SetInt64(width)))

			var colVec []color.Color
			colVec = append(colVec, mandelbrot(t1.Sub(x, xdiff), t2.Sub(y, ydiff)))
			colVec = append(colVec, mandelbrot(t1.Sub(x, xdiff), t2.Add(y, ydiff)))
			colVec = append(colVec, mandelbrot(t1.Add(x, xdiff), t2.Sub(y, ydiff)))
			colVec = append(colVec, mandelbrot(t1.Add(x, xdiff), t2.Add(y, ydiff)))

			// Image point (px, py) represents complex value z.
			img.Set(px, py, average(colVec))
		}
	}
	png.Encode(os.Stdout, img) // NOTE: ignoring errors
}

func mandelbrot(zRe *big.Rat, zIm *big.Rat) color.Color {
	const iterations = 10
	const contrast = 7

	vRe, vIm := new(big.Rat), new(big.Rat)
	for n := uint8(0); n < iterations; n++ {
		t1, t2, t3 := new(big.Rat), new(big.Rat), new(big.Rat)

		// Real part
		vReNew := new(big.Rat)
		t1.Mul(vRe, vRe)
		t2.Mul(vIm, vIm)
		t3.Sub(t1, t2)
		vReNew.Add(t3, zRe)

		// Imaginary part
		vImNew := new(big.Rat)
		t1.Mul(vRe, vIm)
		t2.Mul(big.NewRat(2, 1), t1)
		vImNew.Add(t2, zIm)

		vRe.Set(vReNew)
		vIm.Set(vImNew)

		// Absolute value squared, hence comparing to 4 instead of 2
		t3.Add(t1.Mul(vRe, vRe), t2.Mul(vIm, vIm))

		if cmp := t3.Cmp(big.NewRat(4, 1)); cmp == 1 {
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
