package cgroups

import (
	"fmt"
	v1 "github.com/ctrsploit/ctrsploit/pkg/cgroup/v1"
	"github.com/ctrsploit/ctrsploit/pkg/cgroup/version"
)

func ListSubsystems() (err error) {
	if version.IsCgroupV1() {
		c := v1.CgroupV1{}
		subsystemsSupport, err := c.ListSubsystems(version.UnifiedMountpoint)
		if err != nil {
			return err
		}
		info := fmt.Sprintf("\n------sub systems-------\n%+q", subsystemsSupport)
		var topLevelSubsystems []string
		for _, subsystem := range subsystemsSupport {
			top, err := c.IsTop(version.UnifiedMountpoint, subsystem)
			if err != nil {
				return err
			}
			if top {
				topLevelSubsystems = append(topLevelSubsystems, subsystem)
			}
		}
		info += fmt.Sprintf("\n\n--------top level subsystem----------\n%+q", topLevelSubsystems)
		fmt.Println(info)
	}
	return
}
