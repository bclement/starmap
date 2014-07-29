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
    leftx := 24.0
    lowery := -90.0
    rightx:= 0.0
    uppery := 90.0
    value := r.FormValue(key)
    if value != "" {
        parts := strings.Split(value, ",")
        if len(parts) == 4 {
            x0 := parseFloat(parts[0], leftx)
            y0 := parseFloat(parts[1], lowery)
            x1 := parseFloat(parts[2], rightx)
            y1 := parseFloat(parts[3], uppery)
            /* in stellar coordinates, 24 is left of 0 */
            if x0 > x1 {
                leftx = x0
                rightx = x1
            } else {
                leftx = x1
                rightx = x0
            }
            if y0 > y1 {
                lowery = y1
                uppery = y0
            } else {
                lowery = y0
                uppery = y1
            }
        }
    }
    return geom.NewPoint2D(leftx, lowery), geom.NewPoint2D(rightx, uppery)
}

func handler(w http.ResponseWriter, r *http.Request) {
    data, err := LoadData("data/bright.tsv")
    if err != nil {
        doErr(w, err)
    }
	smlCircle := style.NewPointStyle(0.5, color.White, style.CIRCLE)
	midCircle := style.NewPointStyle(1, color.White, style.CIRCLE)
	lrgCircle := style.NewPointStyle(2, color.White, style.CIRCLE)
	superCircle := style.NewPointStyle(3, color.White, style.CIRCLE)
    width := intParam("WIDTH", 1024, r)
    height := intParam("HEIGHT", 512, r)
    lower, upper := parseBbox("BBOX", r)
    lowerHash, upperHash := geom.BBoxHash(lower, upper, geom.STELLAR)
    trans := geom.CreateTransform(lower, upper, width, height, geom.STELLAR)
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
            gray = uint8((2.0 - mag) * 64.0) + 128
        } else if mag < 4 {
            style = midCircle
            gray = uint8((4.0 - mag) * 64.0) + 128
        } else {
            gray = uint8((6.0 - mag) * 64.0) + 128
        }
        color := color.RGBA{ gray, gray, gray, 255}
        style.Style.Color = color
        render.Render(img, pix, style)
    }
    if err = png.Encode(w, img); err != nil {
        doErr(w, err)
	}
}
