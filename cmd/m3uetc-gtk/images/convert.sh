#!/usr/bin/env bash

for s in 512x512 256x256 128x128 96x96 64x64 48x48 32x32 24x24 16x16
do
    mkdir -p ./$s
    convert ./1024x1024/m3u-etcetera.png -resize $s ./$s/m3u-etcetera.png
done
