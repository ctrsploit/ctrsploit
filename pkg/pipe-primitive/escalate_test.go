package pipeprimitive

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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
