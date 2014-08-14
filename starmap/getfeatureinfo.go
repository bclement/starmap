package starmap

import (
	"fmt"
	"geom"
	"image"
	"net/http"
	"strings"
)

const noDataFeatureInfo = "<html><body>no data</body></html>"

/* feature parameter key/val pair */
type Param struct {
	Key string
	Val string
}

/* struct used by template to generate output */
type Feature struct {
	Name   string
	Params []Param
}

/* convert a contellation object into feature parameters */
func constAsParams(constel *Constellation) []Param {
	rval := make([]Param, 0, 3)
	rval = addParam(rval, "name", constel.Name)
	rval = addParam(rval, "family", constel.Family)
	return rval
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

/* create parameter object and append to dest */
func addParam(dest []Param, key, val string) []Param {
	return append(dest, Param{key, val})
}

/* handler method for WMS get feature info requests */
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
	layersParam := strParam("LAYERS", "stars", r)
	trans := geom.CreateTransform(lower, upper, width, height, geom.STELLAR)
	coord := trans.Reverse(&image.Point{i, j})
	layers := strings.Split(layersParam, ",")
	features := make([]*Feature, 0, 3)
	for _, layer := range layers {
		if layer == "stars" {
			sf := starFeatures(coord)
			features = append(features, sf...)
		} else if layer == "constellations" {
			cf := constelFeatures(coord)
			features = append(features, cf...)
		}
	}
	if len(features) > 0 {
		err := featureTemplate.Execute(w, features)
		if err != nil {
			doErr(w, err)
		}
	} else {
		w.Write([]byte(noDataFeatureInfo))
	}
}

/* get star layer feature info for point */
func starFeatures(point *geom.Point) []*Feature {
	star := data.FindClosest(point)
	if star != nil {
		return []*Feature{&Feature{"star", asParams(star)}}
	} else {
		return []*Feature{}
	}
}

/* get contellation layer feature info for point */
func constelFeatures(point *geom.Point) []*Feature {
	rval := make([]*Feature, 0, 2)
	for _, c := range constelData {
		for _, pi := range c.PolyInfos {
			if pi.Geom.Contains(point) {
				rval = append(rval, &Feature{"constellation",
					constAsParams(c)})
			}
		}
	}
	return rval
}
