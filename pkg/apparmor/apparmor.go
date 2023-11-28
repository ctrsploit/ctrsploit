package apparmor

import (
	"github.com/ctrsploit/ctrsploit/pkg/lsm"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"strings"
)

const (
	dirSysModuleApparmor         = "/sys/module/apparmor/parameters"
	PathSysModuleApparmorEnabled = dirSysModuleApparmor + "/enabled"
	PathSysModuleApparmorMode    = dirSysModuleApparmor + "/mode"
)

/*
Mode
Make sure the apparmor is supported by yourself
*/
func Mode() (mode string, err error) {
	content, err := os.ReadFile(PathSysModuleApparmorMode)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	mode = strings.TrimSpace(string(content))
	return
}

func IsSupport() (support bool) {
	content, err := os.ReadFile(PathSysModuleApparmorEnabled)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") { // not found means not support
			return
		}
		awesome_error.CheckErr(err)
		return
	}
	if strings.TrimSpace(string(content)) == "Y" {
		support = true
	}
	return
}

func IsEnabled() bool {
	return IsSupport() && lsm.IsConfined()
}
