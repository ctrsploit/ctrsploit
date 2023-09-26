package capability

import (
	"fmt"
	"github.com/containerd/containerd/pkg/cap"
	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/ctrsploit/pkg/capability"
	"github.com/ctrsploit/sploit-spec/pkg/colorful"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

const (
	CommandCapabilityName        = "capability"
	standardCaps          uint64 = 0xa80425fb
)

type Cap struct {
	Capabilities item.Long `json:"capabilities"`
	Equal        item.Bool `json:"equal"`
	Additional   item.Long `json:"additional"`
}

type Caps struct {
	Name    result.Title
	Pid1    Cap `json:"pid1"`
	Current Cap `json:"current"`
}

func (c Cap) String() (s string) {
	s += internal.Print(c.Capabilities, c.Equal)
	if !c.Equal.Result {
		s += internal.Print(c.Additional)
	}
	return
}

func (c Caps) String() (s string) {
	s += internal.Print(c.Name)
	s += c.Pid1.String() + "\n"
	s += c.Current.String() + "\n"
	return
}

func getInfoFromCaps(caps uint64, subtitle string) (c Cap) {
	c.Capabilities = item.Long{
		Name:   fmt.Sprintf("[Capabilities (%s)]", subtitle),
		Result: fmt.Sprintf("0x%x", caps),
	}
	c.Equal = item.Bool{
		Name:        "Equal to Docker's Default capability",
		Description: fmt.Sprintf("0x%x", caps),
		Result:      caps == standardCaps,
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
	fmt.Println(c)
	return
}
