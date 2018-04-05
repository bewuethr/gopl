// Web server that computes sufaces and writes SVG data to client
package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	cells   = 100         // number of grid cells
	xyrange = 30.0        // axis ranges (-xyrange..+xyrange)
	angle   = math.Pi / 6 // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		surface(w, r)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}

func surface(w http.ResponseWriter, r *http.Request) {
	width, height := 600, 320 // canvas size in pixels
	useColor := false
	var surf string

	if r != nil {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}
		if v, ok := r.Form["w"]; ok {
			width, _ = strconv.Atoi(v[0]) // NOTE: ignoring conversion errors
		}
		if v, ok := r.Form["h"]; ok {
			height, _ = strconv.Atoi(v[0]) // NOTE: ignoring conversion errors
		}
		if _, ok := r.Form["col"]; ok {
			useColor = true
		}
		if v, ok := r.Form["surf"]; ok {
			surf = v[0]
		}
	}

	xyscale := float64(width) / 2.0 / xyrange // pixels per x or y unit
	zscale := float64(height) * 0.4           // pixels per z unit

	w.Header().Set("Content-Type", "image/svg+xml")

	var min, max, diff float64
	if useColor {
		min, max = getExtremes(surf)
		diff = max - min
	}

	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>\n", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			isNumber := true
			ax, ay, z, ok := corner(i+1, j, width, height, xyscale, zscale, surf)
			isNumber = isNumber && ok
			bx, by, z, ok := corner(i, j, width, height, xyscale, zscale, surf)
			isNumber = isNumber && ok
			cx, cy, z, ok := corner(i, j+1, width, height, xyscale, zscale, surf)
			isNumber = isNumber && ok
			dx, dy, z, ok := corner(i+1, j+1, width, height, xyscale, zscale, surf)
			isNumber = isNumber && ok
			if !isNumber {
				continue
			}
			colString := ""
			if useColor {
				redComp := int((z - min) / diff * 0xff)
				blueComp := int((1 - (z-min)/diff) * 0xff)
				if redComp > 0xff {
					redComp = 0xff
				}
				if blueComp > 0xff {
					blueComp = 0xff
				}
				if redComp < 0 {
					redComp = 0
				}
				if blueComp < 0 {
					blueComp = 0
				}
				colString = fmt.Sprintf("fill='#%02x%02x%02x' ", redComp, 0, blueComp)
			}
			fmt.Fprintf(w, "<polygon %spoints='%g,%g %g,%g %g,%g %g,%g'/>\n",
				colString, ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintln(w, "</svg>")
}

func getExtremes(surf string) (float64, float64) {
	min, max := 0.0, 0.0
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			x := xyrange * (float64(i)/cells - 0.5)
			y := xyrange * (float64(j)/cells - 0.5)
			z, ok := f(x, y, surf)
			if !ok {
				continue
			}
			if z > max {
				max = z
			}
			if z < min {
				min = z
			}
		}
	}
	return min, max
}

func corner(i, j, w, h int, xyscale, zscale float64, surf string) (float64, float64, float64, bool) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z, ok := f(x, y, surf)
	if !ok {
		return 0, 0, 0, false
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := float64(w/2) + (x-y)*cos30*xyscale
	sy := float64(h/2) + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z, true
}

func f(x, y float64, surf string) (float64, bool) {
	res := 0.0

	switch surf {
	case "eggbox":
		res = eggBox(x, y)
	case "moguls":
		res = moguls(x, y)
	case "saddle":
		res = saddle(x, y)
	default:
		r := math.Hypot(x, y) // distance from (0,0)
		res = math.Sin(r) / r
	}

	if math.IsNaN(res) {
		return 0.0, false
	}
	return res, true
}

func eggBox(x, y float64) float64 {
	return (math.Sin(x) + math.Sin(y)) / 10
}

func moguls(x, y float64) float64 {
	return math.Pow(math.Sin(x/xyrange*3*math.Pi), 2) * math.Cos(y/xyrange*3*math.Pi)
}

func saddle(x, y float64) float64 {
	return math.Pow(y/xyrange*2, 2) - math.Pow(x/xyrange*2, 2)
}
