package kernel

import "golang.org/x/sys/unix"

func Uname() {
	u := unix.Utsname{}
	unix.Uname(&u)
}
