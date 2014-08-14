package geom

import (
	"fmt"
	"math"
)

/* generic coordinate storage */
type CoordinateSeq struct {
	Coords []float64
	Dims   int
}

/* returns number of coodinates in sequence */
func (cs *CoordinateSeq) Len() int {
	return len(cs.Coords) / cs.Dims
}

/* returns coordinate at index */
func (cs *CoordinateSeq) Get(index int) []float64 {
	start := index * cs.Dims
	return cs.Coords[start : start+cs.Dims]
}

/* interface for interacting with coordinate based geometries */
type Geometry interface {
	/* return a pointer to the coordinate storage for geometry */
	Coords() *CoordinateSeq
	/* return number of dimensions per coordinate */
	Dims() int
}

/* single coordinate point */
type Point struct {
	/* storage for single coordinate */
	c []float64
}

/* returns a pointer to a newly created 2D point */
func NewPoint2D(x float64, y float64) *Point {
	return &Point{[]float64{x, y}}
}

/* returns a pointer to a newly created 3D point */
func NewPoint3D(x float64, y float64, z float64) *Point {
	return &Point{[]float64{x, y, z}}
}

/* see Geometry interface */
func (p *Point) Coords() *CoordinateSeq {
	return &CoordinateSeq{p.c, len(p.c)}
}

/* see Geometry interface */
func (p *Point) Dims() int {
	return len(p.c)
}

/* access first dimension value */
func (p *Point) X() float64 {
	return p.c[0]
}

/* access second dimension value */
func (p *Point) Y() float64 {
	return p.c[1]
}

/* access third dimension value */
func (p *Point) Z() float64 {
	if len(p.c) < 3 {
		return 0
	} else {
		return p.c[2]
	}
}

/* see Stringer interface */
func (p *Point) String() string {
	return fmt.Sprintf("%v", p.c)
}

/* multi-point polygon */
type Polygon struct {
	/* storage for exterior ring */
	c    CoordinateSeq
	bbox *BoundingBox
}

/* create a new polygon with given dimensions from coordinate sequence
polygon must be closed */
func NewPoly(dims int, coords ...float64) (*Polygon, error) {
	if dims < 2 || len(coords)%dims != 0 {
		return nil, fmt.Errorf("Invalid dimensions: %s", len(coords))
	}
	var bbox *BoundingBox = nil
	if dims == 2 {
		bbox = computeBbox2D(coords)
	}
	return &Polygon{CoordinateSeq{coords, dims}, bbox}, nil
}

