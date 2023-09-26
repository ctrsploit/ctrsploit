package seccomp

import (
	"fmt"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
)

type Mode int

func (m Mode) String() (s string) {
	s = "unknown"
	switch m {
	case 1:
		// The first version of seccomp was merged in 2005 into Linux 2.6.12.
		// It was enabled by writing a "1" to /proc/PID/seccomp.
		// Once that was done, the process could only make four system calls: read(), write(), exit(),
		// and sigreturn().
		s = "strict"
	case 2:
		// Things were calm in seccomp land for the next five years or so until
		// "seccomp mode 2" (or "seccomp filter mode") was added to Linux 3.5 in 2012.
		// It added a second mode for seccomp: SECCOMP_MODE_FILTER. Using that mode,
		// processes can specify which system calls are permitted.
		s = "filter"
	}
	return
}

// CheckSupported
// borrowed from https://github.com/opencontainers/runc/blob/v1.0.0-rc91/libcontainer/seccomp/seccomp_linux.go#L86-L102
// https://github.com/kubernetes-sigs/security-profiles-operator/issues/76
// or on host: `grep CONFIG_SECCOMP= /boot/config-$(uname -r)` https://docs.docker.com/engine/security/seccomp/
// But it must be executed on the host, so the first two are preferred
func CheckSupported() bool {
	// Try to read from /proc/self/status for kernels > 3.8
	s, err := parseStatusFile("/proc/self/status")
	if err != nil {
		// Check if Seccomp is supported, via CONFIG_SECCOMP.
		if err := unix.Prctl(unix.PR_GET_SECCOMP, 0, 0, 0, 0); err != unix.EINVAL {
			// Make sure the kernel has CONFIG_SECCOMP_FILTER.
			if err := unix.Prctl(unix.PR_SET_SECCOMP, unix.SECCOMP_MODE_FILTER, 0, 0, 0); err != unix.EINVAL {
				return true
			}
		}
		return false
	}
	_, ok := s["Seccomp"]
	return ok
}

// CheckEnabled
// borrowed from
// https://github.com/opencontainers/runc/blob/v1.0.0-rc91/libcontainer/seccomp/seccomp_linux.go#L86-L102
// https://serverfault.com/questions/929589/docker-check-if-default-seccomp-profile-is-applied
// why not directly import?
// because of build tags limit, and the function in the runc is designed to check whether kernel supported,
// but here is for check whether container is protected by seccomp
func CheckEnabled() bool {
	mode, _, err := GetStatus()
	if err != nil {
		return CheckSupported()
	}
	return mode > 0
}

func GetStatus() (mode Mode, seccompFilter int, err error) {
	s, err := parseStatusFile("/proc/self/status")
	if err != nil {
		return
	}
	seccompValue, ok := s["Seccomp"]
	if !ok {
		err = fmt.Errorf("not exists filed: seccomp")
		return
	}
	m, err := strconv.Atoi(strings.TrimSpace(seccompValue))
	if err != nil {
		return
	}
	mode = Mode(m)
	seccompFilterValue, ok := s["Seccomp"]
	if !ok {
		err = fmt.Errorf("not exists filed: seccomp")
		return
	}
	seccompFilter, err = strconv.Atoi(strings.TrimSpace(seccompFilterValue))
	if err != nil {
		return
	}
	return
}
