package uname

func byteSliceToString(s [65]byte) (str string) {
	for _, r := range s {
		if r == 0 {
			break
		}
		str += string(r)
	}
	return
}
