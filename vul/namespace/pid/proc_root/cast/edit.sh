#!/usr/bin/env bash

# https://github.com/jdum/asciinema-scene
# https://github.com/marionebl/svg-term-cli

# edit manually
# replace hostname
sed -i s/asus-kali/localhost/g 0.cast
# convert to v2 format
asciinema convert --overwrite -f asciicast-v2 0.cast 1.cast
# resort after edit manually
sciine maximum 1 -i 1.cast -o 2.cast
sciine minimum 3 -s 33.83129 -i 2.cast -o 3.cast
# convert to svg
cat 3.cast | svg-term --out ../video.svg