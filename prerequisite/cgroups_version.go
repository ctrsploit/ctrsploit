package prerequisite

import "github.com/ctrsploit/ctrsploit/pkg/cgroup/version"

type Cgroups struct {
	BasePrerequisite
}

var UsingCgroupsV1 = Cgroups{
	BasePrerequisite: BasePrerequisite{
		Name: "cgroups v1",
		Info: "Cgroups v1 needed",
	},
}

func (p Cgroups) Check() (err error) {
	p.Satisfied = version.IsCgroupV1()
	return
}
