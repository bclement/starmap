package starmap

import (
	"geom"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"render"
	"render/style"
	"testing"
)

func TestPrefix(t *testing.T) {
	if sharedLen("", "") != 0 {
		t.Errorf("Broken for empty strings")
	}
	if sharedLen("a", "b") != 0 {
		t.Errorf("Broken for no shared prefix")
	}
	if sharedLen("aa", "a") != 1 {
		t.Errorf("Broken for single common")
	}
	if sharedLen("thecake", "thecakeisalie") != 7 {
		t.Errorf("Broken for shared prefix")
	}
	if sharedLen("aaa", "aaa") != 3 {
		t.Errorf("Broken for match")
	}
}

/*
func TestClosest(t *testing.T) {
	data, err := LoadData("../data/bright.tsv")
	if err != nil {
		t.Errorf("Can't load test data: %v", err)
	}
	//p := geom.NewPoint2D(14.579589, 25.817871)
	p := geom.NewPoint2D(14.7249, 26.5155)
	res := FindClosest(p)
	if res == nil || res.GeoHash != "ehdyym3b" {
		t.Errorf("no point found")
	} else if res.GeoHash != "ehdyym3b" {
		p, _ = geom.UnHash(res.GeoHash, geom.STELLAR)
		t.Errorf("Expected , got %v", p)
	}
}*/

func TestFilter(t *testing.T) {
	lower := geom.NewPoint2D(24, -90)
	upper := geom.NewPoint2D(22.5, -67.5)
	assertFilter(t, lower, upper, 10, "00", "08")
}

func assertFilter(t *testing.T, lower, upper *geom.Point, num int,
	expLowHash, expUpHash string) {
	data, err := LoadData("../data/bright.tsv")
	if err != nil {
		t.Errorf("Can't load test data: %v", err)
	}
	width := 256
	height := 256
	trans := geom.CreateTransform(lower, upper, width, height, geom.STELLAR)
	lowerHash, upperHash := geom.BBoxHash(lower, upper, geom.STELLAR)
	if lowerHash != expLowHash || upperHash != expUpHash {
		t.Errorf("expected %v-%v, got %v-%v", expLowHash, expUpHash,
			lowerHash, upperHash)
	}
	stars := data.Range(lowerHash, upperHash)
	if len(stars) != num {
		t.Errorf("expected %v stars, got", num, len(stars))
	}
	for _, s := range stars {
		coord, err := geom.UnHash(s.GeoHash, geom.STELLAR)
		if err != nil {
			t.Errorf("problem unhashing: %v", err)
		}
		pix := trans.Transform(coord)
		if pix.X < 0 || pix.X >= width || pix.Y < 0 || pix.Y >= height {
			t.Errorf("problem transforming %v, %v", pix, coord)
		}
	}
}

func TestReadPoly(t *testing.T) {
	p, err := readWktFile("../data/consts/Equuleus.wkt")
	if err != nil {
		t.Errorf("cant read: %v", err)
	}
	coords := p.Coords()
	if coords.Len() != 39 {
		t.Errorf("expected %v coords, got %v", 39, coords.Len())
	}
}

func TestConstFilter(t *testing.T) {
	data, err := LoadConstellations("../data/consts")
	if err != nil || len(data) < 1 {
		t.Errorf("loading: %v", err)
	}
	s := style.NewPolyStyle(1, color.White)
	width := 512
	height := 256
	lower := geom.NewPoint2D(24, -90)
	upper := geom.NewPoint2D(0, 90)
	trans := geom.CreateTransform(lower, upper, width, height, geom.STELLAR)
	for _, c := range data {
		if len(c.PolyInfos) < 1 {
			t.Errorf("Expected poly infos for %v", c.Name)
		}
		img := render.Create(width, height, color.Black)
		for _, pi := range c.PolyInfos {
			if len(pi.LabelPoint) != 2 {
				t.Errorf("Expected label points for %v", c.Name)
			}
			render.RenderPoly(img, pi.Geom, trans, s)
		}
		writeImg(t, img, "/tmp/nexconst/"+c.Name+".png")
	}
}

func writeImg(t *testing.T, img draw.Image, fname string) {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		t.Errorf("open %s", err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		t.Errorf("encode %v", err)
	}
}

func TestReadString(t *testing.T) {
	seqs, err := readStringsWktFile("../data/consts/Orion-strings.wkt")
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Errorf("%v", seqs)
}
