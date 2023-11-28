package lsm

import (
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"strings"
)

const (
	PathAttrCurrent        = "/proc/self/attr/current"
	KernelSupported        = 1
	KernelSupportedNot     = 0
	KernelSupportedUnknown = -1
)

/*
IsKernelSupported
Only useful under the host
*/
func IsKernelSupported(module string) (supported int) {
	// https://www.kernel.org/doc/html/v4.16/admin-guide/LSM/index.html
	content, err := os.ReadFile("/sys/kernel/security/lsm")
	if err != nil {
		if os.IsNotExist(err) {
			// means inside the ctr
			supported = KernelSupportedUnknown
		}
	} else {
		log.Logger.Debug(content)
	}
	return
}

func Current() (current string, err error) {
	content, err := os.ReadFile(PathAttrCurrent)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	current = strings.TrimSpace(string(content))
	return
}

func IsConfined() bool {
	current, err := Current()
	if err != nil {
		return false
	}
	if len(current) > 0 {
		if current != "unconfined" { // TODO: not sure for selinux
			return true
		}
	}
	return false
}
