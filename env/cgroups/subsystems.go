package cgroups

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/log"
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
		info := fmt.Sprintf("------sub systems-------\n%+q", subsystemsSupport)
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
		info += fmt.Sprintf("\n--------top level subsystem----------\n%+q", topLevelSubsystems)
		log.Logger.Info(info)
	}
	return
}
