package selinux

import (
	"github.com/opencontainers/selinux/go-selinux"
)

type TypeMode int

func (m TypeMode) String() (mode string) {
	switch m {
	case -1:
		mode = "disabled"
	case 0:
		mode = "permissive"
	case 1:
		mode = "enforcing"
	default:
		mode = "unknown"
	}
	return
}

func IsEnabled() bool {
	return selinux.GetEnabled()
}

func Mode() TypeMode {
	return TypeMode(selinux.EnforceMode())
}
