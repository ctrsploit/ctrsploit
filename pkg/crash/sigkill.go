package crash

import (
	"context"
	"fmt"
	"syscall"
	"time"
)

var killProcess = syscall.Kill

type Sigkill struct {
	PID         int
	Delay       time.Duration
	ProcRoot    string
	ExitTimeout time.Duration
}

func (s Sigkill) Name() string {
	return MethodSigkill
}

func (s Sigkill) Trigger(ctx context.Context) error {
	pid := normalizedPID(s.PID)
	if err := wait(ctx, s.Delay); err != nil {
		return err
	}
	if err := killProcess(pid, syscall.SIGKILL); err != nil {
		return fmt.Errorf("kill pid %d with SIGKILL: %w", pid, err)
	}
	return waitForProcessExit(ctx, s.ProcRoot, pid, s.ExitTimeout)
}
