package sysctl

import (
	"os"
	"strings"
)

const (
	PathUnprivilegedUsernsClone = "/proc/sys/kernel/unprivileged_userns_clone"
)

func UnprivilegedUsernsCloneEnabled() (bool, error) {
	content, err := os.ReadFile(PathUnprivilegedUsernsClone)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(content)) == "1", nil
}
