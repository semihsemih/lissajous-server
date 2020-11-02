package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	cycles, ok := r.URL.Query()["cycles"]
	if !ok || len(cycles[0]) < 1 {
		cycle := 1
		lissajous(w, float64(cycle))
		return
	}
	cycle, _ := strconv.Atoi(cycles[0])
	lissajous(w, float64(cycle))
}

var palette = []color.Color{
	color.Black,
	color.RGBA{R: 6, G: 249, B: 132, A: 1},
	color.RGBA{R: 253, G: 232, B: 2, A: 1},
	color.RGBA{R: 255, G: 145, B: 26, A: 1},
	color.RGBA{R: 252, G: 93, B: 2, A: 1},
	color.RGBA{R: 255, G: 0, B: 249, A: 1},
}

func lissajous(out io.Writer, cycles float64) {
	const (
		res     = 0.0003 // angular resolution
		size    = 300    // image canvas covers [-size..+size]
		nframes = 512    // number of animation frames
		delay   = 8      // delay between frames in 10 ms units
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
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8((i%(len(palette)-1))+1))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
