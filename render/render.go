package render

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"render/style"
)

/*
take in center point and offset
return rectangle with lower bounds center minus offset
    and upper bounds center plus offset
*/
func centeredRect(p *image.Point, offset int) image.Rectangle {
	return image.Rect(p.X-offset, p.Y-offset, p.X+offset,
		p.Y+offset)
}

/* circular mask for point rendering */
type circle struct {
	p image.Point
	r float64
}

/* see image.Image interface */
func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

/* see image.Image interface */
func (c *circle) Bounds() image.Rectangle {
	return centeredRect(&c.p, int(math.Ceil(c.r)))
}

/* see image.Image interface */
func (c *circle) At(x, y int) color.Color {
	rval := color.Alpha{0}
	if c.r <= 0.5 && x == c.p.X && y == c.p.Y {
		rval = color.Alpha{255}
	} else {
		xx, yy, rr := float64(x-c.p.X)+0.5,
			float64(y-c.p.Y)+0.5, c.r
		if xx*xx+yy*yy < rr*rr {
			rval = color.Alpha{255}
		}
	}
	return rval
}

/* square mask for point rendering */
type square struct {
	image.Rectangle
}

/* see image.Image interface */
func (s *square) ColorModel() color.Model {
	return color.AlphaModel
}

/* see image.Image interface */
func (s *square) Bounds() image.Rectangle {
	return s.Rectangle
}

/* see image.Image interface */
func (s *square) At(x, y int) color.Color {
	min := s.Rectangle.Min
	max := s.Rectangle.Max
	if x >= min.X && x <= max.X && y >= min.Y && y <= max.Y {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

/*
takes in dimensions and background color
returns newly created image
*/
func Create(width, height int, bgcolor color.Color) draw.Image {
	/* create image with background color */
	rval := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rval, rval.Bounds(), &image.Uniform{bgcolor}, image.ZP, draw.Src)
	return rval
}

/*
takes in point to render onto image using style
*/
func Render(img draw.Image, p *image.Point, pstyle *style.PointStyle) {
	var mask image.Image
	/* TODO don't create new mask every time */
	if pstyle.Shape == style.CIRCLE {
		mask = &circle{*p, pstyle.Style.Size}
	} else {
		size := int(math.Ceil(pstyle.Style.Size))
		mask = &square{centeredRect(p, size)}
	}
	draw.DrawMask(img, img.Bounds(), &image.Uniform{pstyle.Color},
		image.ZP, mask, image.ZP, draw.Over)
}

func RenderLine(img draw.Image, p0, p1 *image.Point, s *style.PolygonStyle) {
    run := float64(p1.X - p0.X)
    rise := float64(p1.Y - p0.Y)
    length := math.Hypot(rise, run)
    dx := run / length
    dy := rise / length
    total := int(math.Ceil(length))
    for i := 0; i < total; i += 1 {
        x := float64(p0.X) + float64(i) * dx
        y := float64(p0.Y) + float64(i) * dy
        img.Set(int(x), int(y), s.Color)
    }
}
