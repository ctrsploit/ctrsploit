# Kernel Privesc

## Overview

`pkg/kernel-privesc` is the shared privilege-escalation stage for vulnerabilities
that produce an **arbitrary kernel-memory write primitive** — typically a ROP
chain of `memcpy`/`memset` gadgets that writes controlled bytes to a controlled
kernel virtual address. It is the kernel-memory analogue of
`pkg/pipe-primitive`, which does the same job for file-write primitives.

The split between the vulnerability package and this package is:

- **The CVE** implements `KernelWritePrimitive` — "how do I write bytes to a
  kernel address on *this* kernel". This is the ROP/pivot-specific part.
- **This package** owns the escalation `Method`s — "which kernel global do I
  overwrite to get root, and how do I trigger it from userspace". This is the
  generic, kernel-version-agnostic part.

Methods are loosely coupled: each overwrites one kernel global that holds a path
string (`modprobe_path`, `core_pattern`, …), then triggers the kernel to execute
a root-shell-dropping script as root from userspace. The strategy chain tries
each method in order and returns the first success.

## The primitive

```go
type KernelWritePrimitive interface {
    ExpName() string
    WriteKmem(kaddr uint64, content []byte) error      // memcpy(kaddr, &content, len)
    MemsetKmem(kaddr uint64, fill byte, n uint64) error // memset(kaddr, fill, n)
    Kaddr(symbol string) (uint64, error)                // resolve "modprobe_path" etc.
}
```

`WriteKmem`/`MemsetKmem` **arrange** a ROP segment; they return immediately. The
actual write happens later, when the CVE's pivot thread runs the ROP. The
implementation must keep `content` alive and kernel-readable (e.g. embedded in
the `msg_msg` body) until the pivot fires, and must know the absolute address of
that buffer so it can back-patch each `memcpy`'s `&src` operand.

`Kaddr` resolves logical symbol names to runtime virtual addresses
(`kbase + offset`). Unknown symbols return an error so a method can degrade
gracefully — e.g. skip SELinux zeroing on a non-SELinux kernel.

## Methods

```go
type Method interface {
    Name() string
    Available(p KernelWritePrimitive) (bool, error)
    Prepare(p KernelWritePrimitive) (ShellArtifact, error)
    TriggerAndWait(p KernelWritePrimitive, a ShellArtifact) error
}
```

`Prepare` records the `WriteKmem`/`MemsetKmem` segments and writes the
root-shell-dropping script to disk. `TriggerAndWait` fires the userspace trigger
and polls for the setuid-root shell. The split lets a CVE prepare the writes
*before* a kernel pivot and trigger *after* it — necessary when the pivot thread
blocks (e.g. parks in `msleep`) and the trigger must run on a sibling thread.

| Method | Symbol | Userspace trigger | Notes |
|---|---|---|---|
| `ModprobePathMethod` | `modprobe_path` | exec an unknown-format binary | kernel runs `modprobe_path` as root to "load the module" |
| `CorePatternMethod` | `core_pattern` | child raises `SIGSEGV` | kernel pipes the coredump to `\|/path %P` as root |

Both also `MemsetKmem(selinux_state, 0, 4)` to make SELinux permissive
(best-effort; skipped if the symbol is unknown).

`DefaultMethods()` returns `[ModprobePathMethod, CorePatternMethod]`.
`SelectMethod(name)` picks one by name (used by a consuming CVE's `--method`
flag, which the CVE bridges to its namespace worker via an env var of its own
choosing).

There is deliberately **no `Escalate` entry point** that runs `Prepare` and
`TriggerAndWait` back-to-back. With this package's deferred-write primitive
model, `WriteKmem` only *arranges* a ROP segment — the kernel global is not
overwritten until the CVE's pivot fires the ROP. A back-to-back call would
trigger the userspace action against the real (unmodified) global and silently
fail to get root. The CVE must orchestrate the split itself: `Prepare` → fire
pivot → `TriggerAndWait` on a sibling thread.

### Future methods (documented, not in `DefaultMethods`)

- `poweroff_cmd` — overwrites `poweroff_cmd`; trigger is `orderly_poweroff`.
  Excluded by default — it shuts the box down.
