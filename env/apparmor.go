package env

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/pkg/apparmor"
	"github.com/ctrsploit/ctrsploit/pkg/lsm"
	"github.com/ctrsploit/ctrsploit/util"
	"github.com/fatih/color"
)

const (
	CommandApparmorName = "apparmor"
)

func Apparmor() (err error) {
	info := "===========Apparmor========="
	info += "\nKernel Supported: " + util.ColorfulTickOrBallot(apparmor.IsSupport())
	enabled := apparmor.IsEnabled()
	info += "\nContainer Enabled: " + util.ColorfulTickOrBallot(enabled)
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
	log.Logger.Info(info)
	return
}
