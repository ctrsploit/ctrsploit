package crash

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOOMLimitUsesMemoryMax(t *testing.T) {
	dir := t.TempDir()
	procRoot := filepath.Join(dir, "proc")
	cgroupRoot := filepath.Join(dir, "sys/fs/cgroup")
	pidDir := filepath.Join(procRoot, "9")
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
	if err := os.WriteFile(filepath.Join(cgDir, "memory.max"), []byte("1048576\n"), 0o644); err != nil {
		t.Fatalf("write memory.max fixture: %v", err)
	}

	got, err := (OOM{PID: 9, ProcRoot: procRoot, CgroupRoot: cgroupRoot}).limit()
	if err != nil {
		t.Fatalf("OOM.limit returned error: %v", err)
	}
	want := uint64(1048576 + defaultOOMChunkBytes)
	if got != want {
		t.Fatalf("limit = %d, want %d", got, want)
	}
}

func TestOOMLimitUsesFallbackForUnlimitedMemory(t *testing.T) {
	dir := t.TempDir()
	procRoot := filepath.Join(dir, "proc")
	cgroupRoot := filepath.Join(dir, "sys/fs/cgroup")
	pidDir := filepath.Join(procRoot, "9")
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
	if err := os.WriteFile(filepath.Join(cgDir, "memory.max"), []byte("max\n"), 0o644); err != nil {
		t.Fatalf("write memory.max fixture: %v", err)
	}

	got, err := (OOM{PID: 9, ProcRoot: procRoot, CgroupRoot: cgroupRoot}).limit()
	if err != nil {
		t.Fatalf("OOM.limit returned error: %v", err)
	}
	if got != maxOOMBytesFallback {
		t.Fatalf("limit = %d, want %d", got, maxOOMBytesFallback)
	}
}
