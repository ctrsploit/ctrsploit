package selinux

import (
	"github.com/opencontainers/selinux/go-selinux"
)

func IsEnabled() bool {
	return selinux.GetEnabled()
}

func Mode() int {
	return selinux.EnforceMode()
}

func Translate(mode int) (str string) {
	switch mode {
	case -1:
		str = "disabled"
	case 0:
		str = "permissive"
	case 1:
		str = "enforcing"
	}
	return
}
