package starmap

import (
	"geom"
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

func TestClosest(t *testing.T) {
	data, err := LoadData("../data/bright.tsv")
	if err != nil {
		t.Errorf("Can't load test data: %v", err)
	}
    //p := geom.NewPoint2D(14.579589, 25.817871)
    p := geom.NewPoint2D(14.7249, 26.5155)
    res := data.FindClosest(p)
    if res == nil || res.GeoHash != "ehdyym3b" {
        t.Errorf("no point found")
    } else if res.GeoHash != "ehdyym3b" {
        p, _ = geom.UnHash(res.GeoHash, geom.STELLAR)
        t.Errorf("Expected , got %v", p)
    }
}


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
