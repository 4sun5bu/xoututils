#!/bin/bash

# Extract, rename and convert object files
xarch $1
for f in *.o ; do
    mv $f ${f%.*}.rel
    xout2coff ${f%.*}.rel
done

# Rename the original lib file to preserve
mv  $1 $1.xout

# Make a COFF lib file
z8k-coff-ar r $1 *.o

