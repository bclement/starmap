package starmap

import(
    "fmt"
    "io"
    "geom"
    "os"
    "bytes"
    "bufio"
    "strconv"
    "io/ioutil"
    "strings"
    "path"
    "encoding/json"
)

const(
    wkt_start = iota
    wkt_outer = iota
    wkt_coord = iota
)

type Constellation struct {
    Name string
    Family string
    WktFiles []string
    Geoms []*geom.Polygon
    LabelPoints [][]float64
}

type Constellations []*Constellation

func LoadConstellations(constDir string) (Constellations, error) {
    infos, err := ioutil.ReadDir(constDir)
    if err != nil {
        return nil, err
    }
    suffix := ".json"
    rval := make(Constellations, 0, 88)
    for _, info := range(infos) {
        name := info.Name()
        if strings.HasSuffix(name, suffix) {
            fullPath := path.Join(constDir, name)
            constel, err := readJsonFile(fullPath)
            if err != nil {
                return nil, fmt.Errorf("Unable to parse %v: %v", fullPath, err)
            }
            constel.Geoms = make([]*geom.Polygon, 0, len(constel.WktFiles))
            for _, wktFile := range(constel.WktFiles) {
                fullWktPath := path.Join(constDir, wktFile)
                poly, err := readWktFile(fullWktPath)
                if err != nil {
                    return nil, fmt.Errorf("Unable to parse %v: %v",
                        fullWktPath, err)
                }
                constel.Geoms = append(constel.Geoms, poly)
            }
            rval = append(rval, constel)
        }
    }
    return rval, nil
}

func readJsonFile(path string) (*Constellation, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    dec := json.NewDecoder(f)
    var rval Constellation
    err = dec.Decode(&rval)
    return &rval, err
}

func readWktFile(path string) (*geom.Polygon, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
    state := wkt_start
    var b bytes.Buffer
    coordStrs := make([]string, 0, 16)
    reader := bufio.NewReader(f)
    for {
        r, _, err := reader.ReadRune()
        if err != nil {
            if err == io.EOF {
                break
            } else {
                return nil, err
            }
        }
        switch state {
        case wkt_start:
            if r == '(' {
                state = wkt_outer
            }
        case wkt_outer:
            if r == '(' {
                state = wkt_coord
            } else if r == ')' {
                /* inner holes not supported */
                break
            }
        case wkt_coord:
            if r == ',' {
                coordStrs = append(coordStrs, b.String())
                b.Reset()
            } else if r == ')' {
                state = wkt_outer
            } else {
                b.WriteRune(r)
            }
        }
    }
    coords := make([]float64, 0, len(coordStrs) * 2)
    prevDims := -1
    for _, coordStr := range(coordStrs) {
        floats := strings.Split(coordStr, " ")
        dims := 0
        for _, floatStr := range(floats) {
            if len(floatStr) < 1 {
                continue
            }
            c, err := strconv.ParseFloat(floatStr, 64)
            if err != nil {
                return nil, err
            }
            dims += 1
            coords = append(coords, c)
        }
        if prevDims < 0 {
            prevDims = dims
        } else if prevDims != dims {
            return nil, fmt.Errorf("mismatched dimensions in file: %v", path)
        }
    }
    return geom.NewPoly(prevDims, coords...)
}





