#!/usr/bin/env python3

import json
import subprocess

def tags_contains(commit=None):
    command = ["git", "-P", "tag", "--sort", "taggerdate"]
    if commit is not None:
        command += ["--contains", commit]
    output = subprocess.check_output(command)
    tags = output.strip().split(b"\n")
    return tags

def minor_version(tag):
    version = tag.split(b"-")[0]
    if version.count(b".") > 1:
        return b".".join(version.split(b".")[:2])
    return version

def output():
    # tags = tags_contains("2b188cc1bb857a9d4701ae59aa7768b5124e262e")
    tags = tags_contains()
    while True:
        if len(tags) == 0:
            break
        tag = tags[0]
        minor = minor_version(tag)
        line = [tag] + [t for t in tags[1:] if minor_version(t) == minor]
        tags = [t for t in tags if t not in line]
        line = [i.strip(b"v").decode() for i in line]
        print(json.dumps(line)[1:-1], end=',\n')

output()