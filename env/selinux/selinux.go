package selinux

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/selinux"
	"github.com/ctrsploit/ctrsploit/util"
)

func Selinux() (err error) {
	info := "===========SELinux========="
	info += fmt.Sprintf("\nEnabled: %s", util.ColorfulTickOrBallot(selinux.IsEnabled()))
	mode := selinux.Translate(selinux.Mode())
	info += fmt.Sprintf("\nmode: %v", mode)
	mountPoint := selinux.GetSelinuxMountPoint()
	info += fmt.Sprintf("\nSELinux filesystem mount point: %v", mountPoint)
	fmt.Printf("%s\n\n", info)
	return
}
