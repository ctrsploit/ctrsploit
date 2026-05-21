---
name: record-terminal-session
description: Record terminal demonstrations for ctrsploit vulnerabilities and convert asciinema .cast captures into project-style SVG or GIF artifacts. Use when asked to record, re-record, edit, clean, convert, or validate terminal demo assets under vul/**/cast/*.cast, vul/**/*.svg, vul/**/video.svg, or *.gif, following the workflow from ctrsploit issue #291.
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
- If a vulnerability README demonstrates multiple exploit modes, checks, or
  trigger paths, prefer one named SVG per README subsection instead of forcing
  every path into a single `video.svg`. Use names that match the subsection,
  such as `checksec.svg`, `privilege-escalate.svg`, `image-pollution.svg`, or
  `escape-exec.svg`.
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
- Keep the README commands, the recording, and the verified manual run on the
  same logical path. Do not let the README show one installation or trigger
  method while the SVG shows another.
- Design the demo around the normal operator flow before choosing recording
  automation. A command sequence that is real but optimized for scripting is not
  necessarily a good public demo if a user would naturally perform different
  steps.
- Do not redirect the primary exploit output or proof output into a temporary
  log and then `cat` that log just to make recording easier. The main exploit
  command should usually print its own result directly in the visible terminal.
  Visible log files are appropriate only when reading a log is itself the
  natural operator workflow, such as inspecting a background service.
- Do not present retry or polling loops such as `for`, `seq`, or `until` as the
  public reproduction path unless that loop is genuinely part of the normal user
  workflow. If a loop or wrapper is needed just to make the recording pass,
  first revisit the lab image, shell, dynamic libraries, entrypoint, mounts, and
  trigger choice.
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
- Match the install location to where the exploit actually runs. If the exploit
  is demonstrated inside a container, prefer downloading or copying `ctrsploit`
  into that container and creating input files there. Do not mount host-side
  tools or source files into the container unless that is the realistic,
  documented setup being demonstrated.
- For container-internal exploit demos, prefer showing the natural container
  workflow when that is what a user would do: start or enter the container,
  show the container prompt, prepare input inside the container, run
  `ctrsploit` inside the container, then exit and verify from a fresh container
  or the host as appropriate. Do not collapse those steps into a host-side
  helper script or a long `sh -c` wrapper just because it is easier to record,
  unless the README intentionally documents that wrapper as the user workflow.
- If setup steps are part of the proof's causal chain, make them visible. For
  example, show the command that creates a polluted state before demonstrating a
  cleanup command. Hide only pure environment preparation that a reader does not
  need to understand the result.
- Do not hide key exploit steps in temporary files such as `/tmp/*.sh` for the
  visible demo. Temporary helpers are acceptable for local recording automation,
  but the rendered terminal should show the commands a normal operator would
  type and the real output those commands produced in the lab.
- Do not use `make e2e` or test binaries as the main recorded content unless the
  user explicitly asks for an e2e/test demonstration. E2E is useful before or
  after recording to validate behavior, but it hides the exploit steps and is not
  the default demo narrative.
- Prefer a concise exploit path that proves the vulnerability clearly and avoids
  unnecessary noise. If several exploit modes exist, pick the shortest
  understandable proof unless the user specifies a mode.
- If the concise path is unreliable, fix the reproduction environment or
  documented flow before considering product-code changes. Do not patch core
  exploit or watcher logic only to make a demo easier unless the investigation
  confirms an actual product bug.
- For waiting or trigger-based exploits, record the full causal chain:
  setup, waiting state, trigger command, post-trigger exploit output, and final
  proof. Do not stop at a `Waiting...` message, and do not show the final
  validation before the primary exploit pane has visibly completed.
- For trigger-based exploits, choose a trigger command that is understandable to
  the target user, returns controllably, and avoids misleading failure noise. If
  the first obvious trigger hangs or becomes noisy after the exploit changes the
  runtime, find an equivalent direct trigger or adjust the documented flow; do
  not hide the problem behind log redirection or a recording-only wrapper.
- For destructive flows such as runc overwrite or host file modification, use
  only disposable lab VMs/containers and make the final proof line explicit.
- For interactive shell exploits, especially local privilege escalation flows,
  do not pipe scripted input into the exploit command as the visible demo, such
  as `printf 'id\n...' | ctrsploit ...`. Record a real interactive session:
  run the exploit, wait for the new shell, then visibly enter proof commands
  such as `id` or `head -1 /etc/passwd`.

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
- Do not reuse one fixed terminal height for every SVG. Render each SVG with a
  height that fits its content. For short subsection demos, crop to the smallest
  useful height; for larger demos, increase height only as needed to preserve
  the key commands and proof.
- For multi-terminal flows, prefer a single tmux recording with split panes over
  separate SVGs unless the user asks for separate files. Size pane ratios to the
  content, giving the main exploit or waiting pane enough rows to avoid hiding
  key output.
- In split-pane recordings, align the timing so the reader can follow cause and
  effect: trigger appears after the waiting state, post-trigger output remains
  visible long enough to read, and proof commands run after the exploit visibly
  completes.
- Pace scripted or headless demos like a human demonstration. Leave about
  `0.8s` to `1.2s` between visible commands, pause briefly after important
  output before typing the next proof command, and avoid a burst of commands
  landing in the same instant.
- Leave enough dwell time on the final proof line and final prompt for a reader
  to understand the result.
