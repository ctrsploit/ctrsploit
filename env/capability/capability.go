package capability

import (
	"fmt"
	"github.com/containerd/containerd/pkg/cap"
	"github.com/ctrsploit/ctrsploit/internal/colorful"
	"github.com/ctrsploit/ctrsploit/pkg/capability"
	"github.com/ssst0n3/awesome_libs"
)

const (
	CommandCapabilityName        = "capability"
	standardCaps          uint64 = 0xa80425fb
	PRIVILEGED_CAPS       uint64 = 0x3fffffffff
)

var (
	tpl = `
===========Capability(process 1)===========
{.pid1}
===========Capability(process current)===========
{.current}`
	tplCaps = `[Capabilities]
{.caps} {.equal} 0xa80425fb(docker's default caps)
{.diff}`
	tplDiff = `[Additional Capabilities]
{.diff}
`
)

func getInfoFromCaps(caps uint64) (info string) {
	equal := colorful.Safe("=")
	if caps != standardCaps {
		equal = colorful.Danger("!=")
	}
	diff := ""
	if caps != standardCaps {
		capsDiff, _ := cap.FromBitmap(caps & (^standardCaps))
		diff = awesome_libs.Format(tplDiff, awesome_libs.Dict{
			"diff": colorful.Danger(fmt.Sprintf("\n%q", capsDiff)),
		})
	}
	info = awesome_libs.Format(tplCaps, awesome_libs.Dict{
		"caps":  fmt.Sprintf("0x%x", caps),
		"equal": equal,
		"diff":  diff,
	})

	// Checking for a privileged container
	if caps&PRIVILEGED_CAPS == PRIVILEGED_CAPS {
		info += "\n" + colorful.Title("[Privileged]")
		info += colorful.Danger("\nWARNING: Possible Privileged Container Found!\n")
	}

	return
}

func Capability() (err error) {
	caps, err := capability.GetPid1Capability()
	if err != nil {
		return err
	}
	pid1 := getInfoFromCaps(caps)

	caps, err = capability.GetCurrentCapability()
	if err != nil {
		return err
	}
	current := getInfoFromCaps(caps)

	info := awesome_libs.Format(tpl, awesome_libs.Dict{
		"pid1":    pid1,
		"current": current,
	})

	print(info)
	return
}
