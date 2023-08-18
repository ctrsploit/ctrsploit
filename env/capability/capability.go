package capability

import (
	"fmt"
	"github.com/containerd/containerd/pkg/cap"
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/pkg/capability"
	"github.com/ctrsploit/ctrsploit/util"
	"github.com/fatih/color"
	"github.com/ssst0n3/awesome_libs"
)

const (
	CommandCapabilityName        = "capability"
	standardCaps          uint64 = 0xa80425fb
	PRIVILEGED_CAPS       uint64 = 0x3fffffffff
)

func getInfoFromCaps(caps uint64) (info string) {
	capsParsed, _ := cap.FromBitmap(caps)
	standard := "default(0xa80425fb)"
	if caps == standardCaps {
		standard = color.HiGreenString(" = ") + standard
	} else {
		standard = color.HiRedString(" != ") + standard
	}
	info += awesome_libs.Format(`
{.title_caps}
{.caps}

{.title_caps_parsed}
{.caps_parsed}
`, awesome_libs.Dict{
		"title_caps":        util.TitleWithFgWhiteBoldUnderline("[caps]"),
		"caps":              fmt.Sprintf("0x%x%s", caps, standard),
		"title_caps_parsed": util.TitleWithFgWhiteBoldUnderline("[parsed]"),
		"caps_parsed":       capsParsed,
	})

	if caps != standardCaps {
		capsDiff, _ := cap.FromBitmap(caps & (^standardCaps))
		info += "\n" + util.TitleWithFgWhiteBoldUnderline("[Additional Capabilities]")
		info += color.HiRedString(fmt.Sprintf("\n%q", capsDiff))
	}
	
	// Checking for a privileged container
	if caps & (^PRIVILEGED_CAPS) != 0 {
		info += "\n" + util.TitleWithFgWhiteBoldUnderline("[Privileged]")
		info += color.HiRedString("\nWARNING: Possible Privileged Container Found!")
	}

	return
}

func Capability() (err error) {
	info := "===========Capability========="
	{ // for pid 1
		caps, err := capability.GetPid1Capability()
		if err != nil {
			return err
		}
		info += "\n" + util.TitleWithBgWhiteBold("pid 1") + getInfoFromCaps(caps)
	}
	info += "\n"
	{
		// for current process
		caps, err := capability.GetCurrentCapability()
		if err != nil {
			return err
		}
		info += "\n" + util.TitleWithBgWhiteBold("current process") + getInfoFromCaps(caps)
	}
	log.Logger.Info(info)
	return
}
