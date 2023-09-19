package cgroups

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/cgroup/version"
	"github.com/ctrsploit/ctrsploit/util"
)

const CommandCgroupsName = "cgroups"

func Version() (err error) {
	info := fmt.Sprintf("===========Cgroups=========\n")
	info += fmt.Sprintf("is cgroupv1: %v\n", util.ColorfulTickOrBallot(version.IsCgroupV1()))
	info += fmt.Sprintf("is cgroupv2: %v", util.ColorfulTickOrBallot(version.IsCgroupV2()))
	fmt.Printf("%s\n\n", info)
	return
}
