#!/bin/bash


for i in {1..500}
do
    w=$(shuf -i1-100 -n1)
    h=$(shuf -i1-100 -n1)

    r=$(shuf -i1-255 -n1)
    g=$(shuf -i1-255 -n1)
    b=$(shuf -i1-255 -n1)

    convert -size ${w}x${h} xc:rgba\(${r},${g},${b},0.5\) ${i}.png
done
