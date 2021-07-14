package cgroups

import (
	"ctrsploit/log"
	"ctrsploit/pkg/cgroup"
	"ctrsploit/util"
	"fmt"
)

const CommandName = "cgroups"

func Version() (err error) {
	info := fmt.Sprintf("===========Cgroups=========\n")
	info += fmt.Sprintf("is cgroupv1: %v\n", util.ColorfulTickOrBallot(cgroup.IsCgroupV1()))
	info += fmt.Sprintf("is cgroupv2: %v", util.ColorfulTickOrBallot(cgroup.IsCgroupV2()))
	log.Logger.Info(info)
	return
}
