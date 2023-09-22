package selinux

import (
	"github.com/ctrsploit/ctrsploit/internal/colorful"
	"github.com/ctrsploit/ctrsploit/pkg/selinux"
	"github.com/ssst0n3/awesome_libs"
)

var (
	tplSelinux = `
===========SELinux===========
{.enabled}  Enabled
{.details}`
	tplSelinuxDetails = `---
Mode:			{.mode}
mount point:	{.mount}
`
)

func Selinux() (err error) {
	details := ""
	enabled := selinux.IsEnabled()
	if enabled {
		details = awesome_libs.Format(tplSelinuxDetails, awesome_libs.Dict{
			"mode":  selinux.Translate(selinux.Mode()),
			"mount": selinux.GetSelinuxMountPoint(),
		})
	}
	info := awesome_libs.Format(tplSelinux, awesome_libs.Dict{
		"enabled": colorful.TickOrBallot(enabled),
		"details": details,
	})
	print(info)
	return
}
