package starmap

import (
	"image/color"
	"image/png"
	"net/http"
    "geom"
	"render"
	"render/style"
    "strconv"
    "os"
    "strings"
    "bufio"
    "appengine"
    "appengine/datastore"
    "math"
)

const StarType string = "Star"

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

func loadData(c appengine.Context) error {
    f, err := os.Open("data/bright.tsv")
    if err != nil {
        return err
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    starBuff := make([]*Star, 2048)
    keyBuff := make([]*datastore.Key, 2048)
    index := 0
    for scanner.Scan() {
        star := new(Star)
        line := scanner.Text()
        parts := strings.Split(line, "\t")
        if len(parts) != 6 {
            continue
        }
        id, idErr := strconv.ParseInt(parts[0], 10, 32)
        hip, hipErr := strconv.ParseInt(parts[1], 10, 32)
        if hipErr != nil {
            hip = 0
        }
        star.HipNum = int32(hip)
        star.Name = parts[2]
        ra, raErr := strconv.ParseFloat(parts[3], 64)
        dec, decErr := strconv.ParseFloat(parts[4],64)
        mag, magErr := strconv.ParseFloat(parts[5], 64)
        if idErr != nil || raErr != nil || decErr != nil || magErr != nil {
            continue
        }
        star.Magnitude = mag
        coord := geom.NewPoint2D(ra, dec)
        star.GeoHash = coord.GeoHash(geom.STELLAR)
        key := datastore.NewKey(c, StarType, "", id, nil)
        starBuff[index] = star
        keyBuff[index] = key
        index += 1
        if index >= len(starBuff) {
            _, err = datastore.PutMulti(c, keyBuff, starBuff)
            if err != nil {
                return err
            }
            index = 0
        }
    }
    if index > 0 {
        starBuff = starBuff[:index+1]
        keyBuff = keyBuff[:index+1]
        _, err = datastore.PutMulti(c, keyBuff, starBuff)
        if err != nil {
            return err
        }
    }
    return nil
}

func doErr(w http.ResponseWriter, err error) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
}

func handler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery(StarType)
    count, err := q.Count(c)
    if count < 1000 {
        err = loadData(c)
        if err != nil {
            doErr(w, err)
            return
        }
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
    q = datastore.NewQuery(StarType).
        Filter("GeoHash >=", lowerHash).Filter("GeoHash <", upperHash)
    for t:= q.Run(c); ; {
        var s Star
        _, err = t.Next(&s)
        if err == datastore.Done {
            break
        }
        if err != nil {
            doErr(w, err)
            return
        }
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
