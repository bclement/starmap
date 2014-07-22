package starmap

import (
	"image/color"
	"image/png"
	"net/http"
    "geom"
	"render"
	"render/style"
    "math"
    "strconv"
    "strings"
)

const StarType string = "Star"

func init() {
	/* handler() defined below */
	http.HandleFunc("/", handler)
}

func raToGrid(ra float64) int {
    scaled := ra * (600.0 / 24.0)
    return int(math.Floor(600 - scaled))
}

func decToGrid(dec float64) int {
    scaled := (dec+90.0) * (300.0 / 180.0)
    return int(math.Floor(300 - scaled))
}

func doErr(w http.ResponseWriter, err error) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
}

/* parse integer url parameter
return defaultValue if parameter isn't present or is malformed */
func intParam(key string, defaultValue int, r *http.Request) int {
    value := r.FormValue(key)
    rval := defaultValue
    if value != "" {
        tmp, err := strconv.ParseInt(value, 10, 32)
        if err == nil {
            rval = int(tmp)
        }
    }
    return rval
}

/* parse float from string value
return defaultValue if value is malformed */
func parseFloat(value string, defaultValue float64) float64 {
    rval, err := strconv.ParseFloat(value, 64)
    if err != nil {
        rval = defaultValue
    }
    return rval
}

/* parse bounding box url parameter
return full bounds if parameter isn't present or is malformed */
func parseBbox(key string, r *http.Request) (*geom.Point, *geom.Point) {
    minx := 0.0
    miny := -90.0
    maxx := 24.0
    maxy := 90.0
    value := r.FormValue(key)
    if value != "" {
        parts := strings.Split(value, ",")
        if len(parts) == 4 {
            minx = parseFloat(parts[0], minx)
            miny = parseFloat(parts[1], miny)
            maxx = parseFloat(parts[2], maxx)
            maxy = parseFloat(parts[3], maxy)
        }
    }
    return geom.NewPoint2D(minx, miny), geom.NewPoint2D(maxx, maxy)
}

func handler(w http.ResponseWriter, r *http.Request) {
    data, err := LoadData("data/bright.tsv")
    if err != nil {
        doErr(w, err)
    }
	smlCircle := style.NewPointStyle(1, color.White, style.CIRCLE)
	midCircle := style.NewPointStyle(2, color.White, style.CIRCLE)
	lrgCircle1 := style.NewPointStyle(3, color.White, style.CIRCLE)
    width := intParam("width", 1024, r)
    height := intParam("height", 512, r)
    min, max := parseBbox("bbox", r)
    lowerHash, upperHash := geom.BBoxHash(min, max, geom.STELLAR)
    trans := geom.CreateTransform(min, max, width, height, geom.STELLAR)
	img := render.Create(width, height, color.Black)
    stars := data.Range(lowerHash, upperHash)
    for _, s := range(stars) {
        coord, err := geom.UnHash(s.GeoHash, geom.STELLAR)
        if err != nil {
            doErr(w, err)
            return
        }
        pix := trans.Transform(coord)
        mag := s.Magnitude
        if mag < 2 {
	        render.Render(img, pix, lrgCircle1)
        } else if mag < 3 {
            render.Render(img, pix, midCircle)
        } else if mag < 6 {
            render.Render(img, pix, smlCircle)
        }
    }
    if err = png.Encode(w, img); err != nil {
        doErr(w, err)
	}
}
