#!/usr/bin/env python

import sys
import json

if len(sys.argv) > 1:
    for fname in sys.argv[1:]:
        obj = None
        with open(fname,'r') as f:
            obj = json.load(f)
            polys = []
            for wkt,label in zip(obj['WktFiles'],obj['LabelPoints']):
                poly = {'WktFile':wkt,'LabelPoint':label,'MaxScale':0.012}
                polys.append(poly)
            del obj['WktFiles']
            del obj['LabelPoints']
            obj['Polys'] = polys
        with open(fname, 'w') as f:
            json.dump(obj, f, indent=4)
else :
    print("Too few arguments")
