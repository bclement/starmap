package geom

import "testing"

func TestBBoxPoint(t *testing.T) {
	p1 := NewPoint2D(1, 2)
	p2 := NewPoint3D(2, 3, 42)
	p3 := NewPoint2D(3, 4)
	bbox := NewBBox3D(-1,-1,42, 2,3,43)
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

func TestBBoxPoly(t *testing.T){
    p1,_ := NewPoly2D(0,0, 0,1, 1,1, 1,0)
    p2,_ := NewPoly3D(-1,-1,-1, -1,-2,-1, -2,-2,-1)
    p3,_ := NewPoly2D(3,3, 3,4, 4,4, 4,3)
    bbox := NewBBox3D(-2,-2,-1, 2,2,2)
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
