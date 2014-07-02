package geom

import (
	"bytes"
	"fmt"
	"math"
)

/*
base32Map := [32]rune{'0','1','2','3','4','5','6','7','8','9',
                        'b','c','d','e','f','g','h'.'j','k','m',
                        'n','p','q','r','s','t','u','v','w','x',
                        'y','z'}
*/
const BASE32 = "0123456789bcdefghjkmnpqrstuvwxyz"

/* generic coordinate storage */
type CoordinateSeq struct {
	Coords []float64
	Dims   int
}

/* returns number of coodinates in sequence */
func (cs *CoordinateSeq) Len() int {
	return len(cs.Coords) / cs.Dims
}

/* returns number of dimensions per coordinate */
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

/* interface for getting geohash strings */
type GeoHasher interface {
	/*
	   takes in x and y offsets from center of grid to extents
	   returns bounds as 40bit base32 geohash
	*/
	GeoHash(xcenter, ycenter, xoffset, yoffset float64) string
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

/* see GeoHasher interface */
func (p *Point) GeoHash(xcenter, ycenter, xoffset, yoffset float64) string {
	vals := make([]byte, 8)
	var i byte = 0
	for ; i < 20; i += 1 {
		xGlobalIndex := i * 2
		yGlobalIndex := xGlobalIndex + 1
		xoffset /= 2
		if p.c[0] >= xcenter {
			setBit(xGlobalIndex, vals)
			xcenter += xoffset
		} else {
			xcenter -= xoffset
		}
		yoffset /= 2
		if p.c[1] >= ycenter {
			setBit(yGlobalIndex, vals)
			ycenter += yoffset
		} else {
			ycenter -= yoffset
		}
	}
	var buffer bytes.Buffer
	for i := 0; i < len(vals); i += 1 {
		buffer.WriteByte(BASE32[vals[i]])
	}
	return buffer.String()
}

/*
treats vals as a contiguous bit array where
each byte in val holds 5 bits. Sets bit at
globalIndex to 1
*/
func setBit(globalIndex byte, vals []byte) {
	valIndex := globalIndex / 5
	localIndex := globalIndex % 5
	mask := byte(0x10 >> localIndex)
	vals[valIndex] |= mask
}

/* multi-point polygon */
type Polygon struct {
	/* storage for exterior ring */
	c CoordinateSeq
}

/*
takes in coordinate slice packed as [x1,y1,x2,y2...]
returns newly created 2D polygon
error if length of coords isn't divisible by 2
*/
func NewPoly2D(coords ...float64) (*Polygon, error) {
	if len(coords)%2 != 0 {
		return nil, fmt.Errorf("Invalid dimensions: %s", len(coords))
	}
	return &Polygon{CoordinateSeq{coords, 2}}, nil
}

/*
takes in coordinate slice packed as [x1,y1,z1,x2,y2,z2...]
returns newly created 3D polygon
error if length coords insn't divisible by 3
*/
func NewPoly3D(coords ...float64) (*Polygon, error) {
	if len(coords)%3 != 0 {
		return nil, fmt.Errorf("Invalid dimensions: %s", len(coords))
	}
	return &Polygon{CoordinateSeq{coords, 3}}, nil
}

/* see Geometry interface */
func (p *Polygon) Dims() int {
	return p.c.Dims
}

/* see Geometry interface */
func (p *Polygon) Coords() *CoordinateSeq {
	return &p.c
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

/*
takes in bounds, geometry and a comparison function
return true if bounds has enough dimensions to cover geometry and
all geometry coordinates return true for comp(coord, min) and
true for comp(max, coord) for bounds min and max
*/
func checkAll(bb *BoundingBox, g Geometry, comp func(float64, float64) bool) bool {
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
