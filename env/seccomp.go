package env

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/pkg/seccomp"
	"github.com/ctrsploit/ctrsploit/util"
	"fmt"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const (
	CommandSeccompName = "seccomp"
)

// Seccomp
// reference: https://lwn.net/Articles/656307/
func Seccomp() (err error) {
	seccompMode, _, err := seccomp.GetStatus()
	awesome_error.CheckFatal(err)
	info := "===========Seccomp========="
	info += fmt.Sprintf("\nkernel supported: %v", util.ColorfulTickOrBallot(seccomp.CheckSupported()))
	info += fmt.Sprintf("\nseccomp enabled in current container: %v", util.ColorfulTickOrBallot(seccompMode > 0))
	if seccompMode > 0 {
		var seccompModeString string
		switch seccompMode {
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
		info += fmt.Sprintf("\nseccomp mode: %v", seccompModeString)
	}
	log.Logger.Info(info)
	return
}
