package geom

import (
	"image"
	"math"
)

/* container to hold invariants used in bulk transform operations */
type PointTransform struct {
	Dx     float64
	Dy     float64
	Max    *Point
	Width  int
	Height int
	gd     *GridDef
}

/*
create transform object from lower/upper bounds, image dimensions
and grid definition
*/
func CreateTransform(lowerLeft, upperRight *Point, width, height int,
	gd *GridDef) *PointTransform {
	worldNx := math.Abs(lowerLeft.X() - upperRight.X())
	worldNy := math.Abs(lowerLeft.Y() - upperRight.Y())
	dx := worldNx / float64(width)
	dy := worldNy / float64(height)
	maxx := math.Max(lowerLeft.X(), upperRight.X())
	maxy := math.Max(lowerLeft.Y(), upperRight.Y())
	max := NewPoint2D(maxx, maxy)
	return &PointTransform{dx, dy, max, width, height, gd}
}

/* take in a point in spatial dimensions, return image pixel location */
func (pt *PointTransform) Transform(p *Point) *image.Point {
	return pt.TransformXY(p.X(), p.Y())
}

/* take in a 2D point in spatial dimensions, return image pixel location */
func (pt *PointTransform) TransformCoord(c []float64) *image.Point {
	return pt.TransformXY(c[0], c[1])
}

/* take in a 2D point in spatial dimensions, return image pixel location */
func (pt *PointTransform) TransformXY(x, y float64) *image.Point {
	rawX := (pt.Max.X() - x) / pt.Dx
	if pt.gd.xIncreasesRight {
		rawX = float64(pt.Width) - rawX
	}
	xpix := int(math.Floor(rawX))
	rawY := (pt.Max.Y() - y) / pt.Dy
	if !pt.gd.yIncreasesUp {
		rawY = float64(pt.Height) - rawY
	}
	ypix := int(math.Floor(rawY))
	return &image.Point{xpix, ypix}
}

/* take in an image pixel and return a point in spatial dimensions */
func (pt *PointTransform) Reverse(p *image.Point) *Point {
	tmpX := float64(p.X)
	if pt.gd.xIncreasesRight {
		tmpX += float64(pt.Width)
	}
	tmpX *= pt.Dx
	tmpX -= pt.Max.X()
	tmpX = -tmpX

	tmpY := float64(p.Y)
	if !pt.gd.yIncreasesUp {
		tmpY += float64(pt.Height)
	}
	tmpY *= pt.Dy
	tmpY -= pt.Max.Y()
	tmpY = -tmpY
	return NewPoint2D(tmpX, tmpY)
}
