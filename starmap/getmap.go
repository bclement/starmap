package starmap

import (
    "fmt"
    "appengine"
    "appengine/memcache"
    "bytes"
	"geom"
	"image/color"
	"image/png"
	"net/http"
	"render"
	"render/style"
)


func createKey(width, height int, lower, upper *geom.Point) string {
    return fmt.Sprintf("%v-%v-%v-%v", width, height, lower, upper)
}

func getmap(w http.ResponseWriter, r *http.Request) {
	width := intParam("WIDTH", 1024, r)
	height := intParam("HEIGHT", 512, r)
	lower, upper := parseBbox("BBOX", r)
    ctx := appengine.NewContext(r)
    cacheKey := createKey(width, height, lower, upper)
    item, err := memcache.Get(ctx, cacheKey)
    if err == memcache.ErrCacheMiss {
        tile, err := createTile(w, width, height, lower, upper)
        if err != nil {
            doErr(w, err)
            return
        }
        item = &memcache.Item{Key:cacheKey, Value:tile}
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

func createTile(w http.ResponseWriter, width, height int,
        lower, upper *geom.Point) ([]byte, error) {
	if dataErr != nil {
        return nil, dataErr
	}
	smlCircle := style.NewPointStyle(0.5, color.White, style.CIRCLE)
	midCircle := style.NewPointStyle(1, color.White, style.CIRCLE)
	lrgCircle := style.NewPointStyle(2, color.White, style.CIRCLE)
	superCircle := style.NewPointStyle(3, color.White, style.CIRCLE)
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
		}else if mag < 0 {
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