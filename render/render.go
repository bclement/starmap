package render

import (
	"github.com/bclement/starmap/render/style"
	"image"
	"image/color"
	"image/draw"
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
	r int
}

/* see image.Image interface */
func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

/* see image.Image interface */
func (c *circle) Bounds() image.Rectangle {
	return centeredRect(&c.p, c.r)
}

/* see image.Image interface */
func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5,
		float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
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
		mask = &square{centeredRect(p, pstyle.Style.Size)}
	}
	draw.DrawMask(img, img.Bounds(), &image.Uniform{pstyle.Color},
		image.ZP, mask, image.ZP, draw.Over)
}
