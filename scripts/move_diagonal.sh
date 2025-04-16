#!/bin/bash

curl -X POST http://localhost:17000 -d $'reset\nwhite\nfigure 0.0 0.0\nupdate'

for i in {1..20}
do
    dx=$(echo "scale=2; 0.5 + $i * 0.01" | bc)
    dy=$(echo "scale=2; 0.5 + $i * 0.01" | bc)
    curl -X POST http://localhost:17000 -d "move $dx $dy"$'\nupdate'
    sleep 1
done
