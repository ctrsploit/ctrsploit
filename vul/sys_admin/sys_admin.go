package sys_admin

import (
	cgroupv1_release_agent "github.com/ctrsploit/ctrsploit/exploit/cgroupv1-release_agent"
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite/user"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

type Sysadmin struct {
	vul.BaseVulnerability
}

var (
	SysadminCgroupV1 = Sysadmin{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "cap_sys_admin",
			Description: "Container can be escaped when has cap_sys_admin and use cgroups v1",
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites: prerequisite.Prerequisites{
				&capability.CapSysAdminBnd,
			},
			ExploitablePrerequisites: prerequisite.Prerequisites{
				&capability.CapSysAdminEff,
				&user.MustBeRootToWriteReleaseAgent,
				&cgroups.V1,
			},
		},
	}
)

func (v Sysadmin) Exploit(context *cli.Context) (err error) {
	err = v.BaseVulnerability.Exploit(context)
	if err != nil {
		return
	}
	cgroupv1_release_agent.Exploit()
	return
}
