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

- `privilege-escalate`: write `/etc/passwd` to make local root escalation easy.
- `image-pollution`: write into a file from a container image layer.
- `escape image`: generate a Docker build context for image-based runc capture.
- `escape exec`: wait for `docker exec` to capture and overwrite host runc.
- `escape restart`: trigger container restart to capture and overwrite host runc.

The package should stay generic. CVE-specific code belongs in the vulnerability
package and is injected through the `Primitive` interface or the escape image
provider interfaces below.

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
