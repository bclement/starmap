package starmap

import (
	"appengine"
	"appengine/memcache"
	"bytes"
	"fmt"
	"geom"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"render"
	"render/style"
	"strings"
)

var smlCircle = style.NewPointStyle(0.5, color.White, style.CIRCLE)
var midCircle = style.NewPointStyle(1, color.White, style.CIRCLE)
var lrgCircle = style.NewPointStyle(2, color.White, style.CIRCLE)
var superCircle = style.NewPointStyle(3, color.White, style.CIRCLE)

var labelColors = map[string]color.Color{
	"Heavenly Waters": color.RGBA{0, 154, 205, 255},
	"Hercules":        color.RGBA{34, 139, 34, 255},
	"Ursa Major":      color.RGBA{100, 149, 237, 255},
	"Perseus":         color.RGBA{225, 58, 58, 255},
	"Orion":           color.RGBA{205, 102, 0, 255},
	"Bayer":           color.RGBA{205, 149, 12, 255},
	"La Caille":       color.RGBA{137, 104, 205, 255},
}

/* create the cache key for a WMS tile */
func createKey(layer string, width, height int, lower,
	upper *geom.Point) string {
	return fmt.Sprintf("%v-%v-%v-%v-%v", layer, width, height, lower, upper)
}

/* WMS getmap handler function */
func getmap(w http.ResponseWriter, r *http.Request) {
	width := intParam("WIDTH", 1024, r)
	height := intParam("HEIGHT", 512, r)
	lower, upper := parseBbox("BBOX", r)
	layer := strParam("LAYERS", "stars", r)
	ctx := appengine.NewContext(r)
	cacheKey := createKey(layer, width, height, lower, upper)
	item, err := memcache.Get(ctx, cacheKey)
	if err == memcache.ErrCacheMiss {
		tile, err := createTile(w, layer, width, height, lower, upper)
		if err != nil {
			doErr(w, err)
			return
		}
		item = &memcache.Item{Key: cacheKey, Value: tile}
		err = memcache.Add(ctx, item)
		if err != nil {
			doErr(w, err)
			return
		}
	} else if err != nil {
		doErr(w, err)
		return
	}

	w.Write(item.Value)
}

/* create a new tile image for request */
func createTile(w http.ResponseWriter, layer string, width, height int,
	lower, upper *geom.Point) ([]byte, error) {
	layer = strings.ToLower(layer)
	if layer == "constellations" {
		return createConstTile(w, width, height, lower, upper)
	} else {
		return createStarTile(w, width, height, lower, upper)
	}
}

/* create a constellation layer tile */
func createConstTile(w http.ResponseWriter, width, height int,
	lower, upper *geom.Point) ([]byte, error) {
	if constelErr != nil {
		return nil, constelErr
	}
	scale := math.Abs(upper.X()-lower.X()) / float64(width)
	s := style.NewPolyStyle(1, color.White)
	trans := geom.CreateTransform(lower, upper, width, height, geom.STELLAR)
	img := render.CreateTransparent(width, height)
	bbox := geom.NewBBox2D(lower.X(), lower.Y(), upper.X(), upper.Y())
	for _, c := range constelData {
		txtColor := labelColors[c.Family]
		if txtColor == nil {
			txtColor = color.White
		}
		for _, pi := range c.PolyInfos {
			if bbox.Touches(pi.Geom) {
				render.RenderPoly(img, pi.Geom, trans, s)
				if charsErr == nil && pi.LabelPoint != nil &&
					pi.MaxScale > scale {
					labelPoint := pi.LabelPoint
					pix := trans.TransformXY(labelPoint[0], labelPoint[1])
					render.RenderString(img, chars, 10, pix, c.Name, txtColor)
				}
			}
		}
	}
	var rval bytes.Buffer
	if err := png.Encode(&rval, img); err != nil {
		return nil, err
	}
	return rval.Bytes(), nil
}

/* create a star layer tile */
func createStarTile(w http.ResponseWriter, width, height int,
	lower, upper *geom.Point) ([]byte, error) {
	if dataErr != nil {
		return nil, dataErr
	}
	lowerHash, upperHash := geom.BBoxHash(lower, upper, geom.STELLAR)
	trans := geom.CreateTransform(lower, upper, width, height, geom.STELLAR)
	img := render.Create(width, height, color.Black)
	stars := data.Range(lowerHash, upperHash)
	for _, s := range stars {
		coord, err := geom.UnHash(s.GeoHash, geom.STELLAR)
		if err != nil {
			return nil, err
		}
		pix := trans.Transform(coord)
		mag := s.Magnitude
		style := smlCircle
		var gray uint8
		if mag < -1 {
			style = superCircle
			gray = 255
		} else if mag < 0 {
			style = superCircle
			gray = 200
		} else if mag < 2 {
			style = lrgCircle
			gray = uint8((2.0-mag)*64.0) + 128
		} else if mag < 4 {
			style = midCircle
			gray = uint8((4.0-mag)*64.0) + 128
		} else {
			gray = uint8((6.0-mag)*64.0) + 128
		}
		color := color.RGBA{gray, gray, gray, 255}
		style.Style.Color = color
		render.Render(img, pix, style)
	}
	var rval bytes.Buffer
	if err := png.Encode(&rval, img); err != nil {
		return nil, err
	}
	return rval.Bytes(), nil
}
