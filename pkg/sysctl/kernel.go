package sysctl

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/log"
)

const (
	PathUnprivilegedUsernsClone = "/proc/sys/kernel/unprivileged_userns_clone"
	PathPidMax                  = "/proc/sys/kernel/pid_max"
)

func UnprivilegedUsernsCloneEnabled() (bool, error) {
	content, err := os.ReadFile(PathUnprivilegedUsernsClone)
	if err != nil {
		if os.IsNotExist(err) {
			log.Logger.Warnf("%s does not exist, assuming unprivileged user namespaces are enabled", PathUnprivilegedUsernsClone)
			// Assume unprivileged user namespaces are enabled if the sysctl does not exist
			return true, nil
		}
		return false, err
	}
	return strings.TrimSpace(string(content)) == "1", nil
}

func PidMax() (int, error) {
	content, err := os.ReadFile(PathPidMax)
	if err != nil {
		return 0, fmt.Errorf("failed to read max pid: %w", err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return 0, fmt.Errorf("failed to parse max pid: %w", err)
	}
	return pid, nil
}
