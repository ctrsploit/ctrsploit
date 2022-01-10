package cgroups

import (
	"fmt"
	v1 "github.com/ctrsploit/ctrsploit/pkg/cgroup/v1"
	"github.com/ctrsploit/ctrsploit/pkg/cgroup/version"
	"reflect"
)

func ListSubsystems() (err error) {
	if version.IsCgroupV1() {
		c := v1.CgroupV1{}
		subsystemsSupport, err := c.ListSubsystems("/proc/1/cgroup")
		if err != nil {
			return err
		}
		info := fmt.Sprintf("\n------sub systems-------\n%+q", reflect.ValueOf(subsystemsSupport).MapKeys())
		var topLevelSubsystems []string
		for subsystemName, subsystemPath := range subsystemsSupport {
			if c.IsTop(subsystemPath) {
				topLevelSubsystems = append(topLevelSubsystems, subsystemName)
			}
		}
		info += fmt.Sprintf("\n\n--------top level subsystem----------\n%+q", topLevelSubsystems)
		fmt.Println(info)
	}
	return
}
