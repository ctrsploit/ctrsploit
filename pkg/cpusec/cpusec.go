// Package cpusec detects active CPU and kernel security mitigations for the
// running architecture. It is the shared, reusable replacement for the
// private cpuHasSMEPorSMAP() helper that previously lived in the cve-2026-23111
// ret2usr path, and is consumed by the `env cpusec` subcommand and by
// per-vulnerability prerequisite gates.
//
// All detection is unprivileged userspace file reads: /proc/cpuinfo for
// runtime-active CPU flags, and /boot/config-<release> (fallback
// /proc/config.gz) for compile-time CONFIG_* knobs. No syscalls, no root.
package cpusec

import (
	"fmt"
	"runtime"
)

// Mitigations is the set of CPU/kernel hardening features active on the
// running system. Fields are grouped by architecture; only the group
// matching Arch is populated by Detect.
type Mitigations struct {
	// Arch is runtime.GOARCH ("amd64", "arm64", ...).
	Arch string
	// ConfigSource is where the kernel config was read from
	// ("/boot/config-<release>", "/proc/config.gz"), or "" if no config
	// source was available (cpuinfo-only detection).
	ConfigSource string

	// x86-64 mitigations.
	SMEP     bool // Supervisor Mode Execution Prevention
	SMAP     bool // Supervisor Mode Access Prevention
	KPTI     bool // Kernel Page Table Isolation (x86: "pti" flag)
	IBT      bool // Indirect Branch Tracking (CET)
	KCFI     bool // Clang Control-Flow Integrity (compile-time)
	FgKaslr  bool // Function Granular KASLR (compile-time)

	// arm64 mitigations.
	PAC bool // Pointer Authentication (paca/pacg flags / CONFIG_ARM64_PTR_AUTH)
	BTI bool // Branch Target Identification
	PAN bool // Privileged Access Never
	MTE bool // Memory Tagging Extension

	// CPUFlags is the raw flags (x86) / Features (arm64) tokens from the
	// first /proc/cpuinfo entry, kept for debugging/audit. Empty if cpuinfo
	// could not be read.
	CPUFlags []string
}

// Detect reads /proc/cpuinfo and the kernel config and returns the active
// mitigations for the running arch. It never returns an error for a missing
// kernel config (that degrades gracefully to cpuinfo-only detection); it only
// errors on a failure to read /proc/cpuinfo when the arch needs it.
func Detect() (m Mitigations, err error) {
	m.Arch = runtime.GOARCH
	flags, cpuinfoErr := parseCPUInfo()
	if cpuinfoErr != nil {
		// cpuinfo is the authoritative source for most x86/arm64 flags; a
		// read failure is a hard error for those archs. For unknown archs we
		// don't need it, so swallow the error there.
		if m.Arch != "amd64" && m.Arch != "arm64" {
			return m, nil
		}
		return m, fmt.Errorf("read /proc/cpuinfo: %w", cpuinfoErr)
	}
	m.CPUFlags = flags

	cfg, source, cfgErr := parseKernelConfig()
	if cfgErr == nil {
		m.ConfigSource = source
	}
	// A missing/unreadable kernel config is non-fatal: fall back to cpuinfo.

	switch m.Arch {
	case "amd64":
		fillX86(&m, flags, cfg)
	case "arm64":
		fillArm64(&m, flags, cfg)
	}
	return m, nil
}