- Quantify the final-frame dwell time for generated SVGs instead of only
  checking the exported text. Parse the SVG `animation-duration` and final
  keyframe percentage, and require the last proof/prompt frame to remain visible
  for at least 3 seconds, preferably around 4 to 5 seconds for proof output.
- For headless or scripted recordings, do not let the command exit immediately
  after printing the final prompt. Add an explicit final wait or adjust the cast
  exit event, and ensure `idle_time_limit` does not compress that wait away.
- When using `expect` to hide SSH banners in headless recordings, keep
  `log_user 0` only until the first prompt is matched. Then enable `log_user 1`
  and send an empty carriage return so the first visible frame contains a real
  prompt before the demo commands start.
- Keep recording mechanics out of the visible transcript. When using `expect`
  with SSH or container shells, perform connection setup, `TERM` changes, shell
  startup flags, and `PS1` initialization while `log_user 0` is active. The
  rendered demo should start at a clean prompt followed by the first user-facing
  command, not `spawn ssh`, `PS1=...`, or similar automation details.
- Suppress shell and terminal control noise before recording. If SSH or an
  interactive shell emits title escapes, profile banners, or prompt decoration,
  use a quiet shell setup such as `TERM=dumb`, `--noprofile --norc`, and a fixed
  `PS1` before enabling visible output.
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
   - If the README has multiple reproduction subsections, decide whether each
     subsection needs its own named SVG before recording.
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
   Use the final cast dimensions and rendered content to choose SVG height; do
   not blindly reuse a fixed height from another demo if it cuts off content or
   leaves excessive empty space.
8. If the SVG is intended for a vulnerability README, add or refresh the local
   README reference to it, usually near the subsection it demonstrates. For a
   single-flow README this is often `![](./video.svg)` under `## 5. Reproduce`;
   for multi-flow READMEs, place each named SVG next to the matching subsection.
9. Optionally render a GIF with `agg`:
   ```bash
   docker pull ghcr.io/asciinema/agg:1.7.0
   docker run --rm -it -u "$(id -u):$(id -g)" -v "$PWD/cast:/data" ghcr.io/asciinema/agg:1.7.0 3.cast video.gif
   ```
10. Validate outputs:
   - `head -n 1 cast/<final>.cast` shows asciicast v2 metadata.
   - `video.svg` begins with `<svg` and uses project-like terminal dimensions.
   - The relevant README references the generated SVG when project convention
     calls for it. For multi-flow READMEs, each recorded exploit mode is
     referenced near the subsection it demonstrates.
   - The README, recording, and real verified command path agree on install
     location, input-file location, trigger command, and proof command.
   - Every important visible command result was obtained by actually running the
     same command in the target lab; do not synthesize proof output, copy it
     from a different flow, or record a helper path while documenting a manual
     path.
   - The visible command order follows the target user's normal mental model.
     For example, a container exploit demo should not hide container-internal
     setup in a host-side helper script when a normal user would enter the
     container shell and run the commands there.
   - If the recording used a local build instead of the documented release
     install path, the reason for the deviation was verified and the upload or
     install step is visible in the demo.
   - If the demo uses multiple terminals, the SVG is a tmux split-pane recording
     unless there is a clear reason to use separate files.
   - Split panes have enough visible rows for their roles; the main exploit pane
     does not lose key output because of an even but unsuitable split.
   - Waiting or trigger-based demos show setup, waiting, trigger, post-trigger
     output, and proof in a readable order.
   - The rendered demo does not include local-only hostnames unless intentionally
     shown.
   - Visible command lines are not accidentally hard-wrapped or visually
     garbled.
   - The visible cast events do not contain recording failures such as
     `send: spawn id not open`, `Operation not permitted`,
     `Bad owner or permissions`, unexpected expect debug output, or timeout
     messages. If they do, re-record before rendering the SVG.
   - The visible cast events and SVG do not leak recording setup such as
     `spawn ssh`, `PS1=...`, `TERM=...`, shell startup flags, SSH banners, or
     terminal-title escape sequences unless those are intentionally part of the
     user-facing demo.
   - Header-only text such as an `expect` script is not treated as visible
     terminal output. Validate the actual output events and rendered SVG text
     separately, including checks that unwanted demo shortcuts such as `printf`
     do not appear.
   - Scan the final `.cast` and SVG text for accidental shortcuts and failure
     noise before finishing. For exploit demos, this usually includes negative
     checks for retry loops (`for i in`, `seq`, `until`), hidden helper scripts
     used as the visible flow (`/tmp/*.sh`), terminal detach/control artifacts,
     log relay patterns such as `>/tmp/*.log` followed by `cat /tmp/*.log`, and
     errors such as `No such file`, `docker: Error`, or dynamic loader failures.
   - Review visible redirections to temporary logs. They should be part of the
     user's natural diagnostic workflow, not a way to hide primary exploit
     output, trigger output, or proof output from the terminal.
   - Use distinct, causal names for any visible log redirections so readers can
     tell which command produced which output. Avoid reusing one generic
     `trigger.log` for unrelated trigger and proof steps.
   - The final frames visibly show the vulnerability result or proof and stay on
     screen long enough to read.
   - The final SVG frame dwell has been measured, not guessed. For SVGs generated
     by `svg-term`, compute dwell from `animation-duration` and the final
     keyframe percentage; fix any value under 3 seconds before finishing.
   - Headless browser screenshots of `svg-term` animations can capture the
     initial blank frame. Do not judge the SVG solely from that screenshot; also
     inspect the SVG text, keyframes, dwell time, or open it interactively.

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
