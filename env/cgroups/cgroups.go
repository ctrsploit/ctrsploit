package cgroups

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/internal/colorful"
	v1 "github.com/ctrsploit/ctrsploit/pkg/cgroup/v1"
	"github.com/ctrsploit/ctrsploit/pkg/cgroup/version"
	"github.com/ssst0n3/awesome_libs"
	"reflect"
)

const CommandCgroupsName = "cgroups"

var (
	tplCgroups = `
===========Cgroups===========
{.version}
{.sub}
`
	tplVersion = `{.v1}  cgroups v1
{.v2}  cgroups v2`
	tplSubsystem = `---
sub systems:
{.subsystem}
top level subsystems:
{.top}`
)

func versionInfo() (info string) {
	info = awesome_libs.Format(tplVersion, awesome_libs.Dict{
		"v1": colorful.TickOrBallot(version.IsCgroupV1()),
		"v2": colorful.TickOrBallot(version.IsCgroupV2()),
	})
	return
}

func listSubsystems() (info string, err error) {
	if !version.IsCgroupV1() {
		return
	}
	c := v1.CgroupV1{}
	subsystemsSupport, err := c.ListSubsystems("/proc/1/cgroup")
	if err != nil {
		return
	}
	var topLevelSubsystems []string
	for subsystemName, subsystemPath := range subsystemsSupport {
		if c.IsTop(subsystemPath) {
			topLevelSubsystems = append(topLevelSubsystems, subsystemName)
		}
	}
	info = awesome_libs.Format(tplSubsystem, awesome_libs.Dict{
		"subsystem": fmt.Sprintf("%+q", reflect.ValueOf(subsystemsSupport).MapKeys()),
		"top":       fmt.Sprintf("%+q", topLevelSubsystems),
	})
	return
}

func Cgroups() (err error) {
	sub, err := listSubsystems()
	if err != nil {
		return
	}
	info := awesome_libs.Format(tplCgroups, awesome_libs.Dict{
		"version": versionInfo(),
		"sub":     sub,
	})
	print(info)
	return
}
