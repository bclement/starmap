package geom

import (
	"image"
	"math"
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
	testRoundTrip(t, p, "ezs42", LONLAT)
	pLonLat := NewPoint2D(-180, -90)
	pStellar := NewPoint2D(24, -90)
	hash1 := pLonLat.GeoHash(LONLAT)
	hash2 := pStellar.GeoHash(STELLAR)
	if hash1 != hash2 {
		t.Errorf("expected equal %v != %v", hash1, hash2)
	}
	newP, err := UnHash(hash2, STELLAR)
	if err != nil {
		t.Errorf("invalid: %v", err)
	}
	assertPoint(t, newP, 24, -90)
	p = NewPoint2D(18, 0)
	testRoundTrip(t, p, "d00", STELLAR)
	p = NewPoint2D(0, 0)
	testRoundTrip(t, p, "s00000", LONLAT)
	p = NewPoint2D(12, 0)
	testRoundTrip(t, p, "s00000", STELLAR)
}

func testRoundTrip(t *testing.T, p *Point, prefix string, gd *GridDef) {
	res := p.GeoHash(gd)
	if !strings.HasPrefix(res, prefix) {
		t.Errorf("Expected prefix %v, got %v", prefix, res)
	}
	newP, err := UnHash(res, gd)
	if err != nil {
		t.Errorf("Invalid geohash, %v", err)
	}
	assertPoint(t, newP, p.X(), p.Y())
}

func assertPoint(t *testing.T, p *Point, expX, expY float64) {
	yMargin := 0.0085
	xMargin := 0.0017
	if math.Abs(p.X()-expX) > xMargin {
		t.Errorf("expected x to be %v, got %v", expX, p.X())
	}
	if math.Abs(p.Y()-expY) > yMargin {
		t.Errorf("expected y to be %v, got %v", expY, p.Y())
	}
}

func TestHashBBox(t *testing.T) {
	lower := NewPoint2D(-180, -90)
	upper0 := NewPoint2D(180, 90)
	assertBbox(LONLAT, lower, upper0, "0", "~", t)
	upper1 := NewPoint2D(0, 0)
	assertBbox(LONLAT, lower, upper1, "0", "8", t)
	upper2 := NewPoint2D(-90, 0)
	assertBbox(LONLAT, lower, upper2, "0", "4", t)
	upper3 := NewPoint2D(-135, -45)
	assertBbox(LONLAT, lower, upper3, "00", "0~", t)
	upper4 := NewPoint2D(-180, -90)
	assertBbox(LONLAT, lower, upper4, "00000000", "00000000", t)
	lower = NewPoint2D(24, -90)
	upper := NewPoint2D(22.5, -67.5)
	assertBbox(STELLAR, lower, upper, "00", "08", t)
	upper = NewPoint2D(0, 90)
	assertBbox(STELLAR, lower, upper, "0", "~", t)
	upper = NewPoint2D(12, 0)
	assertBbox(STELLAR, lower, upper, "0", "8", t)
	lower = NewPoint2D(18, 0)
	upper = NewPoint2D(15, 45)
	assertBbox(STELLAR, lower, upper, "d0", "d~", t)
	lower = NewPoint2D(13.5, 0)
	upper = NewPoint2D(12, 22.5)
	assertBbox(STELLAR, lower, upper, "e8", "eh", t)
	lower = NewPoint2D(-22.5, 0)
	upper = NewPoint2D(0, 22.5)
	assertBbox(LONLAT, lower, upper, "e8", "eh", t)
}

func assertBbox(gd *GridDef, lower, upper *Point, expMin, expMax string,
	t *testing.T) {
	min, max := BBoxHash(lower, upper, gd)
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
	assertTrans(trans.Transform(p2), &image.Point{45, 15}, t)
	p3 := NewPoint2D(45, 15)
	assertTrans(trans.Transform(p3), &image.Point{67, 7}, t)
	p4 := NewPoint2D(90, 30)
	assertTrans(trans.Transform(p4), &image.Point{90, 0}, t)
}

func TestTransStellar(t *testing.T) {
	min := NewPoint2D(18, -30)
	max := NewPoint2D(6, 30)
	trans := CreateTransform(min, max, 90, 30, STELLAR)
	p0 := NewPoint2D(18, -30)
	assertTrans(trans.Transform(p0), &image.Point{0, 30}, t)
	p1 := NewPoint2D(15, -15)
	assertTrans(trans.Transform(p1), &image.Point{22, 22}, t)
	p2 := NewPoint2D(12, 0)
	assertTrans(trans.Transform(p2), &image.Point{45, 15}, t)
	p3 := NewPoint2D(9, 15)
	assertTrans(trans.Transform(p3), &image.Point{67, 7}, t)
	p4 := NewPoint2D(6, 30)
	assertTrans(trans.Transform(p4), &image.Point{90, 0}, t)
}

func assertTrans(res, exp *image.Point, t *testing.T) {
	if res.X != exp.X || res.Y != exp.Y {
		t.Errorf("Expected %v, %v. Got %v, %v", exp.X, exp.Y, res.X, res.Y)
	}
}
