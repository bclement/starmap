#!/usr/bin/env python

import sys
import shapely
import shapely.wkt
import shapely.geometry

for arg in sys.argv[1:]:
    print(arg)
    with open(arg, 'r') as f:
        wkt = f.read()
        poly = shapely.wkt.loads(wkt)
        coords = poly.exterior.coords
        left = []
        right = []
        if coords[0][0] > 20:
            left.append(coords[0])
            lastLeft = True
        else:
            right.append(coords[0])
            lastLeft = False
        for coord in coords[1:]:
            if coord[0] > 20:
                if not lastLeft:
                    right.append((0,coord[1]))
                left.append(coord)
                lastLeft = True
            else:
                if lastLeft:
                    left.append((24,coord[1]))
                right.append(coord)
                lastLeft = False
        if left[0] != left[len(left)-1]:
            left.append(left[0])
        if right[0] != right[len(right)-1]:
            right.append(right[0])
        leftPoly = shapely.geometry.Polygon(left)
        rightPoly = shapely.geometry.Polygon(right)
        with open('sep/left-'+arg, 'w') as f:
            f.write(leftPoly.wkt)
        with open('sep/right-'+arg, 'w') as f:
            f.write(rightPoly.wkt)
