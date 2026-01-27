package cgroups

import (
	"fmt"

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
	return p.CheckTemplate(func() {
		if version.IsCgroupV1() {
			topLevelSubSystems, err := v1.ListTopLevelSubSystem()
			if err != nil {
				p.Err = p.WrapErr(fmt.Errorf("listing top level subsystems: %w", err))
				return
			}
			if len(topLevelSubSystems) > 0 {
				p.Satisfied = true
			}
		}
		return
	})
}
