package style

import "image/color"

type Shape int

/* values for Shape type */
const (
    CIRCLE = iota
    SQUARE = iota
)

/* generic style data */
type Style struct {
    /* see specific style types for use */
    Size float64
    Color color.Color
}

/* style data for individual points */
type PointStyle struct {
    Style
    /* render shape of point */
    Shape Shape
}

/* 
takes in the size of the radius of the point, color and shape
returns a pointer to a newly created point style 
*/
func NewPointStyle(size float64, color color.Color, shape Shape) *PointStyle {
    s := Style{size, color}
    return &PointStyle{s, shape}
}

/* style data for polygons */
type PolygonStyle struct {
    Style
}

func NewPolyStyle(size float64, color color.Color) *PolygonStyle {
    return &PolygonStyle{Style{size, color}}
}
