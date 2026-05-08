package runcwatch

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Handle struct {
	File    *os.File
	FD      int
	PID     int
	Cmdline string
	Exe     string
}

type OpenFunc func(pid int, deadline time.Time) (Handle, error)

type ProcessInfo struct {
	IsRunC  bool
	Cmdline string
	Exe     string
}

func (h Handle) Close() error {
	if h.File != nil {
		return h.File.Close()
	}
	if h.FD >= 0 {
		return syscall.Close(h.FD)
	}
	return nil
}

func CaptureHandle(timeout time.Duration, open OpenFunc) (Handle, error) {
	var deadline time.Time
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}
	maxPID, err := currentMaxPID()
	if err != nil {
		return Handle{}, err
	}
	for {
		if deadlineExceeded(deadline) {
			return Handle{}, fmt.Errorf("timeout waiting for runc process")
		}
		pids, err := os.ReadDir("/proc")
		if err != nil {
			return Handle{}, fmt.Errorf("read /proc: %w", err)
		}
		for _, f := range pids {
			pid, err := strconv.Atoi(f.Name())
			if err != nil {
				continue
			}
			if pid <= maxPID {
				continue
			}
			maxPID = pid

			info, err := RunCProcessInfo(pid)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return Handle{}, err
			}
			if !info.IsRunC {
				continue
			}

			handle, err := open(pid, deadline)
			if err != nil {
				continue
			}
			handle.PID = pid
			handle.Cmdline = info.Cmdline
			handle.Exe = info.Exe
			return handle, nil
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func OpenProcExeFile(pid int, deadline time.Time) (Handle, error) {
	path := fmt.Sprintf("/proc/%d/exe", pid)
	for {
		if deadlineExceeded(deadline) {
			return Handle{}, fmt.Errorf("timeout opening runc handle for pid %d", pid)
		}
		handle, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err == nil {
			return Handle{File: handle, FD: int(handle.Fd())}, nil
		}
		if os.IsNotExist(err) {
			return Handle{}, err
		}
		time.Sleep(time.Millisecond)
	}
}

func OpenProcExeFD(pid int, deadline time.Time) (Handle, error) {
	path := fmt.Sprintf("/proc/%d/exe", pid)
	for {
		if deadlineExceeded(deadline) {
			return Handle{}, fmt.Errorf("timeout opening runc handle for pid %d", pid)
		}
		fd, err := syscall.Open(path, syscall.O_RDONLY, 0)
		if err == nil {
			return Handle{FD: fd}, nil
		}
		if os.IsNotExist(err) {
			return Handle{}, err
		}
		time.Sleep(time.Millisecond)
	}
}

func IsRunCPID(pid int) (bool, error) {
	info, err := RunCProcessInfo(pid)
	if err != nil {
		return false, err
	}
	return info.IsRunC, nil
}

func RunCProcessInfo(pid int) (ProcessInfo, error) {
	content, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return ProcessInfo{}, err
	}
	cmdline := string(content)
	exe, _ := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	return ProcessInfo{
		IsRunC:  IsRunCProcess(cmdline, exe),
		Cmdline: cmdline,
		Exe:     exe,
	}, nil
}

func IsRunCProcess(cmdline, exe string) bool {
	argv0 := firstArg(cmdline)
	if isRunCName(filepath.Base(argv0)) {
		return true
	}

	// runc may re-exec itself as `/proc/self/exe init`, which is still a runc
	// process even if the exe symlink disappears before we read it.
	if cmdline == "/proc/self/exe\x00init\x00" {
		return true
	}

	// A bare `/proc/self/exe` is only considered runc when the kernel exe
	// symlink still points at runc. This avoids matching helper processes that
	// intentionally run `/proc/self/exe`, such as the image exec-mode capture.
	if argv0 != "/proc/self/exe" {
		return false
	}

	return isRunCName(filepath.Base(exe))
}

func IsRunCCmdline(cmdline string) bool {
	return IsRunCProcess(cmdline, "")
}

func isRunCName(name string) bool {
	switch name {
	case "runc", "docker-runc":
		return true
	default:
		return false
	}
}

func firstArg(cmdline string) string {
	if n := strings.IndexByte(cmdline, 0); n >= 0 {
		return cmdline[:n]
	}
	return cmdline
}

func currentMaxPID() (int, error) {
	pids, err := os.ReadDir("/proc")
	if err != nil {
		return 0, fmt.Errorf("read /proc: %w", err)
	}
	maxPID := 0
	for _, f := range pids {
		pid, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}
		if pid > maxPID {
			maxPID = pid
		}
	}
	return maxPID, nil
}

func deadlineExceeded(deadline time.Time) bool {
	return !deadline.IsZero() && time.Now().After(deadline)
}
