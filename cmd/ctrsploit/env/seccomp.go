package env

import (
	"ctrsploit/env/seccomp"
	"github.com/urfave/cli/v2"
)

var seccompCommand = &cli.Command{
	Name:    seccomp.CommandName,
	Aliases: []string{"s"},
	Usage:   "show the seccomp info",
	Action: func(context *cli.Context) (err error) {
		err = seccomp.Seccomp()
		if err != nil {
			return
		}
		return
	},
}
