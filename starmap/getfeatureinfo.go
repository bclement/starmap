package starmap

import (
	"fmt"
	"geom"
	"image"
	"net/http"
)

const noDataFeatureInfo = "<html><body>no data</body></html>"

type Param struct {
	Key string
	Val string
}

/* struct used by template to generate output */
type Feature struct {
	Name   string
	Params []Param
}

/* takes in a star and converts it to a parameter slice */
func asParams(star *Star) []Param {
	rval := make([]Param, 0, 6)
	rval = addParam(rval, "hipparchose #", fmt.Sprintf("%v", star.HipNum))
	rval = addParam(rval, "name", star.Name)
	rval = addParam(rval, "magnitude", fmt.Sprintf("%v", star.Magnitude))
	coord, err := geom.UnHash(star.GeoHash, geom.STELLAR)
	if err == nil {
		rval = addParam(rval, "right ascension",
			fmt.Sprintf("%0.5f", coord.X()))
		rval = addParam(rval, "declination", fmt.Sprintf("%0.5f", coord.Y()))
	}
	return rval
}

func addParam(dest []Param, key, val string) []Param {
	return append(dest, Param{key, val})
}

func getfeatureinfo(w http.ResponseWriter, r *http.Request) {
	if templateErr != nil {
		doErr(w, templateErr)
		return
	}
	width := intParam("WIDTH", 1024, r)
	height := intParam("HEIGHT", 512, r)
	lower, upper := parseBbox("BBOX", r)
	i := intParam("X", 0, r)
	j := intParam("Y", 0, r)
	trans := geom.CreateTransform(lower, upper, width, height, geom.STELLAR)
	coord := trans.Reverse(&image.Point{i, j})
	star := data.FindClosest(coord)
	if star != nil {
		features := []Feature{Feature{"star", asParams(star)}}
		err := featureTemplate.Execute(w, features)
		if err != nil {
			doErr(w, err)
		}
	} else {
		w.Write([]byte(noDataFeatureInfo))
	}
}
