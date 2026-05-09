package crash

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCgroupControlPathUsesUnifiedCgroup(t *testing.T) {
	dir := t.TempDir()
	procRoot := filepath.Join(dir, "proc")
	cgroupRoot := filepath.Join(dir, "sys/fs/cgroup")
	pidDir := filepath.Join(procRoot, "42")
	if err := os.MkdirAll(pidDir, 0o755); err != nil {
		t.Fatalf("mkdir pid dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pidDir, "cgroup"), []byte("0::/docker/abc\n"), 0o644); err != nil {
		t.Fatalf("write cgroup fixture: %v", err)
	}

	got, err := cgroupControlPath(procRoot, cgroupRoot, 42, "cgroup.kill")
	if err != nil {
		t.Fatalf("cgroupControlPath returned error: %v", err)
	}
	want := filepath.Join(cgroupRoot, "docker/abc/cgroup.kill")
	if got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}

func TestCgroupKillWritesOne(t *testing.T) {
	dir := t.TempDir()
	procRoot := filepath.Join(dir, "proc")
	cgroupRoot := filepath.Join(dir, "sys/fs/cgroup")
	pidDir := filepath.Join(procRoot, "7")
	cgDir := filepath.Join(cgroupRoot, "demo")
	if err := os.MkdirAll(pidDir, 0o755); err != nil {
		t.Fatalf("mkdir pid dir: %v", err)
	}
	if err := os.MkdirAll(cgDir, 0o755); err != nil {
		t.Fatalf("mkdir cgroup dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pidDir, "cgroup"), []byte("0::/demo\n"), 0o644); err != nil {
		t.Fatalf("write cgroup fixture: %v", err)
	}

	origWrite := writeFile
	t.Cleanup(func() {
		writeFile = origWrite
	})
	var gotPath string
	var gotContent string
	writeFile = func(path string, content []byte, perm os.FileMode) error {
		gotPath = path
		gotContent = string(content)
		if err := os.RemoveAll(pidDir); err != nil {
			t.Fatalf("remove pid dir: %v", err)
		}
		return nil
	}

	err := (CgroupKill{PID: 7, ProcRoot: procRoot, CgroupRoot: cgroupRoot}).Trigger(context.Background())
	if err != nil {
		t.Fatalf("CgroupKill.Trigger returned error: %v", err)
	}
	if gotPath != filepath.Join(cgDir, "cgroup.kill") {
		t.Fatalf("write path = %q", gotPath)
	}
	if gotContent != "1" {
		t.Fatalf("write content = %q, want 1", gotContent)
	}
}

func TestCgroupKillFallsBackWhenPIDStillExists(t *testing.T) {
	dir := t.TempDir()
	procRoot := filepath.Join(dir, "proc")
	cgroupRoot := filepath.Join(dir, "sys/fs/cgroup")
	pidDir := filepath.Join(procRoot, "7")
	cgDir := filepath.Join(cgroupRoot, "demo")
	if err := os.MkdirAll(pidDir, 0o755); err != nil {
		t.Fatalf("mkdir pid dir: %v", err)
	}
	if err := os.MkdirAll(cgDir, 0o755); err != nil {
		t.Fatalf("mkdir cgroup dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pidDir, "cgroup"), []byte("0::/demo\n"), 0o644); err != nil {
		t.Fatalf("write cgroup fixture: %v", err)
	}

	origWrite := writeFile
	t.Cleanup(func() {
		writeFile = origWrite
	})
	writeFile = func(path string, content []byte, perm os.FileMode) error {
		return nil
	}

	err := (CgroupKill{
		PID:         7,
		ProcRoot:    procRoot,
		CgroupRoot:  cgroupRoot,
		ExitTimeout: time.Nanosecond,
	}).Trigger(context.Background())
	if err == nil || !errors.Is(err, ErrUnsupported) {
		t.Fatalf("CgroupKill.Trigger error = %v, want ErrUnsupported", err)
	}
}
