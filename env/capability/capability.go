package capability

import (
	"fmt"
	"github.com/containerd/containerd/pkg/cap"
	"github.com/ctrsploit/ctrsploit/pkg/capability"
	"github.com/ctrsploit/sploit-spec/pkg/colorful"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

const (
	CommandCapabilityName        = "capability"
	standardCaps          uint64 = 0xa80425fb
)

type Cap struct {
	SubTitle     result.SubTitle `json:"-"`
	Capabilities item.Short      `json:"capabilities"`
	NotDefault   item.Bool       `json:"not_default"`
	Additional   item.Long       `json:"additional"`
}

type Caps struct {
	Name    result.Title `json:"name"`
	Pid1    Cap          `json:"pid1"`
	Current Cap          `json:"current"`
}

func getInfoFromCaps(caps uint64, subtitle string) (c Cap) {
	c.SubTitle = result.SubTitle{
		Name: subtitle,
	}
	c.Capabilities = item.Short{
		Name:   "capabilities",
		Result: fmt.Sprintf("0x%x", caps),
	}
	c.NotDefault = item.Bool{
		Name:        fmt.Sprintf("Not Equal to Docker's Default Capability (0x%x)", standardCaps),
		Description: fmt.Sprintf("0x%x", caps),
		Result:      caps != standardCaps,
	}
	if caps != standardCaps {
		capsDiff, _ := cap.FromBitmap(caps & (^standardCaps))
		c.Additional = item.Long{
			Name:        "[Additional]",
			Description: "",
			Result:      colorful.O.Danger(fmt.Sprintf("%q", capsDiff)),
		}
	}
	return
}

func Capability() (err error) {
	caps, err := capability.GetPid1Capability()
	if err != nil {
		return err
	}
	pid1 := getInfoFromCaps(caps, "pid1")

	caps, err = capability.GetCurrentCapability()
	if err != nil {
		return err
	}
	current := getInfoFromCaps(caps, "current")

	c := Caps{
		Name: result.Title{
			Name: "Capability",
		},
		Pid1:    pid1,
		Current: current,
	}
	fmt.Println(printer.Printer.PrintDropAfterFalse(c))
	return
}
