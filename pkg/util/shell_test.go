package util

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestInvokeRootShellBySuTriesNextCandidate(t *testing.T) {
	dir := t.TempDir()
	helper := filepath.Join(dir, "helper")
	if err := os.WriteFile(helper, []byte("#!/bin/sh\necho fallback-ok\n"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err := invokeRootShellBySu([]rootShellCandidate{
		{label: "missing", path: filepath.Join(dir, "missing")},
		{label: "fallback", path: helper},
	}, strings.NewReader(""), &stdout, &stderr)
	if err != nil {
		t.Fatalf("invokeRootShellBySu returned error: %v", err)
	}
	if got := strings.TrimSpace(stdout.String()); got != "fallback-ok" {
		t.Fatalf("stdout = %q, want fallback-ok", got)
	}
}

func TestSelectRootShellBySuRejectsNonSetuidCandidateForNonRoot(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("non-root setuid preflight path only applies to non-root callers")
	}

	dir := t.TempDir()
	helper := filepath.Join(dir, "helper")
	if err := os.WriteFile(helper, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	err := checkSetuidRootExecutable(helper)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "not root") && !strings.Contains(err.Error(), "not setuid") {
		t.Fatalf("error = %q, want owner or setuid context", err)
	}
}

func TestSelectRootShellBySuAcceptsSetuidCandidateForNonRoot(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("non-root setuid preflight path only applies to non-root callers")
	}

	dir := t.TempDir()
	helper := filepath.Join(dir, "helper")
	if err := os.WriteFile(helper, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}
	if err := os.Chown(helper, 0, -1); err != nil {
		t.Skipf("cannot chown helper to root: %v", err)
	}
	if err := syscall.Chmod(helper, 0o4755); err != nil {
		t.Skipf("cannot set setuid bit on helper: %v", err)
	}

	got, err := selectRootShellBySu([]rootShellCandidate{{label: "helper", path: helper}})
	if err != nil {
		t.Fatalf("selectRootShellBySu returned error: %v", err)
	}
	if got.path != helper {
		t.Fatalf("path = %q, want %q", got.path, helper)
	}
}

func TestMountInfoForPathFindsNearestMount(t *testing.T) {
	mountPoint, _, err := mountInfoForPath("/proc/self/status")
	if err != nil {
		t.Fatalf("mountInfoForPath returned error: %v", err)
	}
	if mountPoint != "/proc" {
		t.Fatalf("mountPoint = %q, want /proc", mountPoint)
	}
}

func TestCheckSetuidAllowedReportsNoNewPrivs(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root caller bypasses no_new_privs preflight")
	}

	content, err := os.ReadFile("/proc/self/status")
	if err != nil {
		t.Skipf("cannot read /proc/self/status: %v", err)
	}
	if !strings.Contains(string(content), "NoNewPrivs:\t1") {
		t.Skip("current test process does not have NoNewPrivs set")
	}

	err = checkSetuidAllowed()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "NoNewPrivs") {
		t.Fatalf("error = %q, want NoNewPrivs context", err)
	}
}

func TestInvokeRootShellBySuReturnsContextWhenAllCandidatesFail(t *testing.T) {
	dir := t.TempDir()
	notExecutable := filepath.Join(dir, "not-executable")
	if err := os.WriteFile(notExecutable, []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatalf("write not-executable helper: %v", err)
	}

	err := invokeRootShellBySu([]rootShellCandidate{
		{label: "missing", path: filepath.Join(dir, "missing")},
		{label: "not-executable", path: notExecutable},
	}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error")
	}
	for _, want := range []string{"no su-compatible helper succeeded", "missing unavailable", "not-executable unavailable"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error = %q, want substring %q", err, want)
		}
	}
}
