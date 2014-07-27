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
    xIncreasesRight bool
    yIncreasesUp bool
}

var LONLAT *GridDef = &GridDef{0, 0, 180, 90, true, true}
var STELLAR *GridDef = &GridDef{12, 0, 12, 90, false, true}

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
    px := p.c[0]
    if !gd.xIncreasesRight {
        xcenter = -xcenter
        px = -px
    }
    py := p.c[1]
    if !gd.yIncreasesUp {
        ycenter = -ycenter
        py = -py
    }
	var i byte = 0
	for ; i < 20; i += 1 {
		xGlobalIndex := i * 2
		yGlobalIndex := xGlobalIndex + 1
		xoffset /= 2
		if px >= xcenter {
			setBit(xGlobalIndex, vals)
			xcenter += xoffset
		} else {
			xcenter -= xoffset
		}
		yoffset /= 2
		if py >= ycenter {
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
    xinc := isSet
    if !gd.xIncreasesRight {
        xinc = isUnset
    }
    yinc := isSet
    if !gd.yIncreasesUp {
        yinc = isUnset
    }
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
		if xinc(xGlobalIndex, vals) {
			xcenter += xoffset
		} else {
			xcenter -= xoffset
		}
		yoffset /= 2
		if yinc(yGlobalIndex, vals) {
			ycenter += yoffset
		} else {
			ycenter -= yoffset
		}
	}
	return NewPoint2D(xcenter, ycenter), nil
}

func isUnset(globalIndex byte, vals []byte) bool {
    return !isSet(globalIndex, vals)
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

/*
takes in bounding box and grid definition
returns lower and upper geohash query strings
NOTE: query strings may contain non-base32 characters
*/
func BBoxHash(lower, upper *Point, gd *GridDef) (string, string) {
	/* copy since we change them in the loop */
	xcenter := gd.xcenter
	ycenter := gd.ycenter
	xoffset := gd.xoffset
	yoffset := gd.yoffset
    lowerx := lower.X()
    lowery := lower.Y()
    upperx := upper.X()
    uppery := upper.Y()
    if !gd.xIncreasesRight {
        xcenter = -xcenter
        lowerx = -lowerx
        upperx = -upperx
    }
    if !gd.yIncreasesUp {
        ycenter = -ycenter
        lowery = -lowery
        uppery = -uppery
    }
	/* holds the min geohash prefix */
	var minbuff bytes.Buffer
	/* replaces last byte in lower to create max prefix */
	var maxbyte byte
	done := false
    var i byte
	/* maxes out at 8 character geohash */
	for i = 0; !done && i < 8; i += 1 {
		var curr byte = 0
		var j byte = 0
		/* walk through each 5bit hash character */
		for ; j < 5; j += 1 {
            globalIndex := (i * 5) + j
			if globalIndex % 2 == 0 {
				xoffset /= 2
				if lowerx <= xcenter && upperx <= xcenter {
					xcenter -= xoffset
				} else if lowerx >= xcenter && upperx >= xcenter {
					curr |= 0x10 >> j
					xcenter += xoffset
				} else {
					done = true
				}
			} else {
				yoffset /= 2
				if lowery <= ycenter && uppery <= ycenter {
					ycenter -= yoffset
				} else if lowery >= ycenter && uppery >= ycenter {
					curr |= 0x10 >> j
					ycenter += yoffset
				} else {
					done = true
				}
			}
			if done {
				/* this includes all other characters under this bit prefix */
				maxbyte = curr + (0x20 >> j)
				break
			}
		}
		minbuff.WriteByte(BASE32[curr])
	}
	minhash := minbuff.String()
	minbuff.Truncate(minbuff.Len() - 1)
	if maxbyte < 32 {
		minbuff.WriteByte(BASE32[maxbyte])
	} else {
		/* tilde used as max value when we walk off base32 */
		minbuff.WriteByte(0x7e)
	}
	return minhash, minbuff.String()
}
