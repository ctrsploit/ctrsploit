// Package crash provides composable ways to make the current container exit so
// an external restart policy can start it again.
package crash

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	MethodAuto       = "auto"
	MethodAll        = "all"
	MethodCgroupKill = "cgroup-kill"
	MethodSigkill    = "sigkill"
	MethodKillAll    = "kill-all"
	MethodOOM        = "oom"
)

var ErrUnsupported = errors.New("crash trigger unsupported")

type Options struct {
	PID        int
	Delay      time.Duration
	ProcRoot   string
	CgroupRoot string

	OOMMaxBytes   uint64
	OOMChunkBytes int
}

type Trigger interface {
	Name() string
	Trigger(ctx context.Context) error
}

func DefaultMethods() []string {
	return []string{MethodCgroupKill, MethodSigkill, MethodKillAll}
}

func AllMethods() []string {
	return []string{MethodCgroupKill, MethodSigkill, MethodKillAll, MethodOOM}
}

func ParseMethods(value string) []string {
	var methods []string
	for _, method := range strings.Split(value, ",") {
		method = strings.TrimSpace(method)
		if method == "" {
			continue
		}
		methods = append(methods, method)
	}
	return methods
}

func NewTriggers(methods []string, opts Options) ([]Trigger, error) {
	methods, err := normalizeMethods(methods)
	if err != nil {
		return nil, err
	}

	triggers := make([]Trigger, 0, len(methods))
	for _, method := range methods {
		switch method {
		case MethodCgroupKill:
			triggers = append(triggers, CgroupKill{
				PID:        opts.PID,
				Delay:      opts.Delay,
				ProcRoot:   opts.ProcRoot,
				CgroupRoot: opts.CgroupRoot,
				// A successful write to cgroup.kill is not enough for fallback
				// logic; confirm the target actually disappeared.
				ExitTimeout: defaultExitCheckTimeout,
			})
		case MethodSigkill:
			triggers = append(triggers, Sigkill{
				PID:         opts.PID,
				Delay:       opts.Delay,
				ProcRoot:    opts.ProcRoot,
				ExitTimeout: defaultExitCheckTimeout,
			})
		case MethodKillAll:
			triggers = append(triggers, KillAll{
				PID:      opts.PID,
				Delay:    opts.Delay,
				ProcRoot: opts.ProcRoot,
			})
		case MethodOOM:
			triggers = append(triggers, OOM{
				PID:           opts.PID,
				Delay:         opts.Delay,
				ProcRoot:      opts.ProcRoot,
				CgroupRoot:    opts.CgroupRoot,
				MaxBytes:      opts.OOMMaxBytes,
				ChunkBytes:    opts.OOMChunkBytes,
				TargetOOMBias: 1000,
				SelfOOMBias:   -1000,
			})
		default:
			return nil, fmt.Errorf("unknown crash trigger method %q", method)
		}
	}
	return triggers, nil
}

func TriggerFirst(ctx context.Context, triggers ...Trigger) error {
	var errs []error
	for _, trigger := range triggers {
		if trigger == nil {
			continue
		}
		if err := ctx.Err(); err != nil {
			errs = append(errs, err)
			break
		}
		if err := trigger.Trigger(ctx); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", trigger.Name(), err))
			continue
		}
		return nil
	}
	if len(errs) == 0 {
		return fmt.Errorf("%w: no crash triggers configured", ErrUnsupported)
	}
	return errors.Join(errs...)
}

func normalizeMethods(methods []string) ([]string, error) {
	if len(methods) == 0 {
		return DefaultMethods(), nil
	}

	var normalized []string
	seen := map[string]bool{}
	add := func(method string) {
		if seen[method] {
			return
		}
		seen[method] = true
		normalized = append(normalized, method)
	}

	for _, method := range methods {
		method = strings.TrimSpace(strings.ToLower(method))
		if method == "" {
			continue
		}
		switch method {
		case MethodAuto:
			for _, defaultMethod := range DefaultMethods() {
				add(defaultMethod)
			}
		case MethodAll:
			for _, allMethod := range AllMethods() {
				add(allMethod)
			}
		case "kill", "signal":
			add(MethodSigkill)
		case "killall", "kill-children":
			add(MethodKillAll)
		case "cgroup":
			add(MethodCgroupKill)
		default:
			add(method)
		}
	}
	if len(normalized) == 0 {
		return DefaultMethods(), nil
	}
	return normalized, nil
}

func normalizedPID(pid int) int {
	if pid <= 0 {
		return 1
	}
	return pid
}

func wait(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
