package crash

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type KillAll struct {
	PID      int
	Delay    time.Duration
	ProcRoot string
}

func (k KillAll) Name() string {
	return MethodKillAll
}

func (k KillAll) Trigger(ctx context.Context) error {
	if err := wait(ctx, k.Delay); err != nil {
		return err
	}
	procRoot := k.ProcRoot
	if procRoot == "" {
		procRoot = "/proc"
	}
	targetPID := normalizedPID(k.PID)
	selfPID := os.Getpid()

	entries, err := os.ReadDir(procRoot)
	if err != nil {
		return fmt.Errorf("read %s: %w", procRoot, err)
	}

	var errs []error
	killed := 0
	for _, entry := range entries {
		pid, ok := procPID(entry.Name())
		if !ok || pid == selfPID || pid == targetPID {
			continue
		}
		if err := killProcess(pid, syscall.SIGKILL); err != nil {
			errs = append(errs, fmt.Errorf("kill pid %d with SIGKILL: %w", pid, err))
			continue
		}
		killed++
	}
	if killed == 0 {
		errs = append(errs, fmt.Errorf("%w: no process killed", ErrUnsupported))
	}
	if err := waitForProcessExit(ctx, procRoot, targetPID, defaultExitCheckTimeout); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("%w: %w", ErrUnsupported, errors.Join(errs...))
	}
	return nil
}

func procPID(name string) (int, bool) {
	if name == "" || strings.Trim(name, "0123456789") != "" {
		return 0, false
	}
	pid, err := strconv.Atoi(name)
	return pid, err == nil && pid > 0
}
