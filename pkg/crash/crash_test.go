package crash

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

func TestParseMethods(t *testing.T) {
	got := ParseMethods("auto, sigkill,,oom ")
	want := []string{"auto", "sigkill", "oom"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseMethods = %#v, want %#v", got, want)
	}
}

func TestNewTriggersExpandsAuto(t *testing.T) {
	triggers, err := NewTriggers([]string{MethodAuto}, Options{})
	if err != nil {
		t.Fatalf("NewTriggers returned error: %v", err)
	}
	got := triggerNames(triggers)
	want := []string{MethodCgroupKill, MethodSigkill, MethodKillAll}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("trigger names = %#v, want %#v", got, want)
	}
}

func TestNewTriggersExpandsAll(t *testing.T) {
	triggers, err := NewTriggers([]string{MethodAll}, Options{})
	if err != nil {
		t.Fatalf("NewTriggers returned error: %v", err)
	}
	got := triggerNames(triggers)
	want := []string{MethodCgroupKill, MethodSigkill, MethodKillAll, MethodOOM}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("trigger names = %#v, want %#v", got, want)
	}
}

func TestNewTriggersDeduplicatesExpandedMethods(t *testing.T) {
	triggers, err := NewTriggers([]string{MethodAll, MethodSigkill}, Options{})
	if err != nil {
		t.Fatalf("NewTriggers returned error: %v", err)
	}
	got := triggerNames(triggers)
	want := []string{MethodCgroupKill, MethodSigkill, MethodKillAll, MethodOOM}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("trigger names = %#v, want %#v", got, want)
	}
}

func TestNewTriggersRejectsUnknownMethod(t *testing.T) {
	_, err := NewTriggers([]string{"unknown"}, Options{})
	if err == nil || !strings.Contains(err.Error(), `unknown crash trigger method "unknown"`) {
		t.Fatalf("NewTriggers error = %v, want unknown method", err)
	}
}

func TestTriggerFirstFallsBackToNextTrigger(t *testing.T) {
	firstErr := errors.New("first failed")
	first := &testTrigger{name: "first", err: firstErr}
	second := &testTrigger{name: "second"}

	if err := TriggerFirst(context.Background(), first, second); err != nil {
		t.Fatalf("TriggerFirst returned error: %v", err)
	}
	if first.calls != 1 {
		t.Fatalf("first calls = %d, want 1", first.calls)
	}
	if second.calls != 1 {
		t.Fatalf("second calls = %d, want 1", second.calls)
	}
}

func TestSigkillUsesDefaultPIDAndSIGKILL(t *testing.T) {
	procRoot := t.TempDir()
	pidDir := filepath.Join(procRoot, "1")
	if err := os.Mkdir(pidDir, 0o755); err != nil {
		t.Fatalf("mkdir pid dir: %v", err)
	}

	origKill := killProcess
	t.Cleanup(func() {
		killProcess = origKill
	})

	called := false
	killProcess = func(pid int, sig syscall.Signal) error {
		called = true
		if pid != 1 {
			t.Fatalf("pid = %d, want 1", pid)
		}
		if sig != syscall.SIGKILL {
			t.Fatalf("signal = %v, want SIGKILL", sig)
		}
		if err := os.RemoveAll(pidDir); err != nil {
			t.Fatalf("remove pid dir: %v", err)
		}
		return nil
	}

	if err := (Sigkill{ProcRoot: procRoot}).Trigger(context.Background()); err != nil {
		t.Fatalf("Sigkill.Trigger returned error: %v", err)
	}
	if !called {
		t.Fatal("killProcess was not called")
	}
}

func TestSigkillWrapsKillError(t *testing.T) {
	origKill := killProcess
	t.Cleanup(func() {
		killProcess = origKill
	})

	killErr := errors.New("boom")
	killProcess = func(pid int, sig syscall.Signal) error {
		return killErr
	}

	err := (Sigkill{PID: 42}).Trigger(context.Background())
	if !errors.Is(err, killErr) {
		t.Fatalf("error = %v, want wrapped %v", err, killErr)
	}
	if !strings.Contains(err.Error(), "kill pid 42 with SIGKILL") {
		t.Fatalf("error = %q, want kill context", err)
	}
}

type testTrigger struct {
	name  string
	err   error
	calls int
}

func (t *testTrigger) Name() string {
	return t.name
}

func (t *testTrigger) Trigger(context.Context) error {
	t.calls++
	return t.err
}

func triggerNames(triggers []Trigger) []string {
	names := make([]string, 0, len(triggers))
	for _, trigger := range triggers {
		names = append(names, trigger.Name())
	}
	return names
}
