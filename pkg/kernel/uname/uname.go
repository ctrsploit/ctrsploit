package uname

import (
	"fmt"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"golang.org/x/sys/unix"
)

func Uname() (u unix.Utsname, err error) {
	err = unix.Uname(&u)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	//log.Logger.Debugf("%+v\n", u)
	return
}

func All() (all string, err error) {
	u, err := Uname()
	if err != nil {
		return
	}
	all = fmt.Sprintf("%s %s %s %s %s", sysname(u), nodeName(u), release(u), version(u), machine(u))
	return
}

func sysname(u unix.Utsname) string {
	return byteSliceToString(u.Sysname)
}

func Sysname() (name string, err error) {
	u, err := Uname()
	if err != nil {
		return
	}
	name = sysname(u)
	return
}

func DomainName() (name string, err error) {
	u, err := Uname()
	if err != nil {
		return
	}
	name = domainName(u)
	return
}

func nodeName(u unix.Utsname) string {
	return byteSliceToString(u.Sysname)
}

func NodeName() (name string, err error) {
	u, err := Uname()
	if err != nil {
		return
	}
	name = nodeName(u)
	return
}

func release(u unix.Utsname) string {
	return byteSliceToString(u.Release)
}

func Release() (r string, err error) {
	u, err := Uname()
	if err != nil {
		return
	}
	r = release(u)
	return
}

func version(u unix.Utsname) string {
	return byteSliceToString(u.Version)
}

func Version() (v string, err error) {
	u, err := Uname()
	if err != nil {
		return
	}
	v = version(u)
	return
}

func machine(u unix.Utsname) string {
	return byteSliceToString(u.Machine)
}

func Machine() (m string, err error) {
	u, err := Uname()
	if err != nil {
		return
	}
	m = machine(u)
	return
}
