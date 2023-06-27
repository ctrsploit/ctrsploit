package vul

import (
	cgroupv1_release_agent "github.com/ctrsploit/ctrsploit/exploit/cgroupv1-release_agent"
	"github.com/ctrsploit/ctrsploit/prerequisite"
)

type Sysadmin struct {
	BaseVulnerability
}

var (
	SysadminCgroupV1 = Sysadmin{
		BaseVulnerability: BaseVulnerability{
			Name:        "cap_sys_admin",
			Description: "Container can be escaped when has cap_sys_admin and use cgroups v1",
			CheckSecPrerequisites: prerequisite.Prerequisites{
				prerequisite.ContainsCapSysAdmin,
			},
			ExploitablePrerequisites: prerequisite.Prerequisites{
				prerequisite.MustBeRootToWriteReleaseAgent,
				prerequisite.UsingCgroupsV1,
			},
		},
	}
)

func (v Sysadmin) Exploit() {
	cgroupv1_release_agent.Exploit()
}
