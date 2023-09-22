package apparmor

import (
	"github.com/ctrsploit/ctrsploit/internal/colorful"
	"github.com/ctrsploit/ctrsploit/pkg/apparmor"
	"github.com/ctrsploit/ctrsploit/pkg/lsm"
	"github.com/ssst0n3/awesome_libs"
)

var (
	tplApparmor = `
===========Apparmor===========
{.kernel}  Kernel Supported
{.container}  Container Enabled
{.details}`
	tplApparmorDetails = `---
Profile:	{.profile}
Mode:		{.mode}
`
)

func Apparmor() (err error) {
	details := ""
	enabled := apparmor.IsEnabled()
	if enabled {
		current, err := lsm.Current()
		if err != nil {
			return err
		}
		mode, err := apparmor.Mode()
		if err != nil {
			return err
		}
		details = awesome_libs.Format(tplApparmorDetails, awesome_libs.Dict{
			"profile": colorful.Safe(current),
			"mode":    colorful.Safe(mode),
		})
	}
	info := awesome_libs.Format(tplApparmor, awesome_libs.Dict{
		"kernel":    colorful.TickOrBallot(apparmor.IsSupport()),
		"container": colorful.TickOrBallot(enabled),
		"details":   details,
	})
	print(info)
	return
}
