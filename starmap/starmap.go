package starmap

import (
	"image"
	"image/color"
	"image/png"
	"net/http"
	"render"
	"render/style"
)

func init() {
	/* handler() defined below */
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	/* simple test image output */
	p0 := &image.Point{300, 150}
	smallCircle := style.NewPointStyle(8, color.White, style.CIRCLE)
	p1 := &image.Point{300, 75}
	blue := color.RGBA{0, 0, 255, 255}
	smallSquare := style.NewPointStyle(8, blue, style.SQUARE)
	p3 := &image.Point{300, 225}
	p4 := &image.Point{150, 150}
	red := color.RGBA{255, 0, 0, 255}
	largeCircle := style.NewPointStyle(100, red, style.CIRCLE)
	p5 := &image.Point{450, 150}
	green := color.RGBA{0, 255, 0, 255}
	largeSquare := style.NewPointStyle(100, green, style.SQUARE)
	img := render.Create(600, 300, color.Black)
	render.Render(img, p0, smallCircle)
	render.Render(img, p1, smallSquare)
	render.Render(img, p3, smallSquare)
	render.Render(img, p4, largeCircle)
	render.Render(img, p5, largeSquare)
	if err := png.Encode(w, img); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
