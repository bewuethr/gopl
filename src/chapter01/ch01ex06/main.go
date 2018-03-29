// Produce Lissajous images in multiple colours
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

// Use Solarized Dark palette
var palette = []color.Color{
	color.RGBA{0x00, 0x2b, 0x36, 0xff}, // base03
	color.RGBA{0xb5, 0x89, 0x00, 0xff}, // yellow
	color.RGBA{0xcb, 0x4b, 0x16, 0xff}, // orange
	color.RGBA{0xdc, 0x32, 0x2f, 0xff}, // red
	color.RGBA{0xd3, 0x36, 0x82, 0xff}, // magenta
	color.RGBA{0x6c, 0x71, 0xc4, 0xff}, // violet
	color.RGBA{0x26, 0x8b, 0xd2, 0xff}, // blue
	color.RGBA{0x2a, 0xa1, 0x98, 0xff}, // cyan
	color.RGBA{0x85, 0x99, 0x00, 0xff}, // green
}

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				uint8(t/(cycles*2*math.Pi)*8)+1)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
