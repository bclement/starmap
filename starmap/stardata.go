package starmap

import (
	"bufio"
	"geom"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Star struct {
	/* may be blank */
	Name string
	/* Hipparchos start catalog number */
	HipNum int32
	/* more negative is brighter */
	Magnitude float64
	/* geohash of right declination and ascension */
	GeoHash string
}

type Stardata []*Star

/* sort interface */
func (sd Stardata) Len() int {
	return len(sd)
}

/* sort interface */
func (sd Stardata) Swap(i, j int) {
	sd[i], sd[j] = sd[j], sd[i]
}

/* sort interface */
func (sd Stardata) Less(i, j int) bool {
	return sd[i].GeoHash < sd[j].GeoHash
}

/* takes in a geohash key, returns index into data slice
see sort.Search() for more details */
func (sd Stardata) FindIndex(geohash string) int {
	return sort.Search(len(sd), func(i int) bool {
		return sd[i].GeoHash >= geohash
	})
}

/* takes in a point, returns closest star or nil if not found */
func (sd Stardata) FindClosest(p *geom.Point) *Star {
	lower := geom.NewPoint2D(p.X()+0.5, p.Y()-1)
	upper := geom.NewPoint2D(p.X()-0.5, p.Y()+1)
	lowerHash, upperHash := geom.BBoxHash(lower, upper, geom.STELLAR)

	var rval *Star = nil
	var minDist float64 = math.MaxFloat64

	stars := sd.Range(lowerHash, upperHash)
	for _, s := range stars {
		coord, err := geom.UnHash(s.GeoHash, geom.STELLAR)
		if err != nil {
			continue
		}
		x := math.Abs(p.X() - coord.X())
		xx := x * x
		y := math.Abs(p.Y() - coord.Y())
		yy := y * y
		zz := xx + yy
		if minDist > zz {
			minDist = zz
			rval = s
		}
	}
	return rval
}

/* takes in two strings returns the length of the common prefix */
func sharedLen(a, b string) int {
	length := int(math.Min(float64(len(a)), float64(len(b))))
	rval := 0
	for ; rval < length; rval += 1 {
		if a[rval] != b[rval] {
			break
		}
	}
	return rval
}

/* returns the star that has a matching geohash or nil if not found */
func (sd Stardata) Find(geohash string) *Star {
	i := sd.FindIndex(geohash)
	if i < len(sd) && sd[i].GeoHash == geohash {
		return sd[i]
	} else {
		return nil
	}
}

/* takes in two geohash search keys
returns slice of data that is between those two keys */
func (sd Stardata) Range(start, end string) Stardata {
	startIndex := sd.FindIndex(start)
	endIndex := sd.FindIndex(end)
	return sd[startIndex:endIndex]
}

/* load static star data from tsv file */
func LoadData(datafile string) (Stardata, error) {
	f, err := os.Open(datafile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	rval := make(Stardata, 0, 32)
	for scanner.Scan() {
		star := new(Star)
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) != 6 {
			continue
		}
		//id, idErr := strconv.ParseInt(parts[0], 10, 32)
		hip, hipErr := strconv.ParseInt(parts[1], 10, 32)
		if hipErr != nil {
			hip = 0
		}
		star.HipNum = int32(hip)
		star.Name = parts[2]
		ra, raErr := strconv.ParseFloat(parts[3], 64)
		dec, decErr := strconv.ParseFloat(parts[4], 64)
		mag, magErr := strconv.ParseFloat(parts[5], 64)
		if raErr != nil || decErr != nil || magErr != nil {
			continue
		}
		star.Magnitude = mag
		coord := geom.NewPoint2D(ra, dec)
		star.GeoHash = coord.GeoHash(geom.STELLAR)
		rval = append(rval, star)
	}
	sort.Sort(rval)
	return rval, nil
}
