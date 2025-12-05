package service_account

import (
	"os"

	"github.com/ctrsploit/ctrsploit/prerequisite/mount/mountinfo/mountpoint"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

var ServiceAccount = mountpoint.ContainsMountPoint{
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name:   "service-account mounted",
		Info:   "",
		ExeEnv: exeenv.InContainer,
	},
	ExpectedContains: "/run/secrets/kubernetes.io/serviceaccount",
	Type:             os.ModeDir,
}
