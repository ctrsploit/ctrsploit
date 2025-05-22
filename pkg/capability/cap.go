package capability

import (
	"os"
	"strconv"
	"strings"
)

var (
	Caps22 = []string{
		"CAP_CHOWN",            // 2.2
		"CAP_DAC_OVERRIDE",     // 2.2
		"CAP_DAC_READ_SEARCH",  // 2.2
		"CAP_FOWNER",           // 2.2
		"CAP_FSETID",           // 2.2
		"CAP_KILL",             // 2.2
		"CAP_SETGID",           // 2.2
		"CAP_SETUID",           // 2.2
		"CAP_SETPCAP",          // 2.2
		"CAP_LINUX_IMMUTABLE",  // 2.2
		"CAP_NET_BIND_SERVICE", // 2.2
		"CAP_NET_BROADCAST",    // 2.2
		"CAP_NET_ADMIN",        // 2.2
		"CAP_NET_RAW",          // 2.2
		"CAP_IPC_LOCK",         // 2.2
		"CAP_IPC_OWNER",        // 2.2
		"CAP_SYS_MODULE",       // 2.2
		"CAP_SYS_RAWIO",        // 2.2
		"CAP_SYS_CHROOT",       // 2.2
		"CAP_SYS_PTRACE",       // 2.2
		"CAP_SYS_PACCT",        // 2.2
		"CAP_SYS_ADMIN",        // 2.2
		"CAP_SYS_BOOT",         // 2.2
		"CAP_SYS_NICE",         // 2.2
		"CAP_SYS_RESOURCE",     // 2.2
		"CAP_SYS_TIME",         // 2.2
		"CAP_SYS_TTY_CONFIG",   // 2.2
	}
	Caps24   = append(Caps22, "CAP_MKNOD", "CAP_LEASE")
	Caps2611 = append(Caps24, "CAP_AUDIT_WRITE", "CAP_AUDIT_CONTROL")
	Caps2624 = append(Caps2611, "CAP_SETFCAP")
	Caps2625 = append(Caps2624, "CAP_MAC_OVERRIDE", "CAP_MAC_ADMIN")
	Caps2637 = append(Caps2625, "CAP_SYSLOG")
	Caps30   = append(Caps2637, "CAP_WAKE_ALARM")
	// Caps35 is the caps of kernel 3.5 (37 entries)
	Caps35 = append(Caps30, "CAP_BLOCK_SUSPEND")
	// Caps316 is the caps of kernel 3.16 (38 entries)
	Caps316 = append(Caps35, "CAP_AUDIT_READ")
	// Caps58 is the caps of kernel 5.8 (40 entries)
	Caps58 = append(Caps316, []string{"CAP_PERFMON", "CAP_BPF"}...)
	// Caps59 is the caps of kernel 5.9 (41 entries)
	Caps59     = append(Caps58, "CAP_CHECKPOINT_RESTORE")
	CapsLatest = Caps59
)

// Supported returns the supported capabilities based on the current kernel version.
// Implementation steps:
// 1. Read the kernel version from "/proc/sys/kernel/osrelease".
// 2. Parse the version string (e.g., "5.9.0-xx-generic") to extract major, minor, and patch numbers.
// 3. Compare the current version with known versions and return the corresponding capabilities list.
// If reading or parsing fails, it defaults to returning Caps35.
func Supported() []string {
	// Try to read the kernel version from /proc/sys/kernel/osrelease.
	data, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		// Return Caps35 if reading fails.
		return Caps35
	}

	verStr := strings.TrimSpace(string(data))
	// For example, if verStr is "5.9.0-XX-generic", we extract "5.9.0".
	parts := strings.Split(verStr, "-")
	versionPart := parts[0]

	segments := strings.Split(versionPart, ".")
	if len(segments) < 2 {
		return Caps35
	}

	major, err := strconv.Atoi(segments[0])
	if err != nil {
		return Caps35
	}

	minor, err := strconv.Atoi(segments[1])
	if err != nil {
		return Caps35
	}

	patch := 0
	if len(segments) >= 3 {
		patch, _ = strconv.Atoi(segments[2])
	}

	// Helper function: compare if the current kernel version is greater than or equal to the target version.
	curVersionGE := func(tMajor, tMinor, tPatch int) bool {
		if major > tMajor {
			return true
		} else if major < tMajor {
			return false
		}
		if minor > tMinor {
			return true
		} else if minor < tMinor {
			return false
		}
		return patch >= tPatch
	}

	// Evaluate versions in descending order.
	if curVersionGE(5, 9, 0) {
		return Caps59
	} else if curVersionGE(5, 8, 0) {
		return Caps58
	} else if curVersionGE(3, 16, 0) {
		return Caps316
	} else if curVersionGE(3, 5, 0) {
		return Caps35
	} else if curVersionGE(3, 0, 0) {
		return Caps30
	} else if curVersionGE(2, 6, 37) {
		return Caps2637
	} else if curVersionGE(2, 6, 25) {
		return Caps2625
	} else if curVersionGE(2, 6, 24) {
		return Caps2624
	} else if curVersionGE(2, 6, 11) {
		return Caps2611
	} else if curVersionGE(2, 4, 0) {
		return Caps24
	}
	return Caps22
}
