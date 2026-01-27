package pids

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GetMax reads and returns the cgroup PID limit for the current process.
// It returns -1 if there is no limit (max).
// This function is only effective on Linux.
// https://www.kernel.org/doc/Documentation/cgroup-v1/pids.txt
// https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html
func GetMax() (int64, error) {
	cgroupPath, err := parseCgroupPath()
	if err != nil {
		return 0, fmt.Errorf("failed to parse cgroup path: %w", err)
	}

	// Construct the full path to the pids.max file.
	// The cgroup filesystem is typically mounted at /sys/fs/cgroup.
	pidMaxPath := filepath.Join("/sys/fs/cgroup", cgroupPath, "pids.max")

	// Read the content of the pids.max file.
	content, err := os.ReadFile(pidMaxPath)
	if err != nil {
		// If the file does not exist, it might mean the pids controller is not enabled for this cgroup hierarchy.
		if os.IsNotExist(err) {
			return -1, nil // Treat as no limit
		}
		return 0, fmt.Errorf("failed to read %s: %w", pidMaxPath, err)
	}

	// Parse the file content.
	limitStr := strings.TrimSpace(string(content))
	if limitStr == "max" {
		// "max" indicates that there is no limit.
		return -1, nil
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse '%s' as an integer: %w", limitStr, err)
	}

	return limit, nil
}

func parseCgroupPath() (string, error) {
	// 1. Open /proc/self/cgroup to determine the cgroup path of the current process.
	file, err := os.Open("/proc/self/cgroup")
	if err != nil {
		return "", fmt.Errorf("failed to open /proc/self/cgroup: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// cgroup v2 format: "0::/path/to/cgroup"
		// This is the unified hierarchy, so we can use its path directly.
		if strings.HasPrefix(line, "0::") {
			return strings.TrimPrefix(line, "0::"), nil
		}

		// cgroup v1 format: "id:controller_list:path"
		// e.g., "5:pids:/user.slice"
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}

		controllers := strings.Split(parts[1], ",")
		for _, c := range controllers {
			if c == "pids" {
				// assume that ctrsploit runs in the container
				return c, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("pids cgroup path not found in /proc/self/cgroup")
}
