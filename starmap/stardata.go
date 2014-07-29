package starmap

import (
	"bufio"
	"geom"
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
