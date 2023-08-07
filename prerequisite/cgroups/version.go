package cgroups

import (
	"github.com/ctrsploit/ctrsploit/pkg/cgroup/version"
	"github.com/ctrsploit/ctrsploit/prerequisite"
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
	p.Satisfied = version.IsCgroupV1()
	return
}
