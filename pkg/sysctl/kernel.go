package sysctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/sploit-spec/pkg/log"
)

const (
	PathUnprivilegedUsernsClone = "/proc/sys/kernel/unprivileged_userns_clone"
	PathPidMax                  = "/proc/sys/kernel/pid_max"
	PathThreadsMax              = "/proc/sys/kernel/threads-max"
)

func UnprivilegedUsernsCloneEnabled() (bool, error) {
	content, err := os.ReadFile(PathUnprivilegedUsernsClone)
	if err != nil {
		if os.IsNotExist(err) {
			log.Logger.Debugf("%s does not exist, assuming unprivileged user namespaces are enabled", PathUnprivilegedUsernsClone)
			// Assume unprivileged user namespaces are enabled if the sysctl does not exist
			return true, nil
		}
		return false, err
	}
	return strings.TrimSpace(string(content)) == "1", nil
}

func PidMax() (uint64, error) {
	pid, err := internal.ReadUint64(PathPidMax)
	if err != nil {
		return 0, fmt.Errorf("failed to read pid max: %w", err)
	}
	return pid, nil
}

func ThreadsMax() (uint64, error) {
	threads, err := internal.ReadUint64(PathThreadsMax)
	if err != nil {
		return 0, fmt.Errorf("failed to read threads max: %w", err)
	}
	return threads, nil
}
