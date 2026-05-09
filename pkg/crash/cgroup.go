package crash

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultCgroupRoot = "/sys/fs/cgroup"

var writeFile = os.WriteFile

const (
	defaultExitCheckTimeout  = 2 * time.Second
	defaultExitCheckInterval = 50 * time.Millisecond
)

type CgroupKill struct {
	PID         int
	Delay       time.Duration
	ProcRoot    string
	CgroupRoot  string
	ExitTimeout time.Duration
}

func (c CgroupKill) Name() string {
	return MethodCgroupKill
}

func (c CgroupKill) Trigger(ctx context.Context) error {
	pid := normalizedPID(c.PID)
	if err := wait(ctx, c.Delay); err != nil {
		return err
	}
	path, err := cgroupControlPath(c.ProcRoot, c.CgroupRoot, pid, "cgroup.kill")
	if err != nil {
		return err
	}
	if err := writeFile(path, []byte("1"), 0); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return waitForProcessExit(ctx, c.ProcRoot, pid, c.ExitTimeout)
}

func cgroupControlPath(procRoot, cgroupRoot string, pid int, control string) (string, error) {
	if procRoot == "" {
		procRoot = "/proc"
	}
	if cgroupRoot == "" {
		cgroupRoot = defaultCgroupRoot
	}
	cgroupPath, err := unifiedCgroupPath(filepath.Join(procRoot, fmt.Sprintf("%d/cgroup", pid)))
	if err != nil {
		return "", err
	}
	if cgroupPath == "" {
		return "", fmt.Errorf("%w: unified cgroup path not found for pid %d", ErrUnsupported, pid)
	}
	return filepath.Join(cgroupRoot, strings.TrimPrefix(cgroupPath, "/"), control), nil
}

func unifiedCgroupPath(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !bytes.HasPrefix(line, []byte("0::")) {
			continue
		}
		return string(bytes.TrimPrefix(line, []byte("0::"))), nil
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan %s: %w", path, err)
	}
	return "", nil
}

func waitForProcessExit(ctx context.Context, procRoot string, pid int, timeout time.Duration) error {
	if procRoot == "" {
		procRoot = "/proc"
	}
	if timeout <= 0 {
		timeout = defaultExitCheckTimeout
	}

	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	ticker := time.NewTicker(defaultExitCheckInterval)
	defer ticker.Stop()

	processPath := filepath.Join(procRoot, fmt.Sprint(pid))
	for {
		if _, err := os.Stat(processPath); os.IsNotExist(err) {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline.C:
			return fmt.Errorf("%w: pid %d still exists after cgroup.kill", ErrUnsupported, pid)
		case <-ticker.C:
		}
	}
}
