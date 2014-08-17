#!/usr/bin/env python

import sys

if len(sys.argv) < 2 :
    print("Usage: " + sys.argv[0] + " [input file]")
    sys.exit(1)

minMag = float("inf")
maxMag = float("-inf")
total = 0
count = 0
fname = sys.argv[1]
with open(fname) as src:
    for line in src:
        parts = line.split("\t")
        mag = float(parts[5])
        count += 1
        total += mag
        if mag > maxMag:
            maxMag = mag
        if mag < minMag:
            minMag = mag

av = float(total) / float(count)
print("max %f, min %f, count %f, av %f" % (maxMag, minMag, count, av))
