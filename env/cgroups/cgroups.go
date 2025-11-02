package cgroups

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/cgroup/pids"
	v1 "github.com/ctrsploit/ctrsploit/pkg/cgroup/v1"
	"github.com/ctrsploit/ctrsploit/pkg/cgroup/version"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const CommandCgroupsName = "cgroups"

func Cgroups() (cgroups container.CGroups, err error) {
	cgroups = container.CGroups{
		Subsystems:         []string{},
		TopLevelSubSystems: []string{},
	}
	if version.IsCgroupV1() {
		cgroups.Version = container.CgroupsV1
	}
	if version.IsCgroupV2() {
		cgroups.Version = container.CgroupsV2
	}

	var c v1.CgroupV1
	subsystemsSupport, err := c.ListSubsystems(v1.DefaultMountPoint)
	if err != nil {
		return
	}
	for _, subsystemName := range subsystemsSupport {
		cgroups.Subsystems = append(cgroups.Subsystems, subsystemName)
		is, err := c.IsTop(v1.DefaultMountPoint, subsystemName)
		if err != nil {
			return cgroups, fmt.Errorf("ListTopLevelSubSystem: failed to list sub system %s: %w", subsystemName, err)
		}
		if is {
			cgroups.TopLevelSubSystems = append(cgroups.TopLevelSubSystems, subsystemName)
		}
	}
	cgroups.PidsMax, err = pids.GetMax()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	return
}
