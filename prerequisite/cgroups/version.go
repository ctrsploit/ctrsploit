package cgroups

import (
	"github.com/ctrsploit/ctrsploit/pkg/cgroup/version"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Version struct {
	prerequisite.BasePrerequisite
}

var V1 = Version{
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name: "cgroups v1",
		Info: "Cgroups v1 needed",
	},
}

func (p *Version) Check() (err error) {
	err = p.BasePrerequisite.Check()
	if err != nil {
		return
	}
	p.Satisfied = version.IsCgroupV1()
	return
}
