# ctrsploit Asciinema Recording Notes

## Issue #291 Baseline

Record through tmux so the live demo can be driven separately from the recorder:

```bash
tmux new -s terminal-capture
asciinema rec -c "tmux attach -t terminal-capture" 0.cast
```

Recommended SVG conversion chain:

```bash
# Optional: normalize hostnames before conversion.
sed -i 's/asus-kali/localhost/g' 0.cast

# Convert v3/raw output to asciicast v2 for downstream tools.
asciinema convert --overwrite -f asciicast-v2 0.cast 1.cast

# Re-sort and cap long pauses.
sciine maximum 1 -i 1.cast -o 2.cast

# Example targeted edits from the issue. Adjust ranges for the new capture.
sciine maximum 0.5 -e 89 -i 2.cast -o 3.cast
sciine maximum 0.08 -s 7 -e 37 -i 3.cast -o 4.cast
sciine maximum 1 -s 38 -e 45 -i 4.cast -o 5.cast

# Convert final cast to SVG.
cat 5.cast | svg-term --out ../video.svg
```

GIF conversion uses the pinned agg container:

```bash
docker pull ghcr.io/asciinema/agg:1.7.0
docker run --rm -it -u "$(id -u):$(id -g)" -v "$PWD:/data" ghcr.io/asciinema/agg:1.7.0 demo.cast demo.gif
```

## Existing Artifact Patterns

- Most demos use `vul/<name>/video.svg`.
- Some demos keep intermediates under `vul/<name>/cast/0.cast`,
  `1.cast`, `2.cast`, `3.cast`.
- The `cve-2019-5736` demo has multiple SVG names (`checksec.svg`,
  `image.svg`, `exec.svg`), so copy local naming rather than forcing
  `video.svg`.
- Large captures can produce multi-megabyte SVG/GIF outputs; this is normal for
  long terminal demos.

## Practical Editing Heuristics

- Use `sciine maximum 1` globally first to remove accidental long pauses.
- Use targeted `-s` and `-e` ranges for command downloads, Docker pulls,
  package installation, and long idle waits.
- Keep exploit execution, proof output, and commands that teach the flow slower
  than mechanical setup steps.
- Convert only the final cleaned cast to SVG/GIF.
