# Pipe Primitive

## Overview

`pkg/pipe-primitive` contains shared exploit flows for vulnerabilities that
provide a file write primitive through a pipe-like kernel path. In this package,
a primitive is the ability to write bytes into an existing file at a given
offset without opening that file writable.

Each vulnerability package supplies the primitive-specific implementation:

```go
type Primitive interface {
	GetExpName() string
	MinOffset() int64
	Write(path string, offset int64, content []byte) error
}
```

This package owns the generic CLI and orchestration built on top of that
primitive:

- `privilege-escalate`: obtain a local root shell using `/etc/passwd` plus
  `su`, or a root-owned SUID helper fallback when `su` is unavailable.
- `image-pollution`: write into a file from a container image layer.
- `clean`: drop host page cache to clear page-cache-only primitive effects.
- `escape image`: generate a Docker build context for image-based runc capture.
- `escape exec`: wait for `docker exec` to capture and overwrite host runc.
- `escape restart`: trigger container restart to capture and overwrite host runc.

The package should stay generic. CVE-specific code belongs in the vulnerability
package and is injected through the `Primitive` interface or the escape image
provider interfaces below.

## Privilege Escalation

`privilege-escalate` first checks for a usable SUID `su` implementation,
including BusyBox `su`. Only after that preflight succeeds does it patch the
root password field in `/etc/passwd` and invoke `su` for a root shell.

If no `su`-compatible helper is available, the command tries a fallback based
on other root-owned SUID executables such as `chfn`, `chsh`, `gpasswd`,
`mount`, `newgrp`, `passwd`, `sudo`, or `umount`. That fallback overwrites the
helper's page cache from offset `0` with an embedded linux/amd64 SUID shell
payload and then executes the helper path. It requires a primitive with
`MinOffset() == 0`, a root-owned executable with the SUID bit set, and an
environment that allows SUID execution.

The fallback avoids modifying `/etc/passwd` when `su` is missing, but it is
still destructive to the selected helper's cached executable image. Run it only
in disposable labs or on systems where that side effect is acceptable.

## Clean

`clean` writes `3` to `/proc/sys/vm/drop_caches`, which asks the kernel to drop
page cache plus reclaimable dentry and inode slab objects. This can clear the
visible effect of pipe-primitive exploits that only polluted clean page-cache
entries.

This is not a file restore or rollback operation. It does not undo payloads
that already executed, files that were genuinely written to disk, or other side
effects such as proof files created by an overwritten runc payload.

Use `--sync` to call `sync` before dropping caches. This can make more clean
objects eligible to drop, but it may also flush unrelated dirty data and has a
system-wide performance impact. The command requires permission to write the
host's `/proc/sys/vm/drop_caches`.

## Escape Restart Loader

`escape restart` prepares a running victim container for the next restart. It
writes a small ELF over the image dynamic loader path, overwrites the victim's
current PID 1 executable with `#!/proc/self/exe`, and then runs configurable
restart triggers from `pkg/crash`. The default trigger chain tries cgroup v2
`cgroup.kill` first, then falls back to `SIGKILL` for PID 1 and `kill-all` for
child-process-driven entrypoints. Use `--restart-method all` to also try OOM
after those methods. When Docker restarts the container, runc becomes the script
interpreter and the fake loader captures `/proc/self/exe` as a host runc fd.

The preferred integration is a self-contained restart loader:

```go
type RestartLoaderProvider interface {
	RestartLoader(payload []byte) ([]byte, error)
}
```

The returned loader should contain the primitive-specific write path and the
requested payload. This path does not require special helper files in the
victim image.

Because this path relies on the image dynamic loader, it is intended for
dynamically linked runc.

## Escape Image Writer

`escape image` generates a Docker build context in one of two modes:

- `start` (default): capture runc start by replacing the image dynamic loader.
  This is intended for dynamically linked runc.
- `exec`: capture runc exec by running a long-lived watcher as the entrypoint
  and using Docker `HEALTHCHECK` to trigger recurring runc exec events. This is
  intended for statically linked runc. The healthcheck command executes
  `/proc/self/exe` directly so the watcher can capture the host runc fd before
  the exec turns into a normal container process.

`escape image --mode start` generates these fixed files:

- `Dockerfile`
- `ld.go`
- `writer.go`
- `payload`

`escape image --mode exec` generates these fixed files:

- `Dockerfile`
- `runc-capture.go`
- `writer.go`
- `payload`

`writer.go` is primitive-specific. It is provided by implementing
`EscapeImageWriterProvider`:

```go
type EscapeImageWriterProvider interface {
	EscapeImageWriter() []byte
}
```

The writer runs inside the generated image. Its job is to write `/payload` into
the target path passed as `argv[1]`, usually by calling the primitive's write
implementation.

### Writer Source Files

Prefer keeping writer sources as real `.go` files under the primitive package,
for example:

```go
//go:build ignore
// +build ignore

package main
```

The `ignore` build tag keeps the file out of normal repo builds and tests. This
matters because writer sources often import packages that only exist in the
generated Docker build context.

Before returning the writer from `EscapeImageWriter()`, strip the build tag so
the generated context contains a normal compilable `writer.go`.

### Extra Build Context Files

If `writer.go` imports helper packages, the primitive should also implement
`EscapeImageExtraFileProvider`:

```go
type EscapeImageExtraFileProvider interface {
	EscapeImageExtraFiles() map[string][]byte
}
```

The map key is the relative path inside the generated Docker build context.
Typical files are:

- `go.mod`
- helper package sources, such as `primitive/write.go`

Extra files must use relative paths and cannot overwrite the fixed generated
files. This keeps `pkg/pipe-primitive` generic: the package owns the image
layout and runc loader, while each primitive owns its writer and helper code.
