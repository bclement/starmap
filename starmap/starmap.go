package starmap

import (
	"image/color"
	"image/png"
	"net/http"
    "geom"
	"render"
	"render/style"
    "math"
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

func handler(w http.ResponseWriter, r *http.Request) {
    data, err := LoadData("data/bright.tsv")
    if err != nil {
        doErr(w, err)
    }
	smlCircle := style.NewPointStyle(1, color.White, style.CIRCLE)
	midCircle := style.NewPointStyle(2, color.White, style.CIRCLE)
	lrgCircle1 := style.NewPointStyle(3, color.White, style.CIRCLE)
    width := 1024
    height := 512
    min := geom.NewPoint2D(0, -90)
    max := geom.NewPoint2D(24, 90)
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