- `commit_creds` — direct cred overwrite (`commit_creds(&init_cred)`) plus a
  return-to-userland chain. Different shape (no file artifact); needs a KPTI
  trampoline. Separate effort.

## Shell artifacts

`ShellArtifact` bundles the PID-suffixed `/tmp` paths a method uses:

- `ScriptPath` — the script the kernel runs as root (written into the overwritten
  global / pipe target).
- `RootBashPath` — the setuid-root `bash` copy the script drops.
- `MarkerPath` — a root-owned marker file (proves root ran the script; fallback
  to the setuid-bit check).

The script body is `cp /bin/bash $ROOTBASH; chmod 4755 $ROOTBASH; echo … >
$MARKER; chmod 644 $MARKER`. PID suffixes prevent collisions across concurrent
runs and stale root-owned files from prior runs.

`DropRootShell(path, in, out, err)` execs the setuid-root bash with `-p`
(preserves `euid=0`) via `exec.Command` (fork+exec, **not** `syscall.Exec`, so
the caller survives the shell and can still clean up).

## Parent / worker namespace split

Kernel LPEs frequently run inside a user namespace (`CLONE_NEWUSER`) to gain
`CAP_NET_ADMIN`. Two subtleties govern where the shell is dropped:

1. A setuid-root binary exec'd **inside** `CLONE_NEWUSER` maps uid 0 to the
   caller's host uid — not real root. `DropRootShell` must be called from the
   **parent** process (init user namespace), where the setuid bit is honored for
   real.
2. The pivoting commit thread often parks in `TASK_UNINTERRUPTIBLE` sleep
   (`msleep(0xffffffff)`), so the worker process stays in D state and
   `proc.Wait()` never returns. The parent cannot wait for worker exit before
   dropping the shell.

The proven pattern: the worker emits a success marker (`<expname>-lpe-ok <rootbash>`)
on a pipe; the parent streams that pipe, parses the rootbash path off the marker
line, drops the shell, and best-effort kills the stuck worker. Parsing the path
(not globbing `/tmp`) ensures the parent drops the exact method's shell even when
methods use different path prefixes.

## Known symbols

Methods address symbols by logical name via `Kaddr`:

| Name | Used by | Notes |
|---|---|---|
| `modprobe_path` | `ModprobePathMethod` | path the kernel execs as root on unknown-format bin |
| `core_pattern` | `CorePatternMethod` | `\|/path %P` pipes the coredump to the script as root |
| `selinux_state` | all methods | `.enforcing` zeroed to make SELinux permissive |
| `poweroff_cmd` | (future) | `orderly_poweroff` target; dangerous |

## Implementing a new kernel-write CVE

1. Implement `KernelWritePrimitive` in the vulnerability package. `WriteKmem`
   appends a `memcpy` ROP segment; `MemsetKmem` appends a `memset` segment;
   `Kaddr` switches on the symbol name and returns `kbase + offset`. Keep the
   written content in a controlled, kernel-readable buffer and back-patch each
   `memcpy`'s `&src` once the buffer's absolute address is known.
2. In the CVE's hijack/orchestrator, build the primitive, call
   `method.Prepare(p)` to record the kmem writes, fire the pivot, then call
   `method.TriggerAndWait(p, art)` on a sibling thread.
3. In the parent (init user ns), watch for the `…-lpe-ok <rootbash>` marker and
   call `DropRootShell`.

Keep this package generic: kernel-version-specific offsets, gadget addresses,
and pivot mechanics belong in the vulnerability package, injected through
`KernelWritePrimitive`.

## ROP layout invariant

When a CVE lays the ROP segments into a single buffer (e.g. a `msg_msg` body),
the park-forever tail (or any qwords appended after the `memcpy` segments) **must
be accounted for before back-patching each `memcpy`'s `&src`**. The source
pointer is `fakeRuleAddr + ropBytes + contentOff`, where `ropBytes` is the length
of *all* qwords including the tail. If the tail is appended after the `&src`
patch, every source pointer lands short and `memcpy` reads ROP gadget addresses
into the target global instead of the path string — the trigger fires but the
root shell never appears, with **no oops** (silent failure). Register the tail
inside the layout pass, not after it.
