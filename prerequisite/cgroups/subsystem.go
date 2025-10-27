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
	return p.CheckTemplate(func() (bool, error) {
		if version.IsCgroupV1() {
			topLevelSubSystems, err := v1.ListTopLevelSubSystem()
			if err != nil {
				p.Err = fmt.Errorf("failed to check [%s] caused by listing top level subsystems: %w", p.GetName(), err)
				return p.Satisfied, p.Err
			}
			if len(topLevelSubSystems) > 0 {
				p.Satisfied = true
			}
		}
		return p.Satisfied, p.Err
	})
}
