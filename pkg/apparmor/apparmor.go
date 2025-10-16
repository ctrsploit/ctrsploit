package apparmor

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/lsm"
	"github.com/ssst0n3/awesome_libs/awesome_error"
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
func Mode() (string, error) {
	content, err := os.ReadFile(PathSysModuleApparmorMode)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("seems apparmor not supported: %w", err)
		}
		if errors.Is(err, os.ErrPermission) {
			return "UNKNOWN", fmt.Errorf("read %s permission denied, try run as root", PathSysModuleApparmorMode)
		}
		return "UNKNOWN", fmt.Errorf("failed to get apparmor mode: %w", err)
	}
	return strings.TrimSpace(string(content)), nil
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
