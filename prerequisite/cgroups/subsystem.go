package cgroups

import (
	v1 "github.com/ctrsploit/ctrsploit/pkg/cgroup/v1"
	"github.com/ctrsploit/ctrsploit/pkg/cgroup/version"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type TopLevelSubsystem struct {
	prerequisite.BasePrerequisite
}

var HasTopLevelSubsystem = TopLevelSubsystem{
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name:   "has top level cgroups subsystem",
		Info:   "",
		ExeEnv: exeenv.InContainer,
	},
}

func (p *TopLevelSubsystem) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	if version.IsCgroupV1() {
		topLevelSubSystems, err := v1.ListTopLevelSubSystem()
		if err != nil {
			return false, err
		}
		if len(topLevelSubSystems) > 0 {
			p.Satisfied = true
		}
	}
	p.Checked = true
	return p.Satisfied, nil
}
