// Modify Lissajous server to read parameter values from URL
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func main() {
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "web" {
		handler := func(w http.ResponseWriter, r *http.Request) {
			lissajous(w, r)
		}
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}
	lissajous(os.Stdout, nil)
}

func lissajous(out io.Writer, r *http.Request) {
	var (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)

	if r != nil {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}
		if v, ok := r.Form["cycles"]; ok {
			cycles, _ = strconv.Atoi(v[0]) // NOTE: ignoring conversion errors
		}
		if v, ok := r.Form["res"]; ok {
			res, _ = strconv.ParseFloat(v[0], 64) // NOTE: ignoring conversion errors
		}
		if v, ok := r.Form["size"]; ok {
			size, _ = strconv.Atoi(v[0]) // NOTE: ignoring conversion errors
		}
		if v, ok := r.Form["nframes"]; ok {
			nframes, _ = strconv.Atoi(v[0]) // NOTE: ignoring conversion errors
		}
		if v, ok := r.Form["delay"]; ok {
			delay, _ = strconv.Atoi(v[0]) // NOTE: ignoring conversion errors
		}
	}
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
