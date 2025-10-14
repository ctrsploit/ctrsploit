package env

import (
	"github.com/ctrsploit/ctrsploit/env/sysctl"
	"github.com/urfave/cli/v2"
)

var Sysctl = &cli.Command{
	Name:    sysctl.CommandName,
	Aliases: []string{},
	Usage:   "display sysctl information",
	Action: func(context *cli.Context) (err error) {
		err = sysctl.Print()
		if err != nil {
			return
		}
		return
	},
}
