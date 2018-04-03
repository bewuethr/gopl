// Modify surface to produce different surfaces
package main

import (
	"fmt"
	"math"
	"os"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {

	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>\n", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			isNumber := true
			ax, ay, ok := corner(i+1, j)
			isNumber = isNumber && ok
			bx, by, ok := corner(i, j)
			isNumber = isNumber && ok
			cx, cy, ok := corner(i, j+1)
			isNumber = isNumber && ok
			dx, dy, ok := corner(i+1, j+1)
			isNumber = isNumber && ok
			if !isNumber {
				continue
			}
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int) (float64, float64, bool) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	var (
		arg string
		z   float64
		ok  bool
	)
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	switch arg {
	case "eggbox":
		z, ok = eggBox(x, y)
	case "moguls":
		z, ok = moguls(x, y)
	case "saddle":
		z, ok = saddle(x, y)
	default:
		z, ok = f(x, y)
	}
	if !ok {
		return 0, 0, false
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, true
}

func eggBox(x, y float64) (float64, bool) {
	res := (math.Sin(x) + math.Sin(y)) / 10
	if math.IsNaN(res) {
		return 0, false
	}
	return res, true
}

func moguls(x, y float64) (float64, bool) {
	res := math.Pow(math.Sin(x/xyrange*3*math.Pi), 2) * math.Cos(y/xyrange*3*math.Pi)
	if math.IsNaN(res) {
		return 0, false
	}
	return res, true
}

func saddle(x, y float64) (float64, bool) {
	res := math.Pow(y/xyrange*2, 2) - math.Pow(x/xyrange*2, 2)
	if math.IsNaN(res) {
		return 0, false
	}
	return res, true
}

func f(x, y float64) (float64, bool) {
	r := math.Hypot(x, y) // distance from (0,0)
	res := math.Sin(r) / r
	if math.IsNaN(res) {
		return 0, false
	}
	return res, true
}
