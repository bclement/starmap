package geom

import (
	"math"
    "image"
	"strings"
	"testing"
)

func TestBBoxPoint(t *testing.T) {
	p1 := NewPoint2D(1, 2)
	p2 := NewPoint3D(2, 3, 42)
	p3 := NewPoint2D(3, 4)
	bbox := NewBBox3D(-1, -1, 42, 2, 3, 43)
	if !bbox.Contains(p1) || !bbox.Covers(p1) {
		t.Errorf("Should contain and cover")
	}
	if bbox.Contains(p2) {
		t.Errorf("Shouldn't contain")
	}
	if !bbox.Covers(p2) {
		t.Errorf("Should cover")
	}
	if bbox.Contains(p3) || bbox.Covers(p3) {
		t.Errorf("Shouldn't contain or cover")
	}
}

func TestBBoxPoly(t *testing.T) {
	p1, _ := NewPoly2D(0, 0, 0, 1, 1, 1, 1, 0)
	p2, _ := NewPoly3D(-1, -1, -1, -1, -2, -1, -2, -2, -1)
	p3, _ := NewPoly2D(3, 3, 3, 4, 4, 4, 4, 3)
	bbox := NewBBox3D(-2, -2, -1, 2, 2, 2)
	if !bbox.Contains(p1) || !bbox.Covers(p1) {
		t.Errorf("Should contain and cover")
	}
	if bbox.Contains(p2) {
		t.Errorf("Shouldn't contain")
	}
	if !bbox.Covers(p2) {
		t.Errorf("Should cover")
	}
	if bbox.Contains(p3) || bbox.Covers(p3) {
		t.Errorf("Shouldn't contain or cover")
	}
}

func TestGeoHash(t *testing.T) {
	/* example from http://en.wikipedia.org/wiki/Geohash */
	p := NewPoint2D(-5.6, 42.6)
	res := p.GeoHash(LONLAT)
	if !strings.HasPrefix(res, "ezs42") {
		t.Error("Didn't expect " + res)
	}
	latMargin := 0.000085
	lonMargin := 0.00017
	newP, err := UnHash(res, LONLAT)
	if err != nil {
		t.Error("Invalid geohash")
	}
	if math.Abs(newP.X()-(-5.6)) > lonMargin {
		t.Errorf("Inavlid lon value: %v", newP.X())
	}
	if math.Abs(newP.Y()-42.6) > latMargin {
		t.Errorf("Inavlid lat value: %v", newP.Y())
	}
}

func TestHashBBox(t *testing.T) {
	lower := NewPoint2D(-180, -90)
	upper0 := NewPoint2D(180, 90)
	assertBbox(lower, upper0, "0", "~", t)
	upper1 := NewPoint2D(0, 0)
	assertBbox(lower, upper1, "0", "8", t)
	upper2 := NewPoint2D(-90, 0)
	assertBbox(lower, upper2, "0", "4", t)
	upper3 := NewPoint2D(-135, -45)
	assertBbox(lower, upper3, "00", "0~", t)
	upper4 := NewPoint2D(-180, -90)
	assertBbox(lower, upper4, "00000000", "00000000", t)
}

func assertBbox(lower, upper *Point, expMin, expMax string, t *testing.T) {
	min, max := BBoxHash(lower, upper, LONLAT)
	if min != expMin || max != expMax {
		t.Errorf("Expected %v, %v, got %v %v", expMin, expMax, min, max)
	}
}

func TestTransform(t *testing.T) {
    min := NewPoint2D(-90, -30)
    max := NewPoint2D(90, 30)
    trans := CreateTransform(min, max, 90, 30, LONLAT)
    p0 := NewPoint2D(-90, -30)
    assertTrans(trans.Transform(p0), &image.Point{0, 30}, t)
    p1 := NewPoint2D(-45, -15)
    assertTrans(trans.Transform(p1), &image.Point{22, 22}, t)
    p2 := NewPoint2D(0, 0)
    assertTrans(trans.Transform(p2), &image.Point{45,15}, t)
    p3 := NewPoint2D(45, 15)
    assertTrans(trans.Transform(p3), &image.Point{67, 7}, t)
    p4 := NewPoint2D(90, 30)
    assertTrans(trans.Transform(p4), &image.Point{90, 0}, t)
}

func assertTrans(res, exp *image.Point, t *testing.T) {
    if res.X != exp.X || res.Y != exp.Y {
        t.Errorf("Expected %v, %v. Got %v, %v", exp.X, exp.Y, res.X, res.Y)
    }
}
