package cgroups

import (
	"fmt"
	v1 "github.com/ctrsploit/ctrsploit/pkg/cgroup/v1"
	"github.com/ctrsploit/ctrsploit/pkg/cgroup/version"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
	"reflect"
)

const CommandCgroupsName = "cgroups"

type Result struct {
	Name result.Title `json:"name"`
	V1   item.Bool    `json:"v1"`
	V2   item.Bool    `json:"v2"`
	Sub  item.Long    `json:"sub"`
	Top  item.Long    `json:"top"`
}

func Cgroups() (err error) {
	r := Result{
		Name: result.Title{
			Name: "CGroups",
		},
		V1: item.Bool{
			Name:        "v1",
			Description: "",
			Result:      version.IsCgroupV1(),
		},
		V2: item.Bool{
			Name:        "v2",
			Description: "",
			Result:      version.IsCgroupV2(),
		},
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
	if r.V1.Result {
		r.Sub = item.Long{
			Name:        "sub systems",
			Description: "",
			Result:      fmt.Sprintf("%+q", reflect.ValueOf(subsystemsSupport).MapKeys()),
		}
		r.Top = item.Long{
			Name:        "top level subsystems",
			Description: "",
			Result:      fmt.Sprintf("%+q", topLevelSubsystems),
		}
	}
	fmt.Println(printer.Printer.Print(r))
	return
}
