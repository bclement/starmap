package render

import (
	"github.com/bclement/starmap/render/style"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
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
	img := Create(600, 300, color.Black)
	Render(img, p0, smallCircle)
	Render(img, p1, smallSquare)
	Render(img, p3, smallSquare)
	Render(img, p4, largeCircle)
	Render(img, p5, largeSquare)
	f, err := os.OpenFile("/tmp/res.png", os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		t.Errorf("open %s", err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		t.Error("encode %s", err)
	}
}
