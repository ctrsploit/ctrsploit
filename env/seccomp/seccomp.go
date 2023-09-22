package seccomp

import (
	"github.com/ctrsploit/ctrsploit/internal/colorful"
	"github.com/ctrsploit/ctrsploit/pkg/seccomp"
	"github.com/ssst0n3/awesome_libs"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const (
	CommandSeccompName = "seccomp"
)

var (
	tpl = `
===========Seccomp===========
{.supported}  Kernel Supported
{.enabled}  Container Enabled
{.details}
`
	tplDetails = `---
mode:	{.mode}`
)

// Seccomp
// reference: https://lwn.net/Articles/656307/
func Seccomp() (err error) {
	seccompMode, _, err := seccomp.GetStatus()
	awesome_error.CheckFatal(err)
	details := ""
	if seccompMode > 0 {
		details = awesome_libs.Format(tplDetails, awesome_libs.Dict{
			"mode": TranslateMode(seccompMode),
		})
	}
	info := awesome_libs.Format(tpl, awesome_libs.Dict{
		"supported": colorful.TickOrBallot(seccomp.CheckSupported()),
		"enabled":   colorful.TickOrBallot(seccompMode > 0),
		"details":   details,
	})
	print(info)
	return
}

func TranslateMode(mode int) (seccompModeString string) {
	if mode > 0 {
		switch mode {
		case 1:
			// The first version of seccomp was merged in 2005 into Linux 2.6.12.
			// It was enabled by writing a "1" to /proc/PID/seccomp.
			// Once that was done, the process could only make four system calls: read(), write(), exit(),
			// and sigreturn().
			seccompModeString = "strict"
		case 2:
			// Things were calm in seccomp land for the next five years or so until
			// "seccomp mode 2" (or "seccomp filter mode") was added to Linux 3.5 in 2012.
			// It added a second mode for seccomp: SECCOMP_MODE_FILTER. Using that mode,
			// processes can specify which system calls are permitted.
			seccompModeString = "filter"
		}
	}
	return
}
