package starmap

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"geom"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

/* wkt parsing states */
const (
	wkt_start = iota
	wkt_outer = iota
	wkt_coord = iota
	wkt_end   = iota
)

/* nested json polygon config struct */
type PolyInfo struct {
	WktFile    string
	LabelPoint []float64
	MaxScale   float64
	Geom       *geom.Polygon
}

type StringInfo struct {
	Name    string
	WktFile string
	Lines   []*geom.CoordinateSeq
}

/* top level json constellation config */
type Constellation struct {
	Name        string
	Family      string
	PolyInfos   []*PolyInfo
	StringInfos []*StringInfo
}

type Constellations []*Constellation

/* load contellation objects from static data directory */
func LoadConstellations(constDir string) (Constellations, error) {
	infos, err := ioutil.ReadDir(constDir)
	if err != nil {
		return nil, err
	}
	suffix := ".json"
	rval := make(Constellations, 0, 88)
	for _, info := range infos {
		name := info.Name()
		if strings.HasSuffix(name, suffix) {
			fullPath := path.Join(constDir, name)
			constel, err := readJsonFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("Unable to parse %v: %v", fullPath, err)
			}
			for i := range constel.PolyInfos {
				wktFile := constel.PolyInfos[i].WktFile
				fullWktPath := path.Join(constDir, wktFile)
				poly, err := readWktFile(fullWktPath)
				if err != nil {
					return nil, fmt.Errorf("Unable to parse %v: %v",
						fullWktPath, err)
				}
				constel.PolyInfos[i].Geom = poly
			}
			for i := range constel.StringInfos {
				wktFile := constel.StringInfos[i].WktFile
				fullWktPath := path.Join(constDir, wktFile)
				lines, err := readStringsWktFile(fullWktPath)
				if err != nil {
					return nil, fmt.Errorf("Unable to parse %v: %v",
						fullWktPath, err)
				}
				constel.StringInfos[i].Lines = lines
			}
			rval = append(rval, constel)
		}
	}
	return rval, nil
}

/* parse constellaton JSON config file */
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

func readStringsWktFile(path string) ([]*geom.CoordinateSeq, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	state := wkt_start
	reader := bufio.NewReader(f)
	rval := make([]*geom.CoordinateSeq, 0, 8)
	for state != wkt_end {
		coords, dims, newstate, err := readString(reader, state)
		state = newstate
		if err != nil {
			return nil, fmt.Errorf("Problem reading %v: %v", path, err)
		}
		seq := &geom.CoordinateSeq{coords, dims}
		rval = append(rval, seq)
	}
	return rval, nil
}

func readString(reader *bufio.Reader, state int) ([]float64, int, int, error) {
	var b bytes.Buffer
	coordStrs := make([]string, 0, 16)
	done := false
	for !done {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				state = wkt_end
				break
			} else {
				return nil, 0, state, err
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
				state = wkt_end
				done = true
				break
			}
		case wkt_coord:
			if r == ',' {
				coordStrs = append(coordStrs, b.String())
				b.Reset()
			} else if r == ')' {
				coordStrs = append(coordStrs, b.String())
				b.Reset()
				state = wkt_outer
				done = true
				break
			} else {
				b.WriteRune(r)
			}
		}
	}
	coords := make([]float64, 0, len(coordStrs)*2)
	prevDims := -1
	for _, coordStr := range coordStrs {
		floats := strings.Fields(coordStr)
		dims := 0
		for _, floatStr := range floats {
			if len(floatStr) < 1 {
				continue
			}
			c, err := strconv.ParseFloat(floatStr, 64)
			if err != nil {
				return nil, 0, state, err
			}
			dims += 1
			coords = append(coords, c)
		}
		if prevDims < 0 {
			prevDims = dims
		} else if prevDims != dims {
			err := fmt.Errorf("mismatched dimensions in file")
			return nil, 0, state, err
		}
	}
	return coords, prevDims, state, nil
}

/* parse constellation polygon well known text file */
func readWktFile(path string) (*geom.Polygon, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	state := wkt_start
	reader := bufio.NewReader(f)
	coords, dims, _, err := readString(reader, state)
	if err != nil {
		return nil, fmt.Errorf("Error reading file %v: %v", path, err)
	}
	return geom.NewPoly(dims, coords...)
}
