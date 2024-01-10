#!/usr/bin/env bash
runc=$1
gdb -q -ex "p library_version" -ex "quit" $runc | grep -oP '(?<=\s=\s)\d+' | tr '\n' '.' | sed 's/\.$/\n/'