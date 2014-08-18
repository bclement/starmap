package starmap

import (
	"appengine"
	"appengine/memcache"
	"bytes"
	"fmt"
	"geom"
	"image/color"
	"image/png"
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
func createKey(r *Req) string {
	return fmt.Sprintf("%v-%v-%v-%v-%v", r.Layer, r.Width, r.Height,
		r.Lower, r.Upper)
}

/* WMS getmap handler function */
func getmap(w http.ResponseWriter, r *http.Request) {
	req := ParseReq(r)
	ctx := appengine.NewContext(r)
	cacheKey := createKey(req)
	item, err := memcache.Get(ctx, cacheKey)
	if err == memcache.ErrCacheMiss {
		tile, err := createTile(w, req)
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
func createTile(w http.ResponseWriter, req *Req) ([]byte, error) {
	layer := strings.ToLower(req.Layer)
	if layer == "constellations" {
		return createConstTile(w, req)
	} else if layer == "asterisms" {
		return createAsterTile(w, req)
	} else {
		return createStarTile(w, req)
	}
}

/* create a constellation layer tile */
func createConstTile(w http.ResponseWriter, req *Req) ([]byte, error) {
	if constelErr != nil {
		return nil, constelErr
	}
	scale := req.Scale()
	s := style.NewPolyStyle(1, color.White)
	trans := req.Trans(geom.STELLAR)
	img := render.CreateTransparent(req.Width, req.Height)
	bbox := req.BBox()
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

func createAsterTile(w http.ResponseWriter, req *Req) ([]byte, error) {
	if constelErr != nil {
		return nil, constelErr
	}
	s := style.NewPolyStyle(1, color.White)
	trans := req.Trans(geom.STELLAR)
	img := render.CreateTransparent(req.Width, req.Height)
	bbox := req.BBox()
	for _, c := range constelData {
		for _, si := range c.StringInfos {
			for _, cs := range si.Lines {
				if bbox.TouchesSeq(cs) {
					render.RenderSeq(img, cs, trans, s)
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
func createStarTile(w http.ResponseWriter, req *Req) ([]byte, error) {
	lowerHash, upperHash := geom.BBoxHash(req.Lower, req.Upper, geom.STELLAR)
	trans := req.Trans(geom.STELLAR)
	img := render.Create(req.Width, req.Height, color.Black)

	sr := &StarReq{req, make(chan Stardata)}
	starReqChan <- sr
	for data := range sr.out {
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
			} else if mag < 20 {
				gray = uint8((20.0-mag)*12.0) + 64
			} else {
				gray = 64
			}
			color := color.RGBA{gray, gray, gray, 255}
			style.Style.Color = color
			render.Render(img, pix, style)
		}
	}
	var rval bytes.Buffer
	if err := png.Encode(&rval, img); err != nil {
		return nil, err
	}
	return rval.Bytes(), nil
}
