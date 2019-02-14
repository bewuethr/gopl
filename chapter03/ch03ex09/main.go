// Ch03ex09 implements a web server that renders fractals and writes image data
// to client. The client can specify x, y, size and zoom values in the request.
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"net/http"
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
	handler := func(w http.ResponseWriter, r *http.Request) {
		mandelbrotHandler(w, r)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}

func mandelbrotHandler(w http.ResponseWriter, r *http.Request) {
	var (
		x, y = 0.0, 0.0
		size = 950
		zoom = 0
	)

	if r != nil {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}
		if v, ok := r.Form["x"]; ok {
			x, _ = strconv.ParseFloat(v[0], 64) // NOTE: ignoring conversion errors
		}
		if v, ok := r.Form["y"]; ok {
			y, _ = strconv.ParseFloat(v[0], 64) // NOTE: ignoring conversion errors
		}
		if v, ok := r.Form["size"]; ok {
			size, _ = strconv.Atoi(v[0]) // NOTE: ignoring conversion errors
		}
		if v, ok := r.Form["zoom"]; ok {
			zoom, _ = strconv.Atoi(v[0]) // NOTE: ignoring conversion errors
		}
	}

	axisLength := 2.0 / float64(uint(1)<<uint(zoom))
	xmin := x - axisLength
	xmax := x + axisLength
	ymin := y - axisLength
	ymax := y + axisLength

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for py := 0; py < size; py++ {
		y := float64(py)/float64(size)*(ymax-ymin) + ymin
		ydiff := (ymax - ymin) / (4 * float64(size))
		for px := 0; px < size; px++ {
			x := float64(px)/float64(size)*(xmax-xmin) + xmin
			xdiff := (xmax - xmin) / (4 * float64(size))

			var colVec []color.Color
			colVec = append(colVec, mandelbrot(complex(x-xdiff, y-ydiff)))
			colVec = append(colVec, mandelbrot(complex(x-xdiff, y+ydiff)))
			colVec = append(colVec, mandelbrot(complex(x+xdiff, y-ydiff)))
			colVec = append(colVec, mandelbrot(complex(x+xdiff, y+ydiff)))

			// Image point (px, py) represents complex value z.
			img.Set(px, py, average(colVec))
		}
	}
	png.Encode(w, img) // NOTE: ignoring errors
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
