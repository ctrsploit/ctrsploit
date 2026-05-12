package pipeprimitive

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestEscalateReportsShellHelperFailureAfterPasswdPatch(t *testing.T) {
	dir := t.TempDir()
	passwd := filepath.Join(dir, "passwd")
	if err := os.WriteFile(passwd, []byte("root:x:0:0:root:/root:/bin/bash\n"), 0o644); err != nil {
		t.Fatalf("write passwd fixture: %v", err)
	}

	originalPath := passwdPath
	passwdPath = passwd
	defer func() { passwdPath = originalPath }()

	primitive := &recordingPrimitive{}
	err := escalateWithShellInvoker(primitive, func() error {
		return nil
	}, func() error {
		return errors.New("su missing")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	for _, want := range []string{"passwd was patched", "su-compatible root shell failed", "su missing"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error = %q, want substring %q", err, want)
		}
	}
	if len(primitive.writes) != 1 {
		t.Fatalf("writes = %d, want 1", len(primitive.writes))
	}
	write := primitive.writes[0]
	if write.path != passwd || write.content != ":0:0:root:/root:/bin/bash\n" {
		t.Fatalf("write = %+v", write)
	}
}

func TestEscalatePreflightFailureDoesNotPatchPasswd(t *testing.T) {
	dir := t.TempDir()
	passwd := filepath.Join(dir, "passwd")
	if err := os.WriteFile(passwd, []byte("root:x:0:0:root:/root:/bin/bash\n"), 0o644); err != nil {
		t.Fatalf("write passwd fixture: %v", err)
	}

	originalPath := passwdPath
	passwdPath = passwd
	defer func() { passwdPath = originalPath }()

	primitive := &recordingPrimitive{}
	err := escalateWithShellInvoker(primitive, func() error {
		return errors.New("no setuid helper")
	}, func() error {
		t.Fatal("invoke shell should not run after preflight failure")
		return nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	for _, want := range []string{"before patching", "no setuid helper"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error = %q, want substring %q", err, want)
		}
	}
	if len(primitive.writes) != 0 {
		t.Fatalf("writes = %d, want 0", len(primitive.writes))
	}
}

func TestEscalateFallsBackToSuidOverwriteAfterPasswdSuPreflightFailure(t *testing.T) {
	dir := t.TempDir()
	passwd := filepath.Join(dir, "passwd")
	if err := os.WriteFile(passwd, []byte("root:x:0:0:root:/root:/bin/bash\n"), 0o644); err != nil {
		t.Fatalf("write passwd fixture: %v", err)
	}

	originalPath := passwdPath
	passwdPath = passwd
	defer func() { passwdPath = originalPath }()

	primitive := &recordingPrimitive{}
	err := escalateWithStrategies(primitive, []escalateStrategy{
		{name: "passwd-su", run: func(primitive Primitive) error {
			return escalateWithShellInvoker(primitive, func() error {
				return errors.New("su missing")
			}, func() error {
				t.Fatal("passwd-su invoke should not run after preflight failure")
				return nil
			})
		}},
		{name: "suid-overwrite", run: func(primitive Primitive) error {
			return primitive.Write("/usr/bin/passwd", 0, []byte("suid-payload"))
		}},
	})
	if err != nil {
		t.Fatalf("escalateWithStrategies returned error: %v", err)
	}
	if len(primitive.writes) != 1 {
		t.Fatalf("writes = %d, want 1", len(primitive.writes))
	}
	write := primitive.writes[0]
	if write.path == passwd {
		t.Fatalf("passwd was patched during fallback: %+v", write)
	}
	if write.path != "/usr/bin/passwd" || write.offset != 0 || write.content != "suid-payload" {
		t.Fatalf("write = %+v", write)
	}
}

func TestEscalateReportsAllStrategyFailures(t *testing.T) {
	primitive := &recordingPrimitive{}
	err := escalateWithStrategies(primitive, []escalateStrategy{
		{name: "passwd-su", run: func(Primitive) error { return errors.New("no su") }},
		{name: "suid-overwrite", run: func(Primitive) error { return errors.New("no suid") }},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	for _, want := range []string{"all strategies failed", "passwd-su failed: no su", "suid-overwrite failed: no suid"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error = %q, want substring %q", err, want)
		}
	}
	if len(primitive.writes) != 0 {
		t.Fatalf("writes = %d, want 0", len(primitive.writes))
	}
}

func TestSelectSuidOverwriteTargetRejectsTooSmallCandidate(t *testing.T) {
	dir := t.TempDir()
	helper := filepath.Join(dir, "helper")
	if err := os.WriteFile(helper, []byte("tiny"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}
	if err := os.Chown(helper, 0, -1); err != nil {
		t.Skipf("cannot chown helper to root: %v", err)
	}
	if err := syscall.Chmod(helper, 0o4755); err != nil {
		t.Skipf("cannot set setuid bit on helper: %v", err)
	}

	_, err := selectSuidOverwriteTarget([]string{helper}, len("larger-payload"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "too small") {
		t.Fatalf("error = %q, want too small context", err)
	}
}
