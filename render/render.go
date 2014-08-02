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

func round(x float64) float64 {
    return float64(int(x +0.5))
}

func fpart(x float64) float64 {
    return x - float64(int(x))
}

func rfpart(x float64) float64 {
    return 1 - fpart(x)
}

func plot(img draw.Image, x, y int, c float64) {
    gray := uint8(255 * c)
    color := color.RGBA{gray, gray, gray, 255}
    img.Set(x, y, color)
}

func drawEndpoint(img draw.Image, x, y, gradient float64,
        first, steep bool) float64 {
    xend := round(x)
    yend := y + gradient * (xend - x)
    var xgap float64
    if first {
        xgap = rfpart(x + 0.5)
    } else {
        xgap = fpart(x + 0.5)
    }
    xpix := int(xend)
    ypix := int(yend)
    if steep {
        plot(img, ypix, xpix, rfpart(yend) * xgap)
        plot(img, ypix+1, xpix, fpart(yend) * xgap)
    } else {
        plot(img, xpix, ypix, rfpart(yend) * xgap)
        plot(img, xpix, ypix+1, fpart(yend) * xgap)
    }
    return yend + gradient
}

func RenderLine(img draw.Image, p0, p1 *image.Point, s *style.PolygonStyle) {
    steep := math.Abs(float64(p1.Y-p0.Y)) > math.Abs(float64(p1.X-p0.X))
    x0, y0 := float64(p0.X), float64(p0.Y)
    x1, y1 := float64(p1.X), float64(p1.Y)
    if steep {
        x0, y0 = y0, x0
        x1, y1 = y1, x1
    }
    if x0 > x1 {
        x0, x1 = x1, x0
        y0, y1 = y1, y0
    }
    dx := x1 - x0
    dy := y1 - y0
    gradient := dy / dx
    xpix1 := int(round(x0))
    xpix2 := int(round(x1))
    intery := drawEndpoint(img, x0, y0, gradient, true, steep)
    drawEndpoint(img, x1, y1, gradient, false, steep)

    for x := xpix1 + 1; x < xpix2; x += 1 {
        if steep {
            plot(img, int(intery), x, rfpart(intery))
            plot(img, int(intery)+1, x, fpart(intery))
        } else {
            plot(img, x, int(intery), rfpart(intery))
            plot(img, x, int(intery)+1, fpart(intery))
        }
        intery += gradient
    }
}

