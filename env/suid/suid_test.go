package suid

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

func TestFindSUIDFiles(t *testing.T) {
	dir := t.TempDir()
	helper := filepath.Join(dir, "helper")
	normal := filepath.Join(dir, "normal")
	link := filepath.Join(dir, "helper-link")
	skippedDir := filepath.Join(dir, "skip")
	skippedHelper := filepath.Join(skippedDir, "helper")

	if err := os.WriteFile(helper, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}
	if err := syscall.Chmod(helper, 0o4755); err != nil {
		t.Fatalf("set helper suid bit: %v", err)
	}
	if err := os.WriteFile(normal, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write normal: %v", err)
	}
	if err := os.Symlink(helper, link); err != nil {
		t.Fatalf("symlink helper: %v", err)
	}
	if err := os.Mkdir(skippedDir, 0o755); err != nil {
		t.Fatalf("mkdir skipped dir: %v", err)
	}
	if err := os.WriteFile(skippedHelper, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write skipped helper: %v", err)
	}
	if err := syscall.Chmod(skippedHelper, 0o4755); err != nil {
		t.Fatalf("set skipped helper suid bit: %v", err)
	}

	got, err := Find(Options{
		Paths:    []string{dir},
		SkipDirs: []string{skippedDir},
	})
	if err != nil {
		t.Fatalf("Find returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("files = %+v, want exactly helper", got)
	}
	if got[0].Path != helper {
		t.Fatalf("path = %q, want %q", got[0].Path, helper)
	}
	if got[0].Mode != "-rwsr-xr-x" {
		t.Fatalf("mode = %q, want -rwsr-xr-x", got[0].Mode)
	}
	if !got[0].Executable {
		t.Fatalf("helper should be executable: %+v", got[0])
	}
}

func TestParsePaths(t *testing.T) {
	got := ParsePaths(" /bin, /usr/bin ,,/sbin ")
	want := []string{"/bin", "/usr/bin", "/sbin"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParsePaths = %#v, want %#v", got, want)
	}
}

func TestFindIgnoresMissingPaths(t *testing.T) {
	got, err := Find(Options{Paths: []string{filepath.Join(t.TempDir(), "missing")}})
	if err != nil {
		t.Fatalf("Find returned error for missing path: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("files = %+v, want empty", got)
	}
}

func TestFindScansExplicitRootEvenWhenItIsInSkipDirs(t *testing.T) {
	dir := t.TempDir()
	helper := filepath.Join(dir, "helper")
	if err := os.WriteFile(helper, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}
	if err := syscall.Chmod(helper, 0o4755); err != nil {
		t.Fatalf("set helper suid bit: %v", err)
	}

	got, err := Find(Options{
		Paths:    []string{dir},
		SkipDirs: []string{dir},
	})
	if err != nil {
		t.Fatalf("Find returned error: %v", err)
	}
	if len(got) != 1 || got[0].Path != helper {
		t.Fatalf("files = %+v, want explicit root helper", got)
	}
}

func TestHumanLabelsRootOwnedExecutable(t *testing.T) {
	human := Human([]File{{
		Path:       "/usr/bin/passwd",
		Mode:       "-rwsr-xr-x",
		UID:        0,
		GID:        0,
		Size:       1,
		Executable: true,
		RootOwned:  true,
	}})
	if len(human.Files) != 1 {
		t.Fatalf("human files = %d, want 1", len(human.Files))
	}
	if got := human.Files[0].Usability.Result; !strings.Contains(got, "root-owned executable") {
		t.Fatalf("usability = %q", got)
	}
}

func TestFormatModeUsesLsStyleSpecialBits(t *testing.T) {
	tests := map[os.FileMode]string{
		os.ModeSetuid | 0o755:                 "-rwsr-xr-x",
		os.ModeSetuid | 0o644:                 "-rwSr--r--",
		os.ModeSetgid | 0o755:                 "-rwxr-sr-x",
		os.ModeSetgid | 0o644:                 "-rw-r-Sr--",
		os.ModeSticky | 0o755:                 "-rwxr-xr-t",
		os.ModeSticky | 0o754:                 "-rwxr-xr-T",
		os.ModeSetuid | os.ModeSetgid | 0o755: "-rwsr-sr-x",
	}
	for mode, want := range tests {
		if got := formatMode(mode); got != want {
			t.Fatalf("formatMode(%v) = %q, want %q", mode, got, want)
		}
	}
}
