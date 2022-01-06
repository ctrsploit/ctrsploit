package apparmor

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/apparmor"
	"github.com/ctrsploit/ctrsploit/pkg/lsm"
	"github.com/ctrsploit/ctrsploit/util"
	"github.com/fatih/color"
)

func Apparmor() (err error) {
	info := "===========Apparmor========="
	info += fmt.Sprintf("\nKernel Supported: %s", util.ColorfulTickOrBallot(apparmor.IsSupport()))
	enabled := apparmor.IsEnabled()
	info += fmt.Sprintf("\nContainer Enabled: %s", util.ColorfulTickOrBallot(enabled))
	if enabled {
		current, err := lsm.Current()
		if err != nil {
			return err
		}
		info += "\nApparmor Profile: " + color.HiGreenString(current)
		mode, err := apparmor.Mode()
		if err != nil {
			return err
		}
		info += "\nApparmor Mode: " + color.HiGreenString(mode)
	}
	fmt.Printf("%s\n\n", info)
	return
}
