package starmap

import (
	"geom"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

var data Stardata
var dataErr error

var constelData Constellations
var constelErr error

var featureTemplate *template.Template
var templateErr error

func init() {
	/* handler() defined below */
	http.HandleFunc("/", handler)
	data, dataErr = LoadData("data/bright.tsv")
    constelData, constelErr = LoadConstellations("data/consts")
	featureTemplate, templateErr =
		template.ParseFiles("templates/getfeatureinfo.template")
}

func doErr(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func strParam(key, defaultValue string, r *http.Request) string {
    rval := r.FormValue(key)
    if rval == "" {
        rval = defaultValue
    }
    return rval
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
	rightx := 0.0
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
	request := r.FormValue("REQUEST")
	if strings.EqualFold(request, "GETFEATUREINFO") {
		getfeatureinfo(w, r)
	} else {
		getmap(w, r)
	}
}
