package env

import (
	"ctrsploit/env/apparmor"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/urfave/cli/v2"
)

var apparmorCommand = &cli.Command{
	Name:    apparmor.CommandName,
	Aliases: []string{"a"},
	Usage:   "show the apparmor info",
	Action: func(context *cli.Context) (err error) {
		log.Logger.Debug("")
		err = apparmor.Apparmor()
		if err != nil {
			return
		}
		return
	},
}
