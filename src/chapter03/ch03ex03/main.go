// Modify surface to colour peaks red and valleys blue
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
	min, max := getExtremes()
	diff := max - min
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>\n", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			isNumber := true
			ax, ay, z, ok := corner(i+1, j)
			isNumber = isNumber && ok
			bx, by, z, ok := corner(i, j)
			isNumber = isNumber && ok
			cx, cy, z, ok := corner(i, j+1)
			isNumber = isNumber && ok
			dx, dy, z, ok := corner(i+1, j+1)
			isNumber = isNumber && ok
			if !isNumber {
				continue
			}
			redComp := int((z - min) / diff * 0xff)
			blueComp := int((1 - (z - min)) / diff * 0xff)
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
			fmt.Printf("<polygon fill='#%02x00%02x' points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				redComp, blueComp, ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Println("</svg>")
}

func getExtremes() (float64, float64) {
	min, max := 0.0, 0.0
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			x := xyrange * (float64(i)/cells - 0.5)
			y := xyrange * (float64(j)/cells - 0.5)
			z, ok := f(x, y)
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

func corner(i, j int) (float64, float64, float64, bool) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z, ok := f(x, y)
	if !ok {
		return 0, 0, 0, false
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z, true
}

func f(x, y float64) (float64, bool) {
	r := math.Hypot(x, y) // distance from (0,0)
	res := math.Sin(r) / r
	if math.IsNaN(res) {
		return 0, false
	}
	return res, true
}
