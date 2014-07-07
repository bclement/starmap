package geom

import(
    "image"
    "math"
)

/* container to hold invariants used in bulk transform operations */
type PointTransform struct {
    Dx float64
    Dy float64
    Max *Point
    Width int
    Height int
    gd *GridDef
}

/* 
create transform object from min/max bounds, image dimensions
and grid definition
*/
func CreateTransform(min, max *Point, width, height int,
        gd *GridDef) *PointTransform {
    worldNx := max.X() - min.X()
    worldNy := max.Y() - min.Y()
    dx := worldNx / float64(width)
    dy := worldNy / float64(height)
    return &PointTransform{dx, dy, max, width, height, gd}
}

/* take in a point in spatial dimensions, return image pixel location */
func (pt *PointTransform) Transform(p *Point) *image.Point {
    rawX := (pt.Max.X() - p.X()) / pt.Dx
    if pt.gd.xIncreasesRight {
        rawX = float64(pt.Width) - rawX
    }
    x := int(math.Floor(rawX))
    rawY := (pt.Max.Y() - p.Y()) / pt.Dy
    if !pt.gd.yIncreasesUp {
        rawY = float64(pt.Height) - rawY
    }
    y := int(math.Floor(rawY))
    return &image.Point{x, y}
}

