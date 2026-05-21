---
name: dqd-lab
description: "Manage ctrsploit dqd lab environments from github.com/ctrsploit/dqd. Use when the user asks to start, stop, inspect, or connect to a dqd or dqd_dir experiment environment, run or prepare ctrsploit e2e.yml entries with kind: dqd, verify SSH access to dqd-* hosts, inspect Docker Compose KVM lab status or logs, or tear down disposable ctrsploit test VMs."
---

# DQD Lab

## Overview

Use this skill to manage disposable ctrsploit experiment VMs backed by the dqd repository. Prefer the repository's existing e2e flow when the user wants tests. Prefer the dqd repository's native documented commands when the lab workflow is user-facing, recorded, or intended for README reproduction; use the bundled lifecycle script only as a convenience for local, non-recorded lab management.

## Safety Rules

- Treat dqd machines as disposable lab infrastructure. Never aim these workflows at production hosts or user machines outside the dqd SSH aliases.
- Do not run exploit payloads just because a lab is up. Start and verify the environment first; execute destructive tests only when the user requested that specific test or exploit path.
- Expect network, Docker, KVM, SSH config, and image pulls. If a required command fails because of sandbox or network restrictions, rerun with the appropriate approval instead of inventing a workaround.
- Preserve user changes in `e2e.yml`, compose files, and dqd checkouts. Inspect before editing.

## Choose The Flow

- User asks to record a terminal demo, generate `video.svg`, or reproduce README-visible steps: follow the target README/dqd project commands in the visible session, usually `git clone`, `cd dqd/<dqd_dir>`, `docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d`, and `./ssh`. Do not substitute `scripts/dqd-lab.sh` into a recording unless the user explicitly asks to demonstrate that helper.
- User asks to run ctrsploit e2e tests for a module: use the project flow, usually `TEST_ENV=<name> make e2e DIR=<module-dir>` or `make e2e DIR=<module-dir>`.
- User asks to start a lab from `dqd_dir`, inspect logs/status, SSH into it, or stop it without requesting a recorded or README-faithful flow: using `scripts/dqd-lab.sh` is acceptable for operator convenience.
- User gives a module path with `e2e.yml`: read the file, select the matching `test_envs[]` item, and use its `dqd_dir`, `remote_host`, `start_timeout`, and `name`.
- User gives only a dqd path such as `docker/v28.3.2` or `vul/cve-2026-31431`: start that dqd directory. Ask for a host only if SSH or remote execution is required and it cannot be inferred from `e2e.yml`.

## DQD Hostnames

All direct SSH and SCP commands to dqd lab aliases must use the canonical
`dqd-<host>` hostname. Values in local config, such as `remote_host:
cve-2026-31431` or an environment name, are not necessarily SSH aliases by
themselves; normalize them before connecting. If the value already starts with
`dqd-`, use it as-is. Otherwise prefix `dqd-`, for example:

```bash
ssh dqd-cve-2026-31431 'hostname'
scp ./ctrsploit_linux_amd64 dqd-cve-2026-31431:/usr/bin/ctrsploit
```

When using the native recorded dqd flow from inside the dqd checkout, `./ssh`
is also acceptable because it is part of the dqd project interface. Do not
replace a missing `dqd-` SSH alias with the raw unprefixed host name.

For scripted or headless terminal recordings, prefer the dqd SSH config
directly to avoid unrelated local or system SSH config noise:

```bash
ssh -tt -F /tmp/dqd/ssh_config/config -o LogLevel=ERROR dqd-cve-2026-31431
```

Adjust only the host name for the target lab.

## Native DQD Lifecycle

For recorded demos and documentation reproduction, use the native commands from the dqd checkout:

```bash
git clone https://github.com/ctrsploit/dqd.git
cd dqd/vul/cve-2026-31431
docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
./ssh
```

For local preparation when the checkout already exists, run the same compose command from `${DQD_ROOT:-/tmp/dqd}/<dqd_dir>` rather than showing a helper wrapper in the recording.

## Helper Lab Lifecycle

For non-recorded operator tasks, run these from the skill directory or pass the script by absolute path:

```bash
scripts/dqd-lab.sh up vul/cve-2026-31431 --timeout 120
scripts/dqd-lab.sh ps vul/cve-2026-31431
scripts/dqd-lab.sh logs vul/cve-2026-31431 --tail=120
scripts/dqd-lab.sh ssh dqd-cve-2026-31431 -- 'hostname && docker version --format "{{.Server.Version}}"'
scripts/dqd-lab.sh down vul/cve-2026-31431
```

The script uses `${DQD_ROOT:-/tmp/dqd}` and `${DQD_REPO:-https://github.com/ctrsploit/dqd.git}`. On first `up`, it clones dqd and runs `script/install_ssh_config.sh`; later starts update the checkout unless `--no-update` is passed. Treat this script as local automation, not as the canonical command sequence for public demos.

## ctrsploit E2E Flow

When an `e2e.yml` entry has `kind: dqd`, the repository already knows how to build the test binary, start dqd, upload `/ctrsploit.test`, export `TEST_ENV`, run the configured command, and tear down the compose stack.

Use:

```bash
make e2e DIR=vul/cve-2026-31431
TEST_ENV=cve-2026-31431-image-pollution make e2e DIR=vul/cve-2026-31431
```

Before running, inspect the selected `e2e.yml` and report the environment name, `dqd_dir`, `remote_host`, package, and test command. If `remote_host` is absent or null, the e2e runner uses the environment name. For manual SSH or SCP, normalize the chosen host with the `dqd-` prefix before connecting.

## Verification

After starting a lab, verify at least one of:

- Native flow: from the dqd lab directory, `docker compose -f docker-compose.yml -f docker-compose.kvm.yml ps` shows running services.
- Native flow: from the dqd lab directory, `docker compose -f docker-compose.yml -f docker-compose.kvm.yml logs --tail=80` contains the configured boot success string, by default `Reached target multi-user.target`.
- `scripts/dqd-lab.sh ps <dqd_dir>` shows running compose services.
- `scripts/dqd-lab.sh logs <dqd_dir> --tail=80` contains the configured boot success string, by default `Reached target multi-user.target`.
- `scripts/dqd-lab.sh ssh dqd-<host> -- 'hostname'` succeeds.

If SSH fails right after boot, check compose logs and retry after a short wait before changing configuration.
