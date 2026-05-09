---
name: record-terminal-session
description: Record terminal demonstrations for ctrsploit vulnerabilities and convert asciinema .cast captures into project-style video.svg or video.gif artifacts. Use when asked to record, re-record, edit, clean, convert, or validate terminal demo assets under vul/**/cast/*.cast, vul/**/video.svg, or *.gif, following the workflow from ctrsploit issue #291.
---

# Record Terminal Session

## Goal

Create reproducible terminal demos for this repository using `tmux`, `asciinema`,
`sciine`, `svg-term`, and optionally `agg`, matching the existing `vul/**/video.svg`
and `vul/**/cast/*.cast` artifacts.

## Target Layout

- Work from the vulnerability or demo directory, usually `vul/<name>` or a nested
  directory such as `vul/shared-socket/docker-sock`.
- Store raw and intermediate asciinema files in `cast/`:
  `cast/0.cast`, `cast/1.cast`, `cast/2.cast`, ...
- Write the primary SVG next to the target directory as `video.svg`, unless the
  local pattern uses a specific name such as `exec.svg`, `image.svg`, or
  `checksec.svg`.
- For GIF output, place `video.gif` in `cast/` unless the user requests another
  path.
- When creating or updating `video.svg` for a vulnerability README, update that
  README to reference the SVG, usually with `![](./video.svg)` under
  `## 5. Reproduce`, matching nearby vulnerability docs.

## ctrsploit Vulnerability Demo Content

For `vul/**` demos, record a manual vulnerability exploitation flow in the
appropriate lab target, usually the dqd environment referenced by the
vulnerability's `README.md` or `e2e.yml`.

- The visible recording should show real user-facing steps: starting or entering
  the dqd lab, installing or uploading the current `ctrsploit` binary when
  needed, running `ctrsploit ... check` or `checksec`, running the exploit
  command, and verifying a concrete proof such as modified file content,
  created proof file, root shell, or `[Y]` result.
- Do not assume a source-tree or locally built `ctrsploit` binary should be
  used just because the task originates from this repository. First verify
  whether the installation path documented in the target `README.md` is
  sufficient for the demonstrated module and exploit flow.
- If the demo needs a locally built or unreleased `ctrsploit` binary, make that
  installation path explicit in the recording: show the upload/copy into the VM
  or container and the `chmod +x`/install step. Do not imply a release download
  when the release binary would not contain the demonstrated module.
- Prefer the installation command documented in the target README, such as a
  release `wget`, when it works for the demonstrated module. Use a local build
  upload only after explicitly verifying the release path is unsuitable. Record
  or note the concrete reason for the deviation, such as a missing module,
  unreleased feature, or behavior mismatch, and keep the local build/upload
  step visible in the recording.
- Do not use `make e2e` or test binaries as the main recorded content unless the
  user explicitly asks for an e2e/test demonstration. E2E is useful before or
  after recording to validate behavior, but it hides the exploit steps and is not
  the default demo narrative.
- Prefer a concise exploit path that proves the vulnerability clearly and avoids
  unnecessary noise. If several exploit modes exist, pick the shortest
  understandable proof unless the user specifies a mode.
- For destructive flows such as runc overwrite or host file modification, use
  only disposable lab VMs/containers and make the final proof line explicit.

## Readability Rules

The rendered SVG should look like a deliberate terminal demo, not a noisy raw
shell transcript.

- Keep visible commands within the chosen terminal width. At `width: 118`, avoid
  long one-line SSH/SCP/Docker commands that hard-wrap mid-word or split flags
  such as `ConnectTimeout=3`, file names, or image names across lines.
