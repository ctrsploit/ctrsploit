// Package kernelprivesc provides a reusable privilege-escalation stage for
// kernel exploits that produce an arbitrary kernel-memory write primitive
// (typically via ROP: a memcpy/memset gadget chain that writes controlled
// bytes to a controlled kernel virtual address).
//
// Each CVE implements KernelWritePrimitive (the "how do I write to kmem" part);
// this package owns the escalation methods (the "what do I overwrite to get
// root" part). Methods are loosely coupled: each overwrites one kernel global
// (a path string such as modprobe_path or core_pattern) and triggers the
// kernel to execute a root-shell-dropping script as root from userspace.
//
// This mirrors pkg/pipe-primitive, which does the same for file-write
// primitives (passwd-su, suid-overwrite). The strategy chain tries each
// method in order and returns the first success.
//
// # Known symbols
//
// Methods address kernel symbols by logical name via KernelWritePrimitive.Kaddr:
//
//   - "modprobe_path"  — unknown-format exec triggers modprobe as root
//   - "core_pattern"   — coredump pipe triggers script as root
//   - "selinux_state"  — zeroed to make selinux permissive (best-effort)
//   - "poweroff_cmd"   — orderly_poweroff (not in default methods; dangerous)
//
// # Future methods (not implemented)
//
//   - poweroff_cmd  — overwrites poweroff_cmd; trigger shuts the box down.
//   - commit_creds  — direct cred overwrite (commit_creds(&init_cred)) +
//     return-to-userland; no file artifact, but needs a KPTI trampoline
//     return chain. Different shape from the path-overwrite methods.
//
// # Parent/worker note
//
// Kernel LPEs often run inside a user namespace (CLONE_NEWUSER) to gain
// CAP_NET_ADMIN. A setuid-root binary exec'd from inside CLONE_NEWUSER maps
// uid 0 to the caller's host uid — not real root. DropRootShell must be
// called from the parent process (init user namespace) where the setuid bit
// is honored for real. The CVE's worker emits a success marker; the parent
// sees it and calls DropRootShell.
package kernelprivesc
