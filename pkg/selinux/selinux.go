package selinux

import (
	"os"
	"strings"

	"github.com/opencontainers/selinux/go-selinux"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type TypeMode int

func (m TypeMode) String() (mode string) {
	switch m {
	case -1:
		mode = "disabled"
	case 0:
		mode = "permissive"
	case 1:
		mode = "enforcing"
	default:
		mode = "unknown"
	}
	return
}

func KernelSupported() (supported bool, err error) {
	content, err := os.ReadFile("/proc/filesystems")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	supported = strings.Contains(string(content), "selinuxfs")
	return
}

// IsEnabled detect selinux enabled inside the container
func IsEnabled() bool {
	con, err := selinux.CurrentLabel()
	if err != nil {
		return false
	}
	if con == "kernel" {
		return false
	}
	if strings.Count(con, ":") >= 2 {
		return true
	}
	return false
}

func IsSelinuxPrivileged() bool {
	con, err := selinux.CurrentLabel()
	if err != nil {
		return true
	}
	return isSelinuxPrivileged(con)
}

var (
	knownPrivilegedCon = map[string]struct{}{
		"unconfined_t": {},
		"spc_t":        {},
		"sysadm_t":     {},
		"init_t":       {},
		"kernel_t":     {},
		"kernel":       {},
	}
)

// isSelinuxPrivileged reports whether the given SELinux context string
// represents a privileged or unconfined domain.
func isSelinuxPrivileged(con string) bool {
	if con == "kernel" {
		return true
	}
	parts := strings.Split(con, ":")
	if len(parts) < 3 {
		return true
	} else {
		t := parts[2]
		if _, ok := knownPrivilegedCon[t]; ok {
			return true
		}
		if strings.Contains(t, "unconfined") {
			return true
		}
	}
	return false
}

func Mode() TypeMode {
	return TypeMode(selinux.EnforceMode())
}