/* take in a 2D coordinate sequence
return bounding box containing mins and maxes */
func computeBbox2D(coords []float64) *BoundingBox {
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for i := 1; i < len(coords); i += 2 {
		x, y := coords[i-1], coords[i]
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	return NewBBox2D(minX, minY, maxX, maxY)
}

/*
takes in coordinate slice packed as [x1,y1,x2,y2...]
returns newly created 2D polygon
error if length of coords isn't divisible by 2
*/
func NewPoly2D(coords ...float64) (*Polygon, error) {
	return NewPoly(2, coords...)
}

/*
takes in coordinate slice packed as [x1,y1,z1,x2,y2,z2...]
returns newly created 3D polygon
error if length coords insn't divisible by 3
*/
func NewPoly3D(coords ...float64) (*Polygon, error) {
	return NewPoly(3, coords...)
}

/* see Geometry interface */
func (p *Polygon) Dims() int {
	return p.c.Dims
}

/* see Geometry interface */
func (p *Polygon) Coords() *CoordinateSeq {
	return &p.c
}

/* return true if a is between the b's */
func between(a, b0, b1 float64) bool {
	return (b0 > a) != (b1 > a)
}

/* takes in a point (x,y) and a line defined by (x0,y0)(x1,y1)
returns 0 if point is on the line, negative if point is right of
the line and positive if point is left of the line (assuming x axis
increases right) */
func comp(x, y, x0, y0, x1, y1 float64) float64 {
	return (x1-x0)*(y-y0) - (x-x0)*(y1-y0)
}

/* winding number method
returns false positives for unknown reason
should return true if point in inside polygon
*/
func (p *Polygon) WindingContains(point *Point) bool {
	if p.bbox != nil && !p.bbox.Contains(point) {
		return false
	}
	x := point.X()
	y := point.Y()
	winding := 0
	length := p.c.Len()
	for i := 1; i < length; i += 1 {
		v0 := p.c.Get(i - 1)
		x0, y0 := v0[0], v0[1]
		v1 := p.c.Get(i)
		x1, y1 := v1[0], v1[1]
		if y0 <= y {
			if y1 > y && comp(x, y, x0, y0, x1, y1) > 0 {
				winding += 1
			}
		} else if y1 <= y && comp(x, y, x0, y0, x1, y1) < 0 {
			winding -= 1
		}
	}
	return winding != 0
}

/* return true if point x,y is left of the line between v0 and v1 */
func leftOf(x, y, x0, y0, x1, y1 float64) bool {
	lineX := (x1-x0)*(y-y0)/(y1-y0) + x0
	return x < lineX
}

/* ray casting method
returns true if point in inside polygon
*/
func (p *Polygon) Contains(point *Point) bool {
	if p.bbox != nil && !p.bbox.Contains(point) {
		return false
	}
	x := point.X()
	y := point.Y()
	rval := false
	length := p.c.Len()
	for i := 1; i < length; i += 1 {
		v0 := p.c.Get(i - 1)
		x0, y0 := v0[0], v0[1]
		v1 := p.c.Get(i)
		x1, y1 := v1[0], v1[1]
		if between(y, y0, y1) && leftOf(x, y, x0, y0, x1, y1) {
			rval = !rval
		}
	}
	return rval
}

/* bounds defined by two points */
type BoundingBox struct {
	/* lower bounds point */
	min []float64
	/* upper bounds point */
	max []float64
}

/*
returns pointer to newly created bounding box
uses mins and maxes of x and y
*/
func NewBBox2D(x0, y0, x1, y1 float64) *BoundingBox {
	min := []float64{math.Min(x0, x1), math.Min(y0, y1)}
	max := []float64{math.Max(x0, x1), math.Max(y0, y1)}
	return &BoundingBox{min, max}
}

/*
returns pointer to newly created bounding box
uses mins and maxes of x, y and z
*/
func NewBBox3D(x0, y0, z0, x1, y1, z1 float64) *BoundingBox {
	min := []float64{math.Min(x0, x1), math.Min(y0, y1), math.Min(z0, z1)}
	max := []float64{math.Max(x0, x1), math.Max(y0, y1), math.Max(z0, z1)}
	return &BoundingBox{min, max}
}

/* return true if other bounding box has same points as this */
func (bb *BoundingBox) Equals(other *BoundingBox) bool {
	return equals(bb.min, other.min) && equals(bb.max, other.max)
}

/* return true if a is greater than b */
func gt(a, b float64) bool {
	return a > b
}

/* return true if a is greater than or equal to b */
func gte(a, b float64) bool {
	return a >= b
}

/*
return true if this bounding box contains all points in geometry
does not return true if point lies on bounds, see Covers()
*/
func (bb *BoundingBox) Contains(g Geometry) bool {
	return checkAll(bb, g, gt)
}

/*
return true if this bounding box covers all points in geometry
includes points that lie on bounds
*/
func (bb *BoundingBox) Covers(g Geometry) bool {
	return checkAll(bb, g, gte)
}

/* return true if thsi bounding box includes any points in g */
func (bb *BoundingBox) Touches(g Geometry) bool {
	return checkAny(bb, g, gte)
}

/*
takes in bounds, geometry and a comparison function
return true if bounds has enough dimensions to cover geometry and
all geometry coordinates return true for comp(coord, min) and
true for comp(max, coord) for bounds min and max
*/
func checkAll(bb *BoundingBox, g Geometry,
	comp func(float64, float64) bool) bool {
	bbdims := len(bb.min)
	gdims := g.Dims()
	gc := g.Coords()
	size := gc.Len()
	if gdims <= bbdims {
		for i := 0; i < size; i += 1 {
			coord := gc.Get(i)
			for d := 0; d < bbdims && d < gdims; d += 1 {
				if !comp(coord[d], bb.min[d]) {
					return false
				}
				if !comp(bb.max[d], coord[d]) {
					return false
				}
			}
		}
	}
	return true
}

/*
takes in bounds, geometry and a comparison function
return true if bounds has enough dimensions to cover geometry and
any geometry coordinates return true for comp(coord, min) and
true for comp(max, coord) for bounds min and max
*/
func checkAny(bb *BoundingBox, g Geometry,
	comp func(float64, float64) bool) bool {
	bbdims := len(bb.min)
	gdims := g.Dims()
	gc := g.Coords()
	size := gc.Len()
	if gdims <= bbdims {
		for i := 0; i < size; i += 1 {
			coord := gc.Get(i)
			passes := true
			for d := 0; d < bbdims && d < gdims; d += 1 {
				if !comp(coord[d], bb.min[d]) {
					passes = false
					break
				}
				if !comp(bb.max[d], coord[d]) {
					passes = false
					break
				}
			}
			if passes {
				return true
			}
		}
	}
	return false
}

/* return true if a and b are the same size and have the same elements */
func equals(a, b []float64) bool {
	if len(a) == len(b) {
		for i, val := range a {
			if val != b[i] {
				return false
			}
		}
	}
	return true
}
