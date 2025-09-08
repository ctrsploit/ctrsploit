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

func (p *Version) Check() (bool, error) {
	if !p.Checked {
		p.Satisfied = version.IsCgroupV1()
		p.Checked = true
	}
	return p.Satisfied, nil
}
