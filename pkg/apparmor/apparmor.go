package apparmor

import (
	"ctrsploit/pkg/lsm"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"strings"
)

const (
	dirSysModuleApparmor         = "/sys/module/apparmor/parameters"
	PathSysModuleApparmorEnabled = dirSysModuleApparmor + "/enabled"
	PathSysModuleApparmorMode = dirSysModuleApparmor + "/mode"
)

/*
Mode
Make sure the apparmor is supported by yourself
*/
func Mode() (mode string, err error) {
	content, err := ioutil.ReadFile(PathSysModuleApparmorMode)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	mode = strings.TrimSpace(string(content))
	return
}

func IsSupport() (support bool) {
	content, err := ioutil.ReadFile(PathSysModuleApparmorEnabled)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if strings.TrimSpace(string(content)) == "Y" {
		support = true
	}
	return
}

func IsEnabled() bool {
	return IsSupport() && lsm.IsEnabled()
}
