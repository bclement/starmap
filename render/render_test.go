package render

import (
	"bufio"
	"geom"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"render/style"
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
	writeImg(t, img, "/tmp/res.png")
}

func writeImg(t *testing.T, img draw.Image, fname string) {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		t.Errorf("open %s", err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		t.Error("encode %s", err)
	}
}

func TestCircle(t *testing.T) {
	img := Create(256, 256, color.Black)
	for y := 0; y < 16; y += 1 {
		for i := 0; i < 8; i += 1 {
			size := float64(i+1) / 2.0
			alpha := uint8((16 * (y + 1)) - 1)
			color := color.RGBA{alpha, alpha, alpha, 255}
			circle := style.NewPointStyle(size, color, style.CIRCLE)
			p := &image.Point{(25 * i) + 10, int(alpha)}
			Render(img, p, circle)
		}
	}
	writeImg(t, img, "/tmp/dots.png")
}

func TestLines(t *testing.T) {
	img := Create(256, 256, color.Black)
	s := style.NewPolyStyle(1, color.White)
	p0 := &image.Point{0, 0}
	for x := 0; x < 256; x += 16 {
		p1 := &image.Point{x, 255}
		RenderLine(img, p0, p1, s)
	}
	writeImg(t, img, "/tmp/lines.png")
}

func TestPolys(t *testing.T) {
	img := Create(256, 256, color.Black)
	s := style.NewPolyStyle(1, color.White)
	p1, err1 := geom.NewPoly2D(16.5, 67.5, 16.5, 22.5, 13.5, 22.5, 13.5, 67.5)
	p2, err2 := geom.NewPoly2D(14, 22.5, 13.5, 11, 13, 22.5, 13.5, 33)
	if err1 != nil || err2 != nil {
		t.Errorf("can't create poly %v %v", err1, err1)
	}
	min := geom.NewPoint2D(15, 0)
	max := geom.NewPoint2D(12, 45)
	trans := geom.CreateTransform(min, max, 255, 255, geom.STELLAR)
	RenderPoly(img, p1, trans, s)
	RenderPoly(img, p2, trans, s)
	writeImg(t, img, "/tmp/poly.png")
}

func TestString(t *testing.T) {
	img := Create(256, 256, color.Black)
	f, err := os.Open("../data/chars.png")
	if err != nil {
		t.Errorf("can't open chars: %v", err)
	}
	defer f.Close()
	chars, _, err := image.Decode(bufio.NewReader(f))
	if err != nil {
		t.Errorf("can't decode chars: %v", err)
	}
	p := image.Pt(20, 20)
	c := color.RGBA{255, 0, 0, 255}
	RenderString(img, chars, 10, &p, "Hello World", c)
	writeImg(t, img, "/tmp/strings.png")
}
