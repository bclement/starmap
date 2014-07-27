package starmap

import(
    "geom"
    "testing"
)

func TestFilter(t *testing.T) {
    lower := geom.NewPoint2D(24, -90)
    upper := geom.NewPoint2D(22.5, -67.5)
    assertFilter(t, lower, upper, 10, "00", "08")
    lower = geom.NewPoint2D(13.5, 0)
    upper = geom.NewPoint2D(12, 22.5)
    assertFilter(t, lower, upper, 0, "", "")
}

func assertFilter(t *testing.T, lower, upper *geom.Point, num int,
        expLowHash, expUpHash string) {
    data, err := LoadData("../data/bright.tsv")
    if err != nil {
        t.Errorf("Can't load test data: %v", err)
    }
    width := 256
    height := 256
    trans := geom.CreateTransform(lower, upper, width, height, geom.STELLAR)
    lowerHash, upperHash := geom.BBoxHash(lower, upper, geom.STELLAR)
    if lowerHash != expLowHash || upperHash != expUpHash {
        t.Errorf("expected %v-%v, got %v-%v", expLowHash, expUpHash,
            lowerHash, upperHash)
    }
    stars := data.Range(lowerHash, upperHash)
    if len(stars) != num {
        t.Errorf("expected %v stars, got", num, len(stars))
    }
    for _, s := range(stars) {
        coord, err := geom.UnHash(s.GeoHash, geom.STELLAR)
        if err != nil {
            t.Errorf("problem unhashing: %v", err)
        }
        pix := trans.Transform(coord)
        if pix.X < 0 || pix.X >= width || pix.Y < 0 || pix.Y >= height {
            t.Errorf("problem transforming %v, %v", pix, coord)
        }
    }
}

