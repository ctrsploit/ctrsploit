package runc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

func StraceFGetSeals() (bool, error) {
	parent, child, err := newPipe()
	if err != nil {
		return false, fmt.Errorf("failed to create pipe: %w", err)
	}
	cmd, err := buildCmd(child)
	if err != nil {
		return false, fmt.Errorf("failed to build command: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return false, fmt.Errorf("nsenter failed to start: %w", err)
	}
	r := nl.NewNetlinkRequest(int(libcontainer.InitMsg), 0)
	if _, err := io.Copy(parent, bytes.NewReader(r.Serialize())); err != nil {
		return false, fmt.Errorf("failed to write init message: %w", err)
	}
	has, err := straceFGetSeals(cmd.Process.Pid)
	if err != nil {
		return false, fmt.Errorf("strace failed: %w", err)
	}
	return has, nil
}

func newPipe() (parent *os.File, child *os.File, err error) {
	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, nil, err
	}
	return os.NewFile(uintptr(fds[0]), "parent"), os.NewFile(uintptr(fds[1]), "child"), nil
}

func buildCmd(child *os.File) (*exec.Cmd, error) {
	path, err := LookRunC()
	if err != nil {
		return nil, err
	}
	cmd := &exec.Cmd{
		Path: path,
		// Args:       []string{"nsenter-exec"},
		Args:       []string{"runc"},
		ExtraFiles: []*os.File{child},
		Env:        []string{"_LIBCONTAINER_INITPIPE=3"},
		Stdout:     io.Discard,
		Stderr:     io.Discard,
		SysProcAttr: &syscall.SysProcAttr{
			Ptrace: true,
		},
	}
	return cmd, nil
}

func straceFGetSeals(pid int) (bool, error) {
	var ws syscall.WaitStatus
	if _, err := syscall.Wait4(pid, &ws, 0, nil); err != nil {
		return false, fmt.Errorf("initial wait failed: %w", err)
	}
	inSyscall := true
	for {
		if err := syscall.PtraceSyscall(pid, 0); err != nil {
			return false, fmt.Errorf("PtraceSyscall failed: %w", err)
		}
		if _, err := syscall.Wait4(pid, &ws, 0, nil); err != nil {
			return false, fmt.Errorf("wait failed: %w", err)
		}
		if ws.Exited() {
			return false, nil
		}
		if inSyscall {
			var regs syscall.PtraceRegs
			if err := syscall.PtraceGetRegs(pid, &regs); err != nil {
				return false, fmt.Errorf("PtraceGetRegs failed: %w", err)
			}
			if regs.Orig_rax == unix.SYS_FCNTL {
				if regs.Rsi == unix.F_GET_SEALS {
					return true, nil
				}
			}
		}
		inSyscall = !inSyscall
	}
}
