package seccomp

import (
	"fmt"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
)

// CheckSupported
// borrowed from https://github.com/opencontainers/runc/blob/v1.0.0-rc91/libcontainer/seccomp/seccomp_linux.go#L86-L102
// https://github.com/kubernetes-sigs/security-profiles-operator/issues/76
// or on host: `grep CONFIG_SECCOMP= /boot/config-$(uname -r)` https://docs.docker.com/engine/security/seccomp/
// but it must be executed on the host, so the first two are preferred
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
	seccomp, _, err := GetStatus()
	if err != nil {
		return CheckSupported()
	}
	return seccomp > 0
}

func GetStatus() (seccomp int, seccompFilter int, err error) {
	s, err := parseStatusFile("/proc/self/status")
	if err != nil {
		return
	}
	seccompValue, ok := s["Seccomp"]
	if !ok {
		err = fmt.Errorf("not exists filed: seccomp")
		return
	}
	seccomp, err = strconv.Atoi(strings.TrimSpace(seccompValue))
	if err != nil {
		return
	}
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