- Prefer short helper commands, shorter local file names, environment variables,
  shell aliases, or deliberate multi-line commands with `\` continuations over
  accidental terminal hard wrapping.
- Avoid recording dynamic progress output when it distracts from the exploit.
  For Docker commands, prefer cached images, `--quiet` where useful, or concise
  proof commands after the noisy operation completes.
- After rendering, inspect the SVG or a screenshot for wrapped words, split
  command lines, spinner/progress artifacts, local hostnames, and overlapping or
  visually garbled text. Re-record or edit the cast if the demo is hard to read.

## Workflow

1. Inspect nearby artifacts before recording:
   - List existing `.cast`, `.svg`, and `.gif` files in the target directory.
   - Read the first line of an existing `.cast` to copy dimensions, shell, and
     theme when possible.
   - Prefer dimensions already used by the target; common examples are
     `width: 118`, `height: 55` and SVG width `2130`.
   - Decide the `ctrsploit` install source before recording:
     start from the target `README.md` install command, usually a release
     download. Switch to a local build only after verifying the README path is
     unsuitable for the exact module and demo flow.
2. Start a tmux session for the demo:
   ```bash
   tmux new -s terminal-capture
   ```
3. In another terminal, record the tmux session:
   ```bash
   mkdir -p cast
   asciinema rec -c "tmux attach -t terminal-capture" cast/0.cast
   ```
4. Perform the demo inside tmux. For `vul/**` demos, follow the manual
   vulnerability demo content guidance above. Follow the readability rules:
   keep commands concise, avoid accidental hard wrapping, avoid leaking host
   secrets, and end on a stable proof line such as a `[Y]` result, expected file
   contents, shell prompt, or exploit success message.
5. Convert the raw capture to asciicast v2:
   ```bash
   asciinema convert --overwrite -f asciicast-v2 cast/0.cast cast/1.cast
   ```
6. Normalize and trim timing with `sciine maximum`. Use one or more passes; keep
   each generated file so the edit chain is reviewable:
   ```bash
   sciine maximum 1 -i cast/1.cast -o cast/2.cast
   sciine maximum 0.5 -i cast/2.cast -o cast/3.cast
   ```
7. Convert the final `.cast` to SVG:
   ```bash
   cat cast/3.cast | svg-term --out video.svg
   ```
8. If the SVG is intended for a vulnerability README, add or refresh the local
   README reference to it, usually `![](./video.svg)` under `## 5. Reproduce`.
9. Optionally render a GIF with `agg`:
   ```bash
   docker pull ghcr.io/asciinema/agg:1.7.0
   docker run --rm -it -u "$(id -u):$(id -g)" -v "$PWD/cast:/data" ghcr.io/asciinema/agg:1.7.0 3.cast video.gif
   ```
10. Validate outputs:
   - `head -n 1 cast/<final>.cast` shows asciicast v2 metadata.
   - `video.svg` begins with `<svg` and uses project-like terminal dimensions.
   - The relevant README references the generated SVG when project convention
     calls for it.
   - If the recording used a local build instead of the documented release
     install path, the reason for the deviation was verified and the upload or
     install step is visible in the demo.
   - The rendered demo does not include local-only hostnames unless intentionally
     shown.
   - Visible command lines are not accidentally hard-wrapped or visually
     garbled.
   - The final frames visibly show the vulnerability result or proof.

## Cleaning Casts

Use structured asciinema tools before manual editing. Do not hand-edit timing
JSON unless the tools cannot express the required trim.

- Replace noisy hostnames or paths before conversion if needed:
  ```bash
  sed -i 's/asus-kali/localhost/g' cast/0.cast
  ```
- Use `sciine maximum <seconds> -s <start> -e <end>` to speed up boring ranges
  while preserving important proof steps.
- Keep command typing readable. Use a smaller maximum such as `0.08` for long
  download or install segments, and `0.5` to `1` for ordinary waits.

For ctrsploit-specific examples copied from issue #291, read
`references/ctrsploit-asciinema.md`.

## Safety

- Do not delete existing recording artifacts unless the user explicitly asks.
- If overwriting `video.svg` or `cast/*.cast`, inspect current files first and
  preserve unrelated user changes.
- Avoid recording tokens, private registry names, private SSH targets, shell
  history, or unrelated workspace paths.
- When a command needs network or Docker access and sandboxing blocks it, request
  escalation with the exact command needed.
