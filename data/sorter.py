#!/usr/bin/env python

import sys

if len(sys.argv) < 3 :
    print("Usage: " + sys.argv[0] + " [input file] [pivot]")
    sys.exit(1)

fname = sys.argv[1]
pivot = float(sys.argv[2])
with open(fname) as src, open("out1.tsv", "w") as out1, \
        open("out2.tsv", "w") as out2:
    for line in src:
        parts = line.split("\t")
        mag = float(parts[5])
        if mag < pivot:
            out1.write(line)
        else :
            out2.write(line)
