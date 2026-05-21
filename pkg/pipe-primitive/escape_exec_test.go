package pipeprimitive

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestProcessEntrypointPathPrefersShebangScript(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "ctrsploit-shebang-entrypoint")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nsleep 30\n"), 0o755); err != nil {
		t.Fatalf("write shebang script: %v", err)
	}

	cmd := exec.Command(script)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start shebang script: %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	})

	time.Sleep(200 * time.Millisecond)
	got, err := processEntrypointPath(cmd.Process.Pid)
	if err != nil {
		t.Fatalf("processEntrypointPath returned error: %v", err)
	}
	want, err := filepath.EvalSymlinks(script)
	if err != nil {
		t.Fatalf("eval script symlink: %v", err)
	}
	if got != want {
		t.Fatalf("entrypoint path = %q, want shebang script %q", got, want)
	}
}
