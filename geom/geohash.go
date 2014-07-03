package geom

import (
	"bytes"
	"fmt"
	"strings"
)

const BASE32 = "0123456789bcdefghjkmnpqrstuvwxyz"

/* defines a grid in terms of center point and offsets */
type GridDef struct {
	xcenter float64
	ycenter float64
	/* offset is distance from center to grid boundary */
	xoffset float64
	yoffset float64
}

var LONLAT *GridDef = &GridDef{0, 0, 180, 90}
var STELLAR *GridDef = &GridDef{12, 0, 12, 90}

/* interface for getting geohash strings */
type GeoHasher interface {
	/*
	   takes in grid definition
	   returns bounds as 40bit base32 geohash
	*/
	GeoHash(gd *GridDef) string
}

/* see GeoHasher interface */
func (p *Point) GeoHash(gd *GridDef) string {
	/* copy since we change them in the loop */
	xcenter := gd.xcenter
	ycenter := gd.ycenter
	xoffset := gd.xoffset
	yoffset := gd.yoffset
	vals := make([]byte, 8)
	var i byte = 0
	for ; i < 20; i += 1 {
		xGlobalIndex := i * 2
		yGlobalIndex := xGlobalIndex + 1
		xoffset /= 2
		if p.c[0] >= xcenter {
			setBit(xGlobalIndex, vals)
			xcenter += xoffset
		} else {
			xcenter -= xoffset
		}
		yoffset /= 2
		if p.c[1] >= ycenter {
			setBit(yGlobalIndex, vals)
			ycenter += yoffset
		} else {
			ycenter -= yoffset
		}
	}
	var buffer bytes.Buffer
	for i := 0; i < len(vals); i += 1 {
		buffer.WriteByte(BASE32[vals[i]])
	}
	return buffer.String()
}

/*
treats vals as a contiguous bit array where
each byte in val holds 5 bits. Sets bit at
globalIndex to 1
*/
func setBit(globalIndex byte, vals []byte) {
	valIndex := globalIndex / 5
	localIndex := globalIndex % 5
	mask := byte(0x10 >> localIndex)
	vals[valIndex] |= mask
}

/*
takes in geohash string and grid definition
returns 2D point or error if geohash is invalid
*/
func UnHash(hash string, gd *GridDef) (*Point, error) {
	/* copy since we change them in the loop */
	xcenter := gd.xcenter
	ycenter := gd.ycenter
	xoffset := gd.xoffset
	yoffset := gd.yoffset
	vals := make([]byte, 8)
	for i := 0; i < len(vals); i += 1 {
		index := strings.IndexByte(BASE32, hash[i])
		if index < 0 {
			return nil, fmt.Errorf("Invalid GeoHash character: %v", hash[i])
		}
		vals[i] = byte(index)
	}
	var i byte = 0
	for ; i < 20; i += 1 {
		xGlobalIndex := i * 2
		yGlobalIndex := xGlobalIndex + 1
		xoffset /= 2
		if isSet(xGlobalIndex, vals) {
			xcenter += xoffset
		} else {
			xcenter -= xoffset
		}
		yoffset /= 2
		if isSet(yGlobalIndex, vals) {
			ycenter += yoffset
		} else {
			ycenter -= yoffset
		}
	}
	return NewPoint2D(xcenter, ycenter), nil
}

/*
treats vals as a contiguous bit array where
each byte in val holds 5 bits.
returns true if bit at globalIndex is 1
*/
func isSet(globalIndex byte, vals []byte) bool {
	valIndex := globalIndex / 5
	localIndex := globalIndex % 5
	mask := byte(0x10 >> localIndex)
	return (vals[valIndex] & mask) != 0
}
